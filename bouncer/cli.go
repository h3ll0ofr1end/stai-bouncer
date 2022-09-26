package bouncer

import (
	"errors"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	// the arguments handed over to the cli
	args = os.Args
	// function the get the users home directory
	getUserHomeDir = os.UserHomeDir
	// function to enforce the stai executable
	enforceStaiExecutable = enforceExists
	// function to provide file info
	getFileInfo = os.Stat
)

// Context describes the environment of the tool execution
type Context struct {
	// StaiExecutable the stai executable e.g. /home/steffen/stai-blockchain/venv/bin/stai
	StaiExecutable string
	// Location is the location to filter for
	Location string
	// IsDownThreshold is true in case it was set by the user
	IsDownThreshold bool
	// DownThreshold is the down speed threshold to filter for
	DownThreshold float64
	// Done indicates that we are done (--help, --version...)
	Done bool
}

const DefaultStaiExecutableSuffix = "stai-blockchain/venv/bin/stai"

// defaultStaiExecutable to get the default stai executable from the home directory of the current user
func defaultStaiExecutable() (string, error) {
	dirname, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	return dirname + "/" + DefaultStaiExecutableSuffix, nil
}

// enforceExists enforces that the stai executable can be used
func enforceExists(staiExecutable string) error {
	info, err := getFileInfo(staiExecutable)
	if os.IsNotExist(err) {
		return errors.New("stai executable does not exist")
	}
	if info.IsDir() {
		return errors.New("stai executable can not be a directory")
	}

	// TODO could add check if file is executable for current user

	return nil
}

// RunCli starts the cli which includes validation of parameters.
// The returned context consists of a stai executable and location to filter for
func RunCli() (*Context, error) {
	var staiExecutable string
	var location string
	var done bool
	var isDownThreshold bool
	var downThreshold float64

	cli.HelpFlag = &cli.BoolFlag{
		Name:        "help",
		Aliases:     []string{"h"},
		Usage:       "show help",
		Destination: &done,
	}

	app := &cli.App{
		Name:      "stai-bouncer",
		Usage:     "remove unwanted connections from your Stai Node based on Geo IP Location.",
		UsageText: "stai-bouncer [-e CHIA-EXECUTABLE] [-d DOWN-THRESHOLD] LOCATION\n\t stai-bouncer -e /stai-blockchain/venv/bin/stai -d 0.2 mars",
		ArgsUsage: "LOCATION",
		Description: "Tool will lookup connections via 'stai show -c', get ip locations via geoiplookup and " +
			"remove nodes from specified LOCATION via 'stai show -r' ",
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "stai-exec",
				Aliases:     []string{"e"},
				Required:    false,
				DefaultText: "$HOME/stai-blockchain/venv/bin/stai",
				Usage:       "`CHIA-EXECUTABLE`. normally located inside the bin folder of your venv directory",
				Destination: &staiExecutable,
			},
			&cli.Float64Flag{
				Name:        "down-threshold",
				Aliases:     []string{"d"},
				Required:    false,
				DefaultText: "not active",
				Usage:       "`DOWN-THRESHOLD` defines the additional filter for minimal down speed in MiB for filtering.",
				Destination: &downThreshold,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return errors.New("LOCATION is missing")
			}
			if c.IsSet("down-threshold") {
				isDownThreshold = true
			}
			if staiExecutable == "" {
				defaultExecutable, err := defaultStaiExecutable()
				if err != nil {
					return err
				}
				staiExecutable = defaultExecutable
			}
			if err := enforceStaiExecutable(staiExecutable); err != nil {
				return err
			}

			location = strings.TrimSpace(strings.Join(c.Args().Slice(), " "))
			return nil
		},
		Copyright: "GNU GPLv3",
	}

	err := app.Run(args)
	if err != nil {
		return nil, err
	}

	return &Context{
		StaiExecutable:  staiExecutable,
		Location:        location,
		IsDownThreshold: isDownThreshold,
		DownThreshold:   downThreshold,
		Done:            done,
	}, nil
}
