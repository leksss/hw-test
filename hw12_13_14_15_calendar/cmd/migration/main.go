package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

/*
@Usage
$ go run . -dir=./migrations mysql "user:1234@/calendar?parseTime=true" status
$ go run . -dir=./migrations mysql "user:1234@/calendar?parseTime=true" create init sql
$ go run . -dir=./migrations mysql "user:1234@/calendar?parseTime=true" up
*/

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]

	db, err := goose.OpenDBWithDriver("mysql", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err) //nolint
	}
}
