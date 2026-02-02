#!/usr/bin/env bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

SSH_ASKPASS_REQUIRE=force SSH_ASKPASS=$SCRIPT_DIR/rbw-ssh-askpass SSH_AUTH_SOCK="" ssh -o UserKnownHostsFile=/dev/null "$@"