#
# Summary: Recover domain command
#
# The recover-domain command restores the latest version
# of a single domain.

recover_domain() {
    verify_env

    # options
    force=false

    for var in "$@"; do
        if [[ "$var" =~ ^-h|--help$  ]]; then
            print_recover_domain_help
            exit
        elif [[ "$var" =~ ^--force$  ]]; then
            force=true
        fi
    done

    group=$2
    validate_group $group
    domain=$3
    validate_domain $group $domain
    repo_name=$4
    validate_repo $group $repo_name

    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")

    echo "$group"
    echo "$domain"
    echo "$repo_path"
    promt_for_password
    verify_password $password $repo_path

    echo -n "$password" | restic -r $repo_path restore --tag $domain latest --target "$HOME/$domain"
    printf "\n"
}

print_recover_domain_help() {
    cat <<HEREDOC
Recover files from a domain in a group.

The recover-domain command restores the latest version
of a single domain.

Usage:
  ${PROGRAM} recover [group] [domain] [repo]

Options:
  -h --help     Show this screen.
HEREDOC
}
