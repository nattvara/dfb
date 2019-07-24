#
# Summary: Backup command
#
# The backup command will start a backup of one or more
# domains.

backup_domain() {
    domain_path=$(cat "./$domain" | ggrep -E 'path' | egrep -o '[^:]+$' | tr -d '[:space:]')
    symlink=$(cat "./$domain" | ggrep -E 'symlink' | egrep -o '[^:]+$' | tr -d '[:space:]')
    exclusions=$(cat "./$domain" | ggrep -E 'exclusions' | egrep -o '[^:]+$' | tr " " "\n")
    repos=$(cat "./$domain" | ggrep -E 'repos' | egrep -o '[^:]+$' | tr -d '[:space:]')

    if [[ $symlink != "" ]]; then
        if [ ! -d $symlink ]; then
            print_domain_unavailable $domain
            return
        fi

        cd $symlink
    elif [ -f $domain_path ]; then
        parent_dir=$(dirname $domain_path)
        if [ ! -d $parent_dir ]; then
            print_domain_unavailable $domain
            return
        fi
        cd $parent_dir
    else
        if [ ! -d $domain_path ]; then
            print_domain_unavailable $domain
            return
        fi
        cd $domain_path
    fi

    if [ -f $domain_path ]; then
        restic_backup_path=$(basename $domain_path)
    else
        restic_backup_path="."
    fi

    if [[ $repos != "*" ]] && [[ ",$repos," != *",$repo_name,"* ]]; then
        print_not_this_repo $domain $repo_name
        return
    fi

    echo "$exclusions" > /tmp/dfb_exclusions
    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "begin"
        echo -n "$password" | restic -r $repo_path backup "$restic_backup_path" --tag "$domain"  --exclude-file /tmp/dfb_exclusions --verbose --json >> /tmp/dfb-progress
    else
        echo -n "$password" | restic -r $repo_path backup "$restic_backup_path" --tag "$domain"  --exclude-file /tmp/dfb_exclusions --verbose --json | dfb-progress-parser "$group" "$domain"
        printf "\r"
    fi
}

print_domain_unavailable() {
    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "unavailable"
        return
    fi
    echo -ne "\033[50D\033[0C backing up $domain"
    tput setaf 8;
    echo -e "\033[50D\033[50C unavailable"
    tput sgr0;
}

print_not_this_repo() {
    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "not backed up to $repo_name"
        return
    fi
    echo -ne "\033[50D\033[0C backing up $domain"
    tput setaf 8;
    echo -e "\033[50D\033[50C not backed up to $repo_name"
    tput sgr0;
}

print_message_to_progress_file() {
    group=$1
    domain=$2
    action=$3
    json=$(cat <<END
{"message_type":"dfb","action":"$action","group":"$group","domain":"$domain"}
END
    )
    echo $json >> /tmp/dfb-progress
}

backup() {
    verify_env

    # options
    gui=false

    for var in "$@"; do
        if [[ "$var" =~ ^-h|--help$  ]]; then
            print_backup_help
        elif [[ "$var" =~ ^--gui$  ]]; then
            gui=true
        fi
    done

    group=$2
    validate_group $group
    repo_name=$3
    validate_repo $group $repo_name
    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")
    domains_directory="$DFB_PATH/$group/domains"

    promt_for_password
    verify_password $password $repo_path

    if [ "$gui" = true ]; then
        touch /tmp/dfb-progress
        tail -f /tmp/dfb-progress | dfb-progress-gui &
    fi

    cd $domains_directory
    find . -type f -print0 |
    while IFS= read -r -d '' domain; do
        domain=$(echo $domain | sed -e 's/^\.\///g')
        backup_domain $password $repo_name $repo_path $domain
        cd $domains_directory
    done

    if [ "$gui" = true ]; then
        rm /tmp/dfb-progress
    fi
    printf "\n"
}

print_backup_help() {
    cat <<HEREDOC
Backup a group of domains.

Usage:
  ${PROGRAM} [group] [repo]

Options:
  --gui         Show progress in a graphical user interface.
  -h --help     Show this screen.
HEREDOC
}
