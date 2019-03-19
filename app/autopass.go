package app

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// Run attempts to run the provided command and insert the given password when prompted.
// expectedPrompt is the string to treat as trhe password prompt e.g. "Password: "
// expectedFailure is the string to treat as an indication of failure e.g. "Incorrect password"
func Run(cmd string, password string, expectedPrompt string, expectedFailure string, timeout time.Duration, autoConfirmHostAuthenticity bool) error {

	c := exec.Command("/bin/bash")

	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	defer func() { _ = ptmx.Close() }()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

	if _, err := ptmx.Write([]byte(cmd + "; exit $?\n")); err != nil {
		return err
	}

	errChan := make(chan error)
	readyChan := make(chan struct{})

	go func() {
		data := ""
		buf := make([]byte, 4096)
		confirmed := false
		entered := false
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				errChan <- err
				break
			}
			if n == 0 {
				continue
			}
			data += string(buf[:n])
			if autoConfirmHostAuthenticity && !confirmed && strings.Contains(data, "The authenticity of host ") {
				confirmed = true
				data = ""
				ptmx.Write([]byte("yes\n"))
			} else if !entered && strings.Contains(data, expectedPrompt) {
				entered = true
				data = ""
				ptmx.Write([]byte(password + "\n"))
			} else if entered && len(data) > 5 {
				if strings.Contains(data, expectedPrompt) || strings.Contains(data, expectedFailure) {
					errChan <- fmt.Errorf("authentication failure")
					break
				}
				readyChan <- struct{}{}
				break
			}
		}
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-readyChan:
		go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
		_, _ = io.Copy(os.Stdout, ptmx)
	case err := <-errChan:
		return err
	case <-timer.C:
		return fmt.Errorf("timed out waiting for prompt")
	}

	return nil
}
