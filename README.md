# SSH askpass helper with Bitwarden

This script is a simple SSH askpass helper that uses Bitwarden to retrieve the passphrase for the SSH key.

This will help you to unlock private keys with passphrase stored in your Bitwarden vault.

This does **not** store and extract private keys. Private keys are, well, private, and should be generated uniquely on
each machine.

## Requirements

- [rbw](https://github.com/doy/rbw) installed and configured in your `$PATH`
- `SSH_ASKPASS` environment variable set to the path of this script
- An item named with private_key file name and it has "passphrase" field in Bitwarden vault

## Usage

1. clone the repository
2. run `go build` to build the binary
   (or you can go to release to download it)
3. set the `SSH_ASKPASS` environment variable to the path of the binary, for example:

```shell
# for fish
# ~/.config/fish/config.fish
set -gx SSH_ASKPASS /path/to/rbw-ssh-askpass
set -gx SSH_ASKPASS_REQUIRE prefer
# bash
# ~/.bashrc
export SSH_ASKPASS=/path/to/rbw-ssh-askpass
export SSH_ASKPASS_REQUIRE=prefer
```

## License

MIT