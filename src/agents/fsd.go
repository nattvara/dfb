package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		os.Exit(0)
	}()

	for range time.Tick(2 * time.Second) {
		writeToFile()
	}
}

func writeToFile() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	filename := fmt.Sprintf("%s/Desktop/test.log", usr.HomeDir)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	date := time.Now()

	if _, err := f.Write([]byte(date.String() + "\n")); err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("fsd" + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
