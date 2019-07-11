#!/usr/bin/env bash

DFB_PATH="$HOME/.dfb"
PROGRAM=$(basename "${0}")
VERSION=1.0

main() {
    if [[ "${1:-}" =~ ^-h|--help$  ]]
    then
        print_help
    elif [[ "${1:-}" =~ ^-v|--version$  ]]
    then
        print_version
    elif [ "${1:-}" == "list" ]
    then
        list_domains
    else
        print_help
    fi
}

print_help() {
    cat <<HEREDOC

_________/\\\\\\__________/\\\\\\\\\\___/\\\\\\________
 ________\\/\\\\\\________/\\\\\\///___\\/\\\\\\________
  ________\\/\\\\\\_______/\\\\\\_______\\/\\\\\\________
   ________\\/\\\\\\____/\\\\\\\\\\\\\\\\\\____\\/\\\\\\________
    ___/\\\\\\\\\\\\\\\\\\___\\////\\\\\\//_____\\/\\\\\\\\\\\\\\\\\\__
     __/\\\\\\////\\\\\\______\\/\\\\\\_______\\/\\\\\\////\\\\\\_
      _\\/\\\\\\__\\/\\\\\\______\\/\\\\\\_______\\/\\\\\\__\\/\\\\\\_
       _\\//\\\\\\\\\\\\\\\\\\______\\/\\\\\\_______\\/\\\\\\\\\\\\\\\\\\\\_
        __\\/////////_______\\///________\\//////////__

Domain based Filesystem Backup.

Usage:
  ${PROGRAM} [subcommand] [<arguments>]

Available Commands:
  list          List domains.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
HEREDOC
}

print_version() {
    echo $VERSION
}

verify_env() {
    if [[ ! "$OSTYPE" == "darwin"* ]]; then
        echo "dfb is only availible for macOS"
        exit 1
    fi

    if [ "$(command -v restic)" == "" ]; then
        echo "restic is not installed, visit https://github.com/restic/restic"
        exit 1
    fi

    if [ ! -d "/Library/Filesystems/osxfuse.fs" ]; then
        echo "FUSE for macOS is not installed, visit https://github.com/osxfuse/osxfuse"
        exit 1
    fi

    if [ ! -d "$DFB_PATH" ]; then
        echo "creating dfb root directory at $DFB_PATH"
        mkdir "$DFB_PATH"
    fi
}

list_domains() {
    verify_env
}

main "$@"
