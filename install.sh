#!/usr/bin/env bash

launch_agent="$HOME/Library/LaunchAgents/com.dfb.fsd.plist"
app_path="/Applications/dfb.app"
bins_path="$app_path/Contents/Resources/bin"
symlink_target="/usr/local/bin"

if [ -f $launch_agent ]; then
    launchctl stop $launch_agent
    launchctl unload $launch_agent
fi

if [ -d $app_path ]; then
    echo "removing old version"
    sudo rm -r $app_path
fi

if [ ! -d dfb.app ]; then
    echo "no build of dfb exists at $(pwd)/dfb.app"
    echo "please build dfb before installing."
    exit 1
fi

printf "moving files... "
mv dfb.app /Applications

if [ -f "$symlink_target/dfb" ]; then
    sudo rm "$symlink_target/dfb"
fi

if [ -f "$symlink_target/dfb-progress-parser" ]; then
    sudo rm "$symlink_target/dfb-progress-parser"
fi

if [ -f "$symlink_target/dfb-progress-parser-gui" ]; then
    sudo rm "$symlink_target/dfb-progress-parser-gui"
fi

if [ -f "$symlink_target/dfb-stats" ]; then
    sudo rm "$symlink_target/dfb-stats"
fi

if [ -f "$symlink_target/dfb-fsd" ]; then
    sudo rm "$symlink_target/dfb-fsd"
fi

sudo ln -s "$bins_path/dfb" "$symlink_target/dfb"
sudo ln -s "$bins_path/dfb-progress-parser" "$symlink_target/dfb-progress-parser"
sudo ln -s "$bins_path/dfb-progress-parser-gui" "$symlink_target/dfb-progress-parser-gui"
sudo ln -s "$bins_path/dfb-stats" "$symlink_target/dfb-stats"
sudo ln -s "$bins_path/dfb-fsd" "$symlink_target/dfb-fsd"

if [ ! -d "$HOME/.dfb.logs" ]; then
    echo "creating logs directory at $HOME/.dfb.logs"
    mkdir "$HOME/.dfb.logs"
fi

home=$(echo $HOME | sed 's_/_\\/_g')
sed -e "s/~/$home/g" resources/dfb-fsd.plist > $launch_agent

launchctl load $launch_agent

echo "done."
