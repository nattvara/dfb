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
cat src/main.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/groups.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/domains.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/backup.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/recover.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/fsd.sh >> $OUT && printf "\n" >> $OUT
cat src/helpers/password.sh >> $OUT && printf "\n" >> $OUT
cat src/helpers/validation.sh >> $OUT && printf "\n" >> $OUT

echo 'main "$@"' >> $OUT

chmod +x $OUT

go build -o ./build/dfb-progress-parser -i ./src/tools/progress-parser.go
go build -o ./build/dfb-progress-parser-gui -i ./src/tools/progress-parser-gui.go
go build -o ./build/dfb-fsd -i ./src/agents/fsd.go
