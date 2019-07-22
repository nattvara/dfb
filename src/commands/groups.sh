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
