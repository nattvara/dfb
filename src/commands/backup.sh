#
# Summary: Backup command
#
# The backup command will start a backup of one or more
# domains.

backup_domain() {
    password="$1"
    repo_path="$2"
    domain="$3"
    domain_path=$(cat "./$domain" | ggrep -E 'path' | egrep -o '[^:]+$' | tr -d '[:space:]')
    symlink=$(cat "./$domain" | ggrep -E 'symlink' | egrep -o '[^:]+$' | tr -d '[:space:]')
    exclusions=$(cat "./$domain" | ggrep -E 'exclusions' | egrep -o '[^:]+$' | tr " " "\n")

    if [[ $symlink != "" ]]; then
        if [ ! -d $symlink ]; then
            print_domain_unavailible $domain
            return
        fi

        cd $symlink
    elif [ -f $domain_path ]; then
        parent_dir=$(dirname $domain_path)
        if [ ! -d $parent_dir ]; then
            print_domain_unavailible $domain
            return
        fi
        cd $parent_dir
    else
        if [ ! -d $domain_path ]; then
            print_domain_unavailible $domain
            return
        fi
        cd $domain_path
    fi

    if [ -f $domain_path ]; then
        restic_backup_path=$(basename $domain_path)
    else
        restic_backup_path="."
    fi

    echo "$exclusions" > /tmp/dfb_exclusions
    echo -n "$password" | restic -r $repo_path backup "$restic_backup_path" --tag "$domain"  --exclude-file /tmp/dfb_exclusions --json | dfb-progress-parser "  backing up $domain"
    printf "\r"
}

print_domain_unavailible() {
    echo -ne "\033[50D\033[0C backing up $domain"
    tput setaf 8;
    echo -e "\033[50D\033[50C unavailible"
    tput sgr0;
}

backup() {
    verify_env
    if [[ $2 == "help" ]] || [[ $2 == "" ]]; then
        echo "Usage:"
        echo "  $PROGRAM backup [group] [repo] [<timestamp-file>]"
        exit
    fi
    group=$2
    validate_group $group
    repo_name=$3
    validate_repo $group $repo_name
    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")
    domains_directory="$DFB_PATH/$group/domains"

    promt_for_password
    verify_password $password $repo_path

    cd $domains_directory
    find . -type f -print0 |
    while IFS= read -r -d '' domain; do
        domain=$(echo $domain | sed -e 's/^\.\///g')
        backup_domain $password $repo_path $domain
        cd $domains_directory
    done

    printf "\n"
}
