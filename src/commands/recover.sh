#
# Summary: Recover command
#
# The recover command will allow user to access backed up version
# of domains. A restic repo will be mounted and dfb-fsd daemon will
# make links conveniently availible inside a special directory under
# the domain.

recover() {
    verify_env
    if [[ $2 == "help" ]] || [[ $2 == "" ]]; then
        echo "Usage:"
        echo "  $PROGRAM recover [group] [repo]"
        exit
    fi
    group=$2
    validate_group $group
    repo_name=$3
    validate_repo $group $repo_name
    repo_path=$(cat "$DFB_PATH/$group/repos/$repo_name")

    promt_for_password
    verify_password $password $repo_path

    mountpoint="$DFB_PATH/$group/mountpoint"
    if [ ! -d $mountpoint ]; then
        mkdir $mountpoint
    fi

    if [ "$(df | ggrep $mountpoint)" ]; then
        echo "volume already exists at mountpoint, trying to unmount"
        umount $mountpoint
    fi

    if [ "$(df | ggrep $mountpoint)" ]; then
        tput setaf 1;
        echo "failed to unmount, maybe 'dfb recover' is already running?"
        echo "otherwise, please investigate mountpoint at"
        echo $mountpoint
        tput sgr0;
        exit 1
    else
        echo "successfully unmounted"
    fi

    echo -n "$password" | restic -r $repo_path mount $mountpoint
    printf "\n"

    if [ "$(df | ggrep $mountpoint)" ]; then
        tput setaf 1;
        echo "warning failed to unmount, trying to unmount again"
        tput sgr0;
        umount $mountpoint
    fi

    if [ "$(df | ggrep $mountpoint)" ]; then
        tput setaf 1;
        echo "failed to unmount, please investigate mountpoint at"
        echo $mountpoint
        tput sgr0;
    fi
}
