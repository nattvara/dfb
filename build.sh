#!/usr/bin/env bash

BUILD_DIR="$(pwd)/build"
OUT="$BUILD_DIR/dfb"

if [ ! -d $BUILD_DIR ]; then
    mkdir $BUILD_DIR
fi

if [ -f $OUT ]; then
    rm $OUT
fi

printf "creating dfb script... "

touch $OUT
echo "#!/usr/bin/env bash" > $OUT
cat main.sh >> $OUT && printf "\n" >> $OUT
cat commands/groups.sh >> $OUT && printf "\n" >> $OUT
cat commands/domains.sh >> $OUT && printf "\n" >> $OUT
cat commands/backup.sh >> $OUT && printf "\n" >> $OUT
cat commands/recover.sh >> $OUT && printf "\n" >> $OUT
cat commands/stats.sh >> $OUT && printf "\n" >> $OUT
cat commands/fsd.sh >> $OUT && printf "\n" >> $OUT
cat helpers/password.sh >> $OUT && printf "\n" >> $OUT
cat helpers/validation.sh >> $OUT && printf "\n" >> $OUT
cat helpers/lock.sh >> $OUT && printf "\n" >> $OUT

echo 'main "$@"' >> $OUT

chmod +x $OUT

echo "done."
printf "compiling binaries... "

go build -o ./build/dfb-progress-parser -i ./tools/progress-parser/cmd.go
FYNE_FONT=/Applications/dfb.app/Contents/Resources/fonts/Lato-Black.ttf go build -o ./build/dfb-progress-parser-gui -i ./tools/progress-parser-gui/cmd.go
go build -o ./build/dfb-stats -i ./tools/stats/cmd.go
go build -o ./build/dfb-fsd -i ./agents/fsd.go

echo "done."
printf "packaging application... "

fyne package -executable build/dfb-progress-parser-gui -icon resources/icon.png -name dfb

mkdir dfb.app/Contents/Resources/fonts
cp resources/fonts/*.ttf dfb.app/Contents/Resources/fonts

mkdir dfb.app/Contents/Resources/bin
cp build/dfb dfb.app/Contents/Resources/bin/dfb
cp build/dfb-progress-parser dfb.app/Contents/Resources/bin/dfb-progress-parser
cp build/dfb-stats dfb.app/Contents/Resources/bin/dfb-stats
cp build/dfb-fsd dfb.app/Contents/Resources/bin/dfb-fsd
echo "FYNE_SCALE=0.9 FYNE_FONT=/Applications/dfb.app/Contents/Resources/fonts/Lato-Black.ttf /Applications/dfb.app/Contents/MacOS/dfb-progress-parser-gui" > dfb.app/Contents/Resources/bin/dfb-progress-parser-gui
chmod +x dfb.app/Contents/Resources/bin/dfb-progress-parser-gui

echo "done."
