#!/usr/bin/env bash

OUT="./build/dfb"

if [ -f $OUT ]; then
    rm $OUT
fi

echo "#!/usr/bin/env bash" > $OUT
cat src/main.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/groups.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/domains.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/backup.sh >> $OUT && printf "\n" >> $OUT
cat src/commands/fsd.sh >> $OUT && printf "\n" >> $OUT
cat src/helpers/password.sh >> $OUT && printf "\n" >> $OUT
cat src/helpers/validation.sh >> $OUT && printf "\n" >> $OUT

echo 'main "$@"' >> $OUT

chmod +x $OUT

go build -o ./build/dfb-progress-parser -i ./src/tools/progress-parser.go
go build -o ./build/dfb-fsd -i ./src/agents/fsd.go
