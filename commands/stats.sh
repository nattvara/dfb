#
# Summary: stats command
#
# The stats command will allow a user to view metrics about
# the backed up data. This command is just a wrapper
# around the tool written in go (see tools/stats) but
# will try to open the output in quick look. Note, if a
# non default output path is used no preview will be availible.

stats() {
    verify_env
    default_output_path="/tmp/dfb-metric.png"

    dfb-stats "${@:2}"
    if [ -f $default_output_path ]; then
        qlmanage -p $default_output_path 2> /dev/null 1> /dev/null
        rm $default_output_path
    fi
}
