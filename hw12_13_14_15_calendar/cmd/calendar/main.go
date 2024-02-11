package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mmart-pro/otusGo/hw12_13_14_15_calendar/internal/app"
	flag "github.com/spf13/pflag"
)

func main() {
	configFlag := flag.StringP("config", "c", "", "json config")
	version := flag.BoolP("version", "v", false, "print app version")
	help := flag.BoolP("help", "h", false, "usage help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	if *version {
		printVersion()
		return
	}
	if *configFlag == "" {
		fmt.Println("config flag required")
		flag.Usage()
		return
	}

	err := app.NewCalendar(*configFlag).
		Startup(context.Background())
	if err != nil {
		log.Fatal()
	}
}
