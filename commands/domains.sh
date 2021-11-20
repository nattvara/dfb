#
# Summary: Domains command
#
# A domain is a directory in the home directory to backup,
# this could be a symlink to some other directory on another
# volume. This command will display and manipulate domains
# in a group.

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
    elif [ "${2:-}" == "not-added" ]
    then
        list_not_added "$3"
    elif [ "${2:-}" == "rm" ]
    then
        remove_domain "$3" "$4"
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
  ls            List domains.
  not-added     List directories not added as domains in home directory.
  add           Add new domain.
  rm            Remove a domain.

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

        find . -type f -print0 | sort -z |
        while IFS= read -r -d '' domain; do
            repos=$(cat "$domain" | ggrep -E 'repos:' | egrep -o '[^:]+$' | tr -d '[:space:]')
            symlink_path=$(cat "$domain" | ggrep -E 'symlink:' | egrep -o '[^:]+$' | tr -d '[:space:]')
            domain="$(echo $domain | sed -e 's/^\.\///g')"
            printf "$group:$domain$(head -c 100 < /dev/zero | tr '\0' ' ')" | awk '{print substr($0, 1, 100)}' | tr -d '\n'
            printf "$repos$(head -c 100 < /dev/zero | tr '\0' ' ')" | awk '{print substr($0, 1, 30)}' | tr -d '\n'
            printf "$symlink_path\n"
        done
    done
}

list_not_added() {
    if [[ $1 == "help" ]]; then
        echo "Usage:"
        echo "  $ $PROGRAM domains not-added [group]"
        exit
    fi
    group=$1

    printf "Directories not backed up: \n\n"

    cd "$HOME"
    find . -print0 -maxdepth 1 | sort -z |
    while IFS= read -r -d '' file; do
        domain="$(echo $file | sed -e 's/^\.\///g')"
        if [ ! -f "$DFB_PATH/$group/domains/$domain" ]; then
            printf "~/$domain\n"
        fi
    done
}

add_domain() {
    if [[ $1 == "help" ]]; then
        echo "Usage:"
        echo "  $ $PROGRAM domains add [group] [domain] [<symlink>]"
        exit
    fi
    group=$1
    path=$2
    symlink=$3
    domain=$(basename "$path")

    validate_group $group

    if [[ $path == "" ]]; then
        echo "please provide a domain"
        exit 1
    fi
    if [ ! -d $path ] && [ ! -f $path ] && [[ $symlink == "" ]]; then
        echo "domain is not a valid directory or file"
        exit 1
    fi

    if [[ $symlink != "" ]]; then
        create_symlink $domain "$symlink"
    fi

    content=$(cat <<CONTENT
path: $path
symlink: $symlink
exclusions: **/node_modules **/.DS_Store
repos: *
CONTENT
)
    echo "$content" > "$DFB_PATH/$group/domains/$domain"
}

create_symlink() {
    domain=$1
    symlink=$2

    symlinks="$DFB_PATH/$group/symlinks"

    if [ ! -d $symlinks ]; then
        echo "creating symlinks directory"
        mkdir $symlinks
    fi

    if [ -f "$symlinks/$domain" ]; then
        rm "$symlinks/$domain"
    fi

    if [ -f "$HOME/$domain" ] && [ ! -L "$HOME/$domain" ]; then
        echo "$HOME/$domain already exists, please remove manually if a link should exist here"
        exit 1
    elif [ -d "$HOME/$domain" ] && [ ! -L "$HOME/$domain" ]; then
        echo "$HOME/$domain already exists, please remove manually if a link should exist here"
        exit 1
    fi

    ln -vs "$symlink" "$symlinks/$domain"
    ln -vs "$symlinks/$domain" "$HOME/$domain"
}

remove_domain() {
    if [[ $1 == "help" ]]; then
        echo "Usage:"
        echo "  $ $PROGRAM domains rm [group] [domain]"
        exit
    fi
    group=$1
    domain=$2

    validate_group $group
    validate_domain $group $domain

    echo "deleting record of domain $domain"
    rm "$DFB_PATH/$group/domains/$domain"

    if [ -L "$DFB_PATH/$group/symlinks/$domain" ]; then
        echo "deleting symlink to real directory"
        rm "$DFB_PATH/$group/symlinks/$domain"
    fi

    printf "\n\nNOTE: this does not delete the actual directory"
    printf " it will simply not be included in any more backups"
    printf " neither will it be removed from previous backups\n\n"
    printf "you will have to delete the data of \"$domain\" yourself\n"
}
