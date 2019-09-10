package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"

	"github.com/boilerplate/backend/app/cmd"
	"github.com/boilerplate/backend/app/utils"
)

// Opts with all cli commands and flags
type Opts struct {
	ServerCmd cmd.ServerCommand `command:"server"`

	SharedSecret string `long:"secret" env:"SECRET" required:"true" description:"shared secret key"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "unknown"

func main() {
	utils.InitRandom()

	fmt.Printf("Backend %s\n", revision)
	cmd.Revision = revision

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)

		commonOpts := cmd.CommonOpts{}

		c := command.(cmd.CommonOptionsCommander)

		c.SetCommon(commonOpts)

		err := c.Execute(args)

		if err != nil {
			log.Printf("[ERROR] failed with %+v", err)
		}

		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if dbg {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	log.SetOutput(filter)
}
