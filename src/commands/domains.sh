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
    group=$1
    path=$2
    symlink=$3
    domain=$(basename "$path")

    validate_group $group

    if [[ $path == "" ]]; then
        echo "please provide a domain"
        exit 1
    fi
    if [ ! -d $path ] && [[ $symlink == "" ]]; then
        echo "domain is not a valid directory"
        exit 1
    fi

    if [[ $symlink != "" ]]; then
        create_symlink $domain $symlink
    fi

    content=$(cat <<CONTENT
path: $path
symlink: $symlink
exclusions: node_modules vendor
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

    ln -s $symlink "$symlinks/$domain"
    ln -s "$symlinks/$domain" "$HOME/$domain"
}
