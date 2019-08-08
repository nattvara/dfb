#
# Summary: password functions
#

promt_for_password() {
    repo_name=$1
    password=$(osascript -s o <<END
set x to display dialog "Enter the password for the repo $repo_name" default answer "" with hidden answer
set y to (text returned of x)
END
    )
    if [ -z "$(echo ${password//[[:blank:]]/})" ]; then
        terminal-notifier -group "dfb" -title "dfb" -subtitle "Backup" -message "password cannot be empty" -sender "com.example.dfb" > /dev/null
        exit 1
    fi
    if [[ "$password" =~ "User cancelled" ]]; then
        terminal-notifier -group "dfb" -title "dfb" -subtitle "Backup" -message "no password entered" -sender "com.example.dfb" > /dev/null
        exit
    fi
}

verify_password() {
    password="$1"
    repo_path="$2"

    if echo "$password" | restic -r "$repo_path" key list 2> /dev/null 1> /dev/null; then
        terminal-notifier -group "dfb" -title "dfb" -subtitle "Backup" -message "Password correct" -sender "com.example.dfb" > /dev/null
    else
        terminal-notifier -group "dfb" -title "dfb" -subtitle "Backup" -message "Invalid password for $repo_path" -sender "com.example.dfb" > /dev/null
        exit 1
    fi
}
