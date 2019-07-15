#!/usr/bin/env bash

launch_agent="$HOME/Library/LaunchAgents/com.dfb.fsd.plist"

if [ -f $launch_agent ]; then
    launchctl stop $launch_agent
    launchctl unload $launch_agent
fi

cp build/dfb /usr/local/bin/dfb
cp build/dfb-progress-parser /usr/local/bin/dfb-progress-parser
cp build/dfb-fsd /usr/local/bin/dfb-fsd
cp resources/dfb-fsd.plist $launch_agent

launchctl load $launch_agent
