#
# Summary: fsd command
#
# The fsd command will allow user to start and stop the
# dfb-fsd daemon.

launch_agent="$HOME/Library/LaunchAgents/com.dfb.fsd.plist"

fsd() {
    verify_env
    if [[ "${2:-}" =~ ^-h|--help$  ]]
    then
        print_fsd_help
    elif [ "${2:-}" == "start" ]
    then
        fsd_start
    elif [ "${2:-}" == "stop" ]
    then
        fsd_stop
    else
        print_fsd_help
    fi
}

print_fsd_help() {
    cat <<HEREDOC
dfb filesystem agent.

The dfb-fsd agent watches domains and the mount point
for the restic filesystem, updating symlinks to show
backed up versions of domain, as well as cleaning up
on unmount.

Usage:
  ${PROGRAM} fsd <subcommand>

Available Commands:
  start       Start the agent.
  stop        Stop the agent.

Options:
  -h --help     Show this screen.
HEREDOC
}

fsd_start() {
    echo "starting agent"
    launchctl load $launch_agent
}

fsd_stop() {
    echo "stopping agent"
    if [ $(launchctl unload $launch_agent) ]; then
        echo "agent stopped"
    else
        echo "agent wasn't running"
    fi
}
