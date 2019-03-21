package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liamg/sshotp/app"
	"github.com/spf13/cobra"
)

var password string
var useEnv bool
var timeout time.Duration
var disableConfirmHostAuthenticity bool

var rootCmd = &cobra.Command{
	Use:   "sshotp",
	Short: "Enter passwords to commands non-interactively",
	Long:  `SSHOTP is essentially a go implementation of sshpass (https://linux.die.net/man/1/sshpass), though unlike sshpass it doesn't restrict itself to SSH logins. It can supply a password to any process with an identifiable password prompt.`,
	Run: func(cmd *cobra.Command, args []string) {

		command := strings.Join(args, " ")

		if command == "" {
			fmt.Println("You must specify a command.")
			os.Exit(1)
		}

		if useEnv {
			password = os.Getenv("AUTOPASS")
		}

		if err := app.Run(command, password, "assword", "denied", timeout, !disableConfirmHostAuthenticity); err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "plaintext password (not recommended)")
	rootCmd.PersistentFlags().BoolVar(&useEnv, "env", false, "use value of $AUTOPASS environment variable as password")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", time.Second*10, "timeout length to wait for prompt/confirmation")
	rootCmd.PersistentFlags().BoolVar(&disableConfirmHostAuthenticity, "disable-ssh-host-confirm", false, "autopass will automatically confirm the authenticity of SSH hosts unless this option is specified")
}
