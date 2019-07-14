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

    if [[ $symlink != "" ]]; then
        if [ ! -d $symlink ]; then
            echo -ne "\033[50D\033[0C backing up $domain"
            tput setaf 8;
            echo -e "\033[50D\033[50C unavailible"
            tput sgr0;
            return
        fi

        cd $symlink
    elif [ -f $domain_path ]; then
        parent_dir=$(dirname $domain_path)
        cd $parent_dir
    else
        cd $domain_path
    fi

    if [ -f $domain_path ]; then
        restic_backup_path=$(basename $domain_path)
    else
        restic_backup_path="."
    fi

    echo -n "$password" | restic -r $repo_path backup "$restic_backup_path" --tag "$domain" --json | dfb-progress-parser "  backing up $domain"
    printf "\r"
}

backup() {
    verify_env
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
        backup_domain $password $repo_path $domain
        cd $domains_directory
    done

    printf "\n"
}
