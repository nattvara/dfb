#
# Summary: lock functions
#
# The lock is used to prevent the user from running dfb commands
# that should not run simultaneously.

LOCK_FILE="$HOME/.dfb.lock"

check_lock() {
    if [ -f "$LOCK_FILE" ]; then
        printf "ERR: lock file found at $LOCK_FILE\n"
        printf "\"$(cat $LOCK_FILE)\""
        printf "\n\ndfb is locked. this is probably because it is already running.\n"
        printf "if you are sure you want to proceed re-run this command with the --force flag\n"
        exit 1
    fi
}

lock_dfb() {
    locked_by=$1
    printf "locked at $(gdate) by $locked_by command" > "$LOCK_FILE"
}

unlock_dfb() {
    rm "$LOCK_FILE"
}
