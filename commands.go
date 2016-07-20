//
// commands.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package main

import (
	"fmt"
	"os"

	"github.com/jkawamoto/fgo/command"
	"github.com/urfave/cli"
)

// GlobalFlags defines global flags.
var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "dest, d",
		Value: "pkg",
		Usage: "directory `NAME` to store package files",
	},
	cli.StringFlag{
		Name:  "brew, b",
		Value: "brew",
		Usage: "directory `NAME` to store brew file",
	},
}

// Commands defines sub commands.
var Commands = []cli.Command{
	{
		Name:      "init",
		Usage:     "create Makefile and other related directories.",
		ArgsUsage: "[username]",
		Action:    command.CmdInit,
	},
	{
		Name:      "build",
		Usage:     "build binaries, upload them, an update brew formula.",
		ArgsUsage: "[version]",
		Action:    command.CmdBuild,
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}