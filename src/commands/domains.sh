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
