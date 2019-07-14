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
            printf "backing up $domain\t unavailible\n"
            return
        fi

        cd $symlink
    else
        cd $domain_path
    fi

    echo "$password" | restic -r $repo_path backup . --tag "$domain" --json | dfb-progress-parser "backing up $domain"
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
}
