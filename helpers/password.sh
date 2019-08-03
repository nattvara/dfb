#
# Summary: password functions
#

promt_for_password() {
    password=$(osascript <<END
set x to display dialog "What is your password?" default answer "" with hidden answer
set y to (text returned of x)
END
    )
    if [ -z "$(echo ${password//[[:blank:]]/})" ]; then
        osascript -e "display notification with title \"Password cannot be empty\""
        exit 1
    fi
}

verify_password() {
    password="$1"
    repo_path="$2"

    if echo "$password" | restic -r "$repo_path" key list > /dev/null; then
        echo "valid password"
    else
        osascript -e "display notification \"for $repo_path\" with title \"Invalid password\""
        exit 1
    fi
}
