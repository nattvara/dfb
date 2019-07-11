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

    if [ "$(command -v ggrep)" == "" ]; then
        echo "GNU grep is not installed, run: brew install grep"
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
    elif [ "${2:-}" == "repos" ]
    then
        list_group_repos "$3"
    elif [ "${2:-}" == "add-repo" ]
    then
        add_group_repo "$3"
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
  repos     List restic repos for a group.
  add-repo  Add restic repo for a group.

Options:
  -h --help     Show this screen.
  -v --version  Print version information.
HEREDOC
}

list_groups() {
    printf "Groups: \n\n"

    cd $DFB_PATH
    find . -type d ! -path . -maxdepth 1 -print0 |
    while IFS= read -r -d '' dir; do
        dir="${dir//.}"
        dir="${dir/\/}"
        echo "$dir"
    done
}

add_group() {
    printf "Enter name of new group: "
    read group

    cd $DFB_PATH

    if [ -d "$group" ]; then
        echo "group already exists!"
        exit 1
    fi

    mkdir "$group"
    mkdir "$group/repos"
    mkdir "$group/domains"
}

list_group_domains() {
    echo "list domains for group"
}

list_group_repos() {
    validate_group "$@"
    cd "$DFB_PATH/$1/repos/"
    find . -type f -maxdepth 1 -print0 |
    while IFS= read -r -d '' file; do
        printf "$(echo $file | sed -e 's/^\.\///g'): "
        cat $file
    done
}

add_group_repo() {
    validate_group "$@"

    printf "Enter name of repo: "
    read name
    printf "Enter repo path: "
    read repo

    echo "$repo" > "$DFB_PATH/$1/repos/$name"
}

validate_group() {
    if [ "$1" == "" ]; then
        echo "please provide a group."
        exit 1
    fi

    if [ ! -d "$DFB_PATH/$1" ]; then
        echo "please provide a valid group."
        exit 1
    fi
}

main "$@"
