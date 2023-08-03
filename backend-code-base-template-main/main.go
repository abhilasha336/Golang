package main

import (
	"backend-code-base-template/app"
	"backend-code-base-template/migrations"
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	runserver = flag.Bool("runserver", false, "This is a string argument for running server")
	migration = flag.Bool("migration", false, "This is a string argument for running migration")
	up        = flag.Bool("up", false, "This is a string argument for running migration up")
	down      = flag.Bool("down", false, "This is a string argument for running migration down")
)

func main() {
	// logrus init
	log := logrus.New()

	flag.Parse()

	if !*runserver && !*migration {
		log.Fatalf("Please specify the file you want to execute")
	}
	if *runserver == true {
		app.Run()
	}
	if *migration == true {
		if !*down && !*up {
			log.Fatal("Please specify the migration type")
		}
		if *up == true {
			migrations.Migration("up")
		}
		if *down == true {
			migrations.Migration("down")
		}
	}

}
