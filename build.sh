#!/usr/bin/env bash

BUILD_DIR="$(pwd)/build"
OUT="$BUILD_DIR/dfb"

if [ ! -d $BUILD_DIR ]; then
    mkdir $BUILD_DIR
fi

if [ -f $OUT ]; then
    rm $OUT
fi

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

echo 'main "$@"' >> $OUT

chmod +x $OUT

go build -o ./build/dfb-progress-parser -i ./tools/progress-parser/cmd.go
go build -o ./build/dfb-progress-parser-gui -i ./tools/progress-parser-gui/cmd.go
go build -o ./build/dfb-stats -i ./tools/stats/cmd.go
go build -o ./build/dfb-fsd -i ./agents/fsd.go
