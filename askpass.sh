#!/usr/bin/env bash
# A small shell script to provide an alternative to the Go program
# rewritten from the fish script by copilot, and it's working.

prompt="$*"
if [ -z "$prompt" ]; then
    echo "I'm a tiny ssh askpass helper and I'm not designed to work standalone"
    exit 1
fi

get_passphrase_from_rbw() {
    local key="$1"
    local pass
    pass=$(rbw get -f "passphrase" "$key")
    if [ -z "$pass" ]; then
        echo "No passphrase found, exiting"
        return 1
    fi
    echo "$pass"
}

ask_trust() {
    read -r -p "$prompt" response
    case "$response" in
        yes|y|Y)
            echo "yes"
            ;;
        no|n|N)
            echo "no"
            ;;
        *)
            echo "$response"
            ;;
    esac
}

unknown() {
    read -r -p "$prompt" response
    echo "$response"
}

if [[ "$prompt" =~ for\ \(?:\ key\)?\ \'?\(.+?\)\'?: ]]; then
    echo "The host is known, asking the passphrase" >&2
    key=$(echo "$prompt" | grep -oP "for(?: key)? '?(.+?)'?: " | sed -e "s/ for(?: key)? '?//" -e "s/'?: //")
    get_passphrase_from_rbw "$key"
elif [[ "$prompt" == *"The authenticity"* ]]; then
    echo "The host authenticity can't be established, asking for trust" >&2
    ask_trust
else
    echo "The prompt is unknown, asking user" >&2
    unknown
fi