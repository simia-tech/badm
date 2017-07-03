package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/simia-tech/badm"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configurationPath := filepath.Join(user.HomeDir, ".badm")

	if err := badm.RegisterTypes(configurationPath); err != nil {
		log.Fatal(err)
	}

	typeNames := badm.TypeNames()

	app := &cli.App{
		Name:    "badm",
		Usage:   "administration of bolt db files",
		Version: "0.1",
		Commands: []*cli.Command{
			{
				Name:    "select",
				Aliases: []string{"s"},
				Usage:   "select a bolt db file",
				Action: func(c *cli.Context) error {
					return badm.Select(configurationPath, c.Args().First())
				},
			},
			{
				Name:    "buckets",
				Aliases: []string{"b"},
				Usage:   "lists buckets in the database",
				Action: func(c *cli.Context) error {
					return badm.ListBuckets(configurationPath)
				},
			},
			{
				Name:    "keys",
				Aliases: []string{"k"},
				Usage:   "lists keys in the bucket",
				Action: func(c *cli.Context) error {
					return badm.ListKeys(configurationPath, c.Args().First())
				},
			},
			{
				Name:    "values",
				Aliases: []string{"v"},
				Usage:   "lists values in the bucket",
				Action: func(c *cli.Context) error {
					return badm.ListValues(configurationPath, c.Args().First())
				},
			},
			{
				Name:    "key-values",
				Aliases: []string{"kv"},
				Usage:   "lists keys and values in the bucket",
				Action: func(c *cli.Context) error {
					return badm.ListKeyValues(configurationPath, c.Args().First())
				},
			},
			{
				Name:  "set",
				Usage: "set different values",
				Subcommands: []*cli.Command{
					{
						Name:    "key-type",
						Aliases: []string{"kt"},
						Usage:   fmt.Sprintf("sets the key type for a bucket (possible types are: %s)", strings.Join(typeNames, ", ")),
						Action: func(c *cli.Context) error {
							return badm.SetKeyType(configurationPath, c.Args().Get(0), c.Args().Get(1))
						},
					},
					{
						Name:    "value-type",
						Aliases: []string{"vt"},
						Usage:   fmt.Sprintf("sets the value type for a bucket (possible types are: %s)", strings.Join(typeNames, ", ")),
						Action: func(c *cli.Context) error {
							return badm.SetValueType(configurationPath, c.Args().Get(0), c.Args().Get(1))
						},
					},
				},
			},
			{
				Name:    "clear",
				Aliases: []string{"c"},
				Usage:   "clears a bucket configuration",
				Action: func(c *cli.Context) error {
					return badm.Clear(configurationPath, c.Args().First())
				},
			},
			{
				Name:    "plugin",
				Aliases: []string{"p"},
				Usage:   "manage the plugins",
				Subcommands: []*cli.Command{
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "lists all plugins",
						Action: func(c *cli.Context) error {
							return badm.ListPlugins(configurationPath)
						},
					},
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "add a plugin",
						Action: func(c *cli.Context) error {
							return badm.AddPlugin(configurationPath, c.Args().First())
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
