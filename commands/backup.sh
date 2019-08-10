#
# Summary: Backup command
#
# The backup command will start a backup of one or more
# domains.

backup() {
    verify_env
    backup_start=$(gdate +%s.%N)

    # options
    gui=false
    confirm=false

    for var in "$@"; do
        if [[ "$var" =~ ^-h|--help$  ]]; then
            print_backup_help
            exit
        elif [[ "$var" =~ ^--gui$  ]]; then
            gui=true
        elif [[ "$var" =~ ^--confirm$  ]]; then
            confirm=true
        fi
    done

    group=$2
    validate_group $group
    repo_name=$3
    validate_repo $group $repo_name
    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")
    domains_directory="$DFB_PATH/$group/domains"
    STATS_PATH="$DFB_PATH/$group/stats"

    if [ "$confirm" = true ]; then
        confirm_backup_should_start $repo_name $group
    fi

    promt_for_password $repo_name
    verify_password $password $repo_path

    if [ "$gui" = true ]; then
        touch /tmp/dfb-progress
        unbuffer tail -f /tmp/dfb-progress | dfb-progress-parser-gui > /dev/stdout 2>&1 &
    fi

    cd $domains_directory
    find . -type f -print0 | sort -z |
    while IFS= read -r -d '' domain; do
        domain=$(echo $domain | sed -e 's/^\.\///g')
        backup_domain $password $repo_name $repo_path $domain
        cd $domains_directory
    done

    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$repo_name" "gathering_stats"
    else
        printf "\n\n  gathering stats for repo $repo_name... "
    fi

    repo_raw_data_csv="$STATS_PATH/repo_raw_data.csv"
    repo_time_took_csv="$STATS_PATH/repo_time_took.csv"

    echo -n "$password" \
    | restic -r "$repo_path" stats --mode raw-data --json \
    | ggrep "{" \
    | jq -r '[.[]] | @csv' \
    | tr -d '\n' >> "$repo_raw_data_csv" \
    && echo ",$group,$repo_name,$(gdate +%Y-%m-%dT%H:%M:%S%z)" >> "$repo_raw_data_csv"

    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "repo" "gathering_stats_done"
        print_message_to_progress_file "$group" "null" "done"
    else
        echo "done"
    fi

    if [ "$gui" = true ]; then
        ps aux | ggrep "[t]ail -f /tmp/dfb-progress" | awk '{print $2}' | xargs kill -9
        rm /tmp/dfb-progress
    fi

    backup_done=$(gdate +%s.%N)
    backup_took=$(echo "$backup_done - $backup_start" | bc)
    echo "$backup_took,$group,$repo_name,$(gdate +%Y-%m-%dT%H:%M:%S%z)" >> $repo_time_took_csv

    terminal-notifier -group "dfb" -title "dfb" -subtitle "Done" -message "Backup of $group to $repo_name is done" -sender "com.example.dfb" -sound default > /dev/null

    printf "\n"
}

print_backup_help() {
    cat <<HEREDOC
Backup a group of domains.

Usage:
  ${PROGRAM} [group] [repo]

Options:
  --gui         Show progress in a graphical user interface.
  --confirm     Show a dialogue that the user have to confirm for backup to start. Useful if backup is started by a cron job.
  -h --help     Show this screen.
HEREDOC
}

confirm_backup_should_start() {
    repo_name=$1
    group=$2
    script_result=$(osascript -s o <<END
display dialog "A backup is about to start \n\ngroup: \t$group\nrepo: \t$repo_name" buttons {"Proceed", "Abort"} default button "Proceed" cancel button "Abort" with title "dfb"
END
    )
    if [ "$script_result" != "button returned:Proceed" ]; then
        terminal-notifier -group "dfb" -title "dfb" -subtitle "Aborted" -message "Backup of $group was aborted" -sender "com.example.dfb" > /dev/null
        exit 1
    fi
}

backup_domain() {
    domain_path=$(cat "./$domain" | ggrep -E 'path:' | egrep -o '[^:]+$' | tr -d '[:space:]')
    symlink=$(cat "./$domain" | ggrep -E 'symlink:' | egrep -o '[^:]+$' | tr -d '[:space:]')
    exclusions=$(cat "./$domain" | ggrep -E 'exclusions:' | egrep -o '[^:]+$' | tr " " "\n")
    repos=$(cat "./$domain" | ggrep -E 'repos:' | egrep -o '[^:]+$' | tr -d '[:space:]')

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

    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "begin"
    fi

    snapshots_csv="$STATS_PATH/snapshots.csv"
    domain_restore_size_csv="$STATS_PATH/domain_restore_size.csv"
    domain_raw_data_csv="$STATS_PATH/domain_raw_data.csv"
    if [ ! -d "$STATS_PATH" ]; then
        mkdir "$STATS_PATH"
    fi

    echo "$exclusions" > /tmp/dfb_exclusions

    echo -n "$password" \
        | restic -r "$repo_path" \
        backup "$restic_backup_path" \
        --tag "$domain" \
        --exclude-file /tmp/dfb_exclusions \
        --verbose \
        --json \
        2>&1 \
        | unbuffer -p tee >( \
            ggrep "summary" \
            | jq -r 'select(.message_type=="summary") | [.[]] | @csv' \
            | tr -d '\n' >> "$snapshots_csv" \
            && echo ",$group,$domain,$repo_name,$(gdate +%Y-%m-%dT%H:%M:%S%z)" >> "$snapshots_csv"
        ) \
        | if [ "$gui" = true ]; \
            then \
                unbuffer -p ggrep "" >> /tmp/dfb-progress; \
            else \
                dfb-progress-parser "$group" "$domain"; \
        fi

    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "gathering_stats"
    fi

    echo -n "$password" \
        | restic -r "$repo_path" stats latest --mode restore-size --json \
        | ggrep "{" \
        | jq -r '[.[]] | @csv' \
        | tr -d '\n' >> "$domain_restore_size_csv" \
        && echo ",$group,$domain,$repo_name,$(gdate +%Y-%m-%dT%H:%M:%S%z)" >> "$domain_restore_size_csv"

    echo -n "$password" \
        | restic -r "$repo_path" stats latest --mode raw-data --json \
        | ggrep "{" \
        | jq -r '[.[]] | @csv' \
        | tr -d '\n' >> "$domain_raw_data_csv" \
        && echo ",$group,$domain,$repo_name,$(gdate +%Y-%m-%dT%H:%M:%S%z)" >> "$domain_raw_data_csv"

    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "gathering_stats_done"
    else
        echo "done."
    fi

    printf "\r"
}

print_domain_unavailable() {
    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "unavailable"
        return
    fi
    printf "\033[50D\033[0C backing up $domain "
    tput setaf 8;
    printf "\033[50D\033[60Cunavailable \n"
    tput sgr0;
}

print_not_this_repo() {
    if [ "$gui" = true ]; then
        print_message_to_progress_file "$group" "$domain" "not_this_repo"
        return
    fi
    printf "\033[50D\033[0C backing up $domain "
    tput setaf 8;
    printf "\033[50D\033[60Cnot backed up to $repo_name \n"
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
    printf "$json\n\n" >> /tmp/dfb-progress
}
