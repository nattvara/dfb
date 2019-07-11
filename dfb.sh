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
    elif [ "${1:-}" == "domains" ]
    then
        domains "$@"
    elif [ "${1:-}" == "backup" ]
    then
        backup "$@"
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
  domain      Domain commands.
  backup      Backup a group of domains to a repo.

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
    elif [ "${2:-}" == "add" ]
    then
        add_group
    elif [ "${2:-}" == "ls" ]
    then
        list_groups
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
  repos     List restic repos for a group.
  add-repo  Add restic repo for a group.

Options:
  -h --help     Show this screen.
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

validate_repo() {
    if [ ! -f "$DFB_PATH/$1/repos/$2" ]; then
        echo "please provide a valid repo"
        exit 1
    fi
}

domains() {
    verify_env
    if [[ "${2:-}" =~ ^-h|--help$  ]]
    then
        print_domains_help
    elif [ "${2:-}" == "add" ]
    then
        add_domain "$3" "$4" "$5"
    elif [ "${2:-}" == "ls" ]
    then
        list_domains
    else
        print_domains_help
    fi
}

print_domains_help() {
    cat <<HEREDOC
Domains to backup.

A domain is a directory in the home directory to backup,
this could be a symlink to some other directory on another
volume.

Usage:
  ${PROGRAM} domains <subcommand> [parameters]

Available Commands:
  ls        List domains.
  add       Add new domain.

Options:
  -h --help     Show this screen.
HEREDOC
}

list_domains() {
    printf "Domains: \n\n"

    cd "$DFB_PATH"
    find . -type d ! -path . -maxdepth 1 -print0 |
    while IFS= read -r -d '' group; do
        cd "$DFB_PATH/$group/domains"
        group="$(echo $group | sed -e 's/^\.\///g')"

        find . -type f -print0 |
        while IFS= read -r -d '' domain; do
            domain="$(echo $domain | sed -e 's/^\.\///g')"
            echo "$group:$domain"
        done
    done
}

add_domain() {
    if [[ $1 == "help" ]]; then
        echo "Usage:"
        echo "  $ $PROGRAM domains add [group] [domain] [<symlink>]"
        exit
    fi
    validate_group "$1"

    if [[ $2 == "" ]]; then
        echo "please provide a domain"
        exit 1
    fi
    if [ ! -d $2 ]; then
        echo "domain is not a valid directory"
        exit 1
    fi

    domain=$(basename "$2")
    content=$(cat <<CONTENT
path: $2
symlink:
exclusions: node_modules vendor
CONTENT
)
    echo "$content" > "$DFB_PATH/$1/domains/$domain"
}

backup_domain() {
    password="$1"
    repo_path="$2"
    domain="$3"
    domain_path=$(cat "./$domain" | ggrep -E 'path' | egrep -o '[^:]+$' | tr -d '[:space:]')

    cd $domain_path
    echo "$password" | restic -r $repo_path backup . --tag "$domain" --json
}

backup() {
    if [[ $2 == "help" ]]; then
        echo "Usage:"
        echo "  $PROGRAM backup [group] [repo] [<timestamp-file>]"
        exit
    fi
    group=$2
    repo_name=$3
    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")
    validate_group $group
    validate_repo $group $repo_name
    domains_directory="$DFB_PATH/$group/domains"

    promt_for_password
    verify_password $password $repo_path

    cd $domains_directory
    find . -type f -print0 |
    while IFS= read -r -d '' domain; do
        domain=$(echo $domain | sed -e 's/^\.\///g')

        printf "Backing up "
        echo $domain

        backup_domain $password $repo_path $domain
        cd $domains_directory
    done
}

promt_for_password() {
    password=$(osascript <<END
set x to display dialog "What is your password?" default answer "" with hidden answer
set y to (text returned of x)
END
    )
    if [ -z "$(echo ${password//[[:blank:]]/})" ]; then
        osascript -e "display notification with title \"Password cannot be empty\""
        exit 1
    fi
}

verify_password() {
    password="$1"
    repo_path="$2"

    if echo "$password" | restic -r "$repo_path" key list; then
        echo "valid password"
    else
        osascript -e "display notification \"for $repo_path\" with title \"Invalid password\""
        exit 1
    fi
}

main "$@"
