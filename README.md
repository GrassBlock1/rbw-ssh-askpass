# SSH askpass helper with Bitwarden

This script is a simple SSH askpass helper that uses Bitwarden to retrieve the passphrase for the SSH key.

This will help you to unlock private keys with passphrase stored in your Bitwarden vault.

This does **not** store and extract private keys. Private keys are, well, private, and should be generated uniquely on
each machine.

## Requirements

- [rbw](https://github.com/doy/rbw) installed and configured in your `$PATH`
- `SSH_ASKPASS` environment variable set to the path of this script
- In your Bitwarden vault, an item named with `private_key` file name(eg: `id_rsa`) exists, and it has a "passphrase" field.

## Usage
> [!NOTE]
> A non-binary version is available at [askpass.sh](askpass.sh).
> To get started, simply download it in somewhere and set the `SSH_ASKPASS` environment variable to the path of the script.

> [!NOTE]
> The script is not designed to use with git + https.
> For Git credential helper, you can use [git-credential-rbw](https://github.com/doy/rbw/blob/main/bin/git-credential-rbw) to use & store them.

1. clone the repository
2. run `go build` to build the binary
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

### Working with rbw's ssh agent
After adding a key by following [Bitwarden SSH Agent | Bitwarden](https://bitwarden.com/help/ssh-agent/), rbw-agent should work after running `rbw sync`.

However, if you put things like `IdentityFile ~/.ssh/id_ed25519` in your SSH config, you may be still asked to input passphrase. If so, editing the title of SSH key item to the private key filename and add a hidden field named 'passphrase' of should work.

## License

MIT