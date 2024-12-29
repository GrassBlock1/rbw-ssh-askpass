#!/usr/bin/env fish
# a attempt to make a small shell script to provide a alternative to the go program
# currently it's not working, it can't read user input.

set prompt (echo $argv | string join " ")
if test -z $prompt
   echo "I'm a tiny ssh askpass helper and I'm not designed to work standalone"
   return 1
end

function get_passphrase_from_rbw --description 'get passphrase from rbw'
    set key $argv[1]
    set pass (rbw get -f "passphrase" $key)
    if test -z $pass
        echo "No passphrase found, exiting"
        return 1
    end
    echo $pass
end

function ask_trust
    read -l response --prompt-str="$prompt" ;or return 1
    switch response
        case yes y Y
            echo "yes"
        case no n N
            echo "no"
        case *
            echo $response
    end
end

function unknown
    read -l response --prompt-str="$prompt" or return 1
end

if string match -q -r " for(?: key)? '?(.+?)'?: " $prompt
    echo "The host is known, asking the passphrase" 1>&2
    string match -q -r " for(?: key)? '?(.+?)'?: " $prompt $key
    get_passphrase_from_rbw $key
else if string match -q -- "*The authenticity*" $prompt
    echo "the host authenticity can't be established, asking for trust" 1>&2
    ask_trust $prompt
else
    echo "the prompt is unknown, asking user" 1>&2
    unknown $prompt
end



