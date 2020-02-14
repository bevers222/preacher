package main

import (
	"bufio"
	"fmt"
	"github.com/bevers222/preacher/internal/commands"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	home, configDirPath, tokenFilePath, orgFilePath, pwd, directory string
)

func init() {
	home = os.Getenv("HOME")
	configDirPath = filepath.Join(home, ".preacher")
	tokenFilePath = filepath.Join(configDirPath, "token")
	orgFilePath = filepath.Join(configDirPath, "org")
	pwd, _ = os.Getwd()
	directory = filepath.Join(pwd, "preacher")
}

func main() {
	app := &cli.App{
		Name:  "preacher",
		Usage: "spread your message",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "org",
				Aliases:     []string{"o"},
				Usage:       "github organization",
				DefaultText: "set from org config file",
				FilePath:    orgFilePath,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "github access token",
				DefaultText: "set from token config file",
				FilePath:    tokenFilePath,
			},
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "interact with files from the `BRANCH` branch",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "show debug logs",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Usage:   "store fetched files in folders in `DIRECTORY` relative to the present working directory",
				Value:   directory,
			},
		},
		CommandNotFound: func(c *cli.Context, command string) {
			fmt.Fprintf(c.App.Writer, "preacher: '%s' is not a preacher command.\n See 'preacher --help'\n", command)
		},
		Commands: []*cli.Command{
			{
				Name:        "fetch",
				Usage:       "fetch all files",
				UsageText:   "fetch [file name] (default: Jenkinsfile)",
				Description: "fetch will pull all occurences of the given file on the given branch",
				ArgsUsage:   "[file name]",
				Action:      commands.Fetch,
			},
		},
		Before: func(c *cli.Context) error {
			var debug = c.Bool("debug")
			_, err := os.Stat(configDirPath)
			if os.IsNotExist(err) {
				scanner := bufio.NewScanner(os.Stdin)
				if debug {
					fmt.Fprintf(c.App.Writer, "Preacher config not found at %v, running first time setup", configDirPath)
				}
				fmt.Fprintf(c.App.Writer, "Welcome to preacher. Let's get a few things configured!\n")
				fmt.Fprintf(c.App.Writer, "Github Access Token: ")
				scanner.Scan()
				token := scanner.Text()
				fmt.Fprintf(c.App.Writer, "Github Organization: ")
				scanner.Scan()
				org := scanner.Text()
				err = os.MkdirAll(configDirPath, os.ModePerm)
				if err != nil {
					fmt.Fprintf(c.App.Writer, "Error making config directory (%v): %v\n", configDirPath, err)
					os.Exit(1)
				}
				err = ioutil.WriteFile(tokenFilePath, []byte(token), 0644)
				if err != nil {
					fmt.Fprintf(c.App.Writer, "Error writing token file (%v): %v\n", tokenFilePath, err)
					os.Exit(1)
				}
				err = ioutil.WriteFile(orgFilePath, []byte(org), 0644)
				if err != nil {
					fmt.Fprintf(c.App.Writer, "Error writing org file (%v): %v\n", orgFilePath, err)
					os.Exit(1)
				}
				fmt.Fprintf(c.App.Writer, "\nPreacher has added config files at %v\nThese values can be overridden at any time by using the --token (-t) and --org (-o) flags.\n\nHappy Preaching!\n\n", configDirPath)
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
