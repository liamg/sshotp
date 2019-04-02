# SSHOTP: Automatic entry of non-interactive passwords

Autopass is essentially a go implementation of [sshpass](https://linux.die.net/man/1/sshpass), though unlike sshpass it doesn't restrict itself to SSH logins. It can supply a password to any process with an identifiable password prompt.

**Do not use this unless you understand the risks involved - ssh prompts for a password interactively for a reason!**

The original use case for this was needing to automate the acquisition and use of an SSH OTP (via vault) in a nice script.

## Requirements

- Mac/Linux
- Go 1.11+ (to build)

## Install

```bash
go get -u github.com/liamg/sshotp
```

## Example

```bash
sshotp --password mypassword123 "ssh me@myserver.mine -p 2222"
```

## Usage

```
Usage:
  sshotp [flags]

Flags:
      --disable-ssh-host-confirm   sshotp will automatically confirm the authenticity of SSH hosts unless this option is specified
      --env                        use value of $SSHOTP environment variable as password
  -h, --help                       help for autopass
      --password string            plaintext password (not recommended)
      --timeout duration           timeout length to wait for prompt/confirmation (default 10s)
```
