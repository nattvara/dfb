#!/usr/bin/env bash

DFB_PATH="$HOME/.dfb"
PROGRAM=$(basename "${0}")
VERSION=1.0

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

print_version() {
    echo $VERSION
}

main() {
    if [[ "${1:-}" =~ ^-h|--help$  ]]
    then
        print_main_help
    elif [[ "${1:-}" =~ ^-v|--version$  ]]
    then
        print_version
    elif [ "${1:-}" == "groups" ]
    then
        groups "$@"
    else
        print_main_help
    fi
}

print_main_help() {
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
  ${PROGRAM} <command> <subcommand> [parameters]

Available Commands:
  groups      Group commands.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
HEREDOC
}

groups() {
    verify_env
    if [[ "${2:-}" =~ ^-h|--help$  ]]
    then
        print_groups_help
    elif [[ "${2:-}" =~ ^-v|--version$  ]]
    then
        print_version
    elif [ "${2:-}" == "add" ]
    then
        add_group
    elif [ "${2:-}" == "ls" ]
    then
        list_groups
    elif [ "${2:-}" == "domains" ]
    then
        list_group_domains
    else
        print_groups_help
    fi
}

print_groups_help() {
    cat <<HEREDOC
Domain groups.

A group contians a number of domains, and restic repositories
to backup those domains to.

Usage:
  ${PROGRAM} groups <subcommand> [parameters]

Available Commands:
  ls        List groups.
  add       Add new group.
  domains   List domains for a group.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
HEREDOC
}

list_groups() {
    echo "list_groups"
}

add_group() {
    echo "add_group"
}

list_group_domains() {
    echo "list domains for group"
}

main "$@"
