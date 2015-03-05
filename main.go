package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

var (
	// Waitgroup for all go routines
	wg sync.WaitGroup

	// Database object
	db *Database

	// Debug/Verbose switch
	Debug bool

	// Tray icon switch
	ShowIcon bool

	// Generate sameple configuration switch
	GenConfig bool

	// Configuration filepath
	ConfigFile string

	// Cache filepath
	DatabaseFile string
)

func init() {

	flag.BoolVar(&Debug, "v", false, "debugging/verbose information")
	flag.BoolVar(&ShowIcon, "i", false, "Show a icon in the statusbar/tray")
	flag.BoolVar(&GenConfig, "g", false, "Generate a configuration file")
	flag.StringVar(&ConfigFile, "f", "", "The configuration file in which the user settings are stored")
	flag.StringVar(&DatabaseFile, "d", "", "The directory in which the database will be stored")
	flag.Parse()

	if GenConfig {

		cfg, err := GenerateConfig()
		if err != nil {

			log.Fatal(err)
		}

		_, err = os.Stdout.Write(cfg)
		if err != nil {

			log.Fatal(err)
		}

		os.Exit(0)
	}

	if len(ConfigFile) == 0 {

		// Search for config directory
		if len(os.Getenv("XDG_CONFIG_HOME")) != 0 {

			ConfigFile = os.ExpandEnv("$XDG_CONFIG_HOME/mnd.conf")

		} else {

			ConfigFile = os.ExpandEnv("$HOME/.config/mnd.conf")
		}
	}

	if len(DatabaseFile) == 0 {

		// Search for data directory
		if len(os.Getenv("XDG_DATA_HOME")) != 0 {

			DatabaseFile = os.ExpandEnv("$XDG_DATA_HOME/mnd.db")

		} else {

			DatabaseFile = os.ExpandEnv("$HOME/.local/share/mnd.db")
		}
	}
}

func main() {

	d, err := InitDB(DatabaseFile)
	if err != nil {

		log.Fatal(err)
	}
	db = d

	cfg, err := GetCFG()
	if err != nil {

		log.Fatal(err)
	}

	for i, _ := range cfg.Accounts {

		wg.Add(1)
		go cfg.LaunchPOP3Client(&wg, &cfg.Accounts[i])
	}

	if ShowIcon {

		ShowStatusIcon()
	}

	wg.Wait()
}
