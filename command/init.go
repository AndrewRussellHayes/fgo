//
// command/init.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/deiwin/interact"
	"github.com/jkawamoto/fgo/fgo"
	colorable "github.com/mattn/go-colorable"
	"github.com/ttacon/chalk"
	"github.com/urfave/cli"

	"github.com/tcnksm/go-gitconfig"
)

// InitOpt defines options for cmdInit.
type InitOpt struct {
	// Configuration
	Config Config
	// GitHub user name.
	UserName string
	// GitHub repository name.
	Repository string
	// Description of the target application.
	Description string
}

// Generater is an interface provides Generate method.
type Generater interface {
	Generate() ([]byte, error)
}

// CmdInit parses options and run cmdInit.
func CmdInit(c *cli.Context) error {

	opt := InitOpt{
		Config: Config{
			Package:  c.GlobalString(PackageFlag),
			Homebrew: c.GlobalString(HomebrewFlag),
		},
		UserName:    c.Args().First(),
		Repository:  c.Args().Get(1),
		Description: c.String("desc"),
	}
	return cmdInit(&opt)

}

// cmdInit defines init command action.
func cmdInit(opt *InitOpt) (err error) {

	stdout := colorable.NewColorableStdout()

	// Check user name.
	if opt.UserName == "" {
		fmt.Fprintf(stdout, "Checking git configuration to get the user name: ")
		opt.UserName, err = gitconfig.GithubUser()
		if err != nil {
			fmt.Fprintf(stdout, chalk.Red.Color("Cannot find user name (%v)"), err.Error())
			return cli.NewExitError("", 1)
		}
		fmt.Fprintln(stdout, chalk.Yellow.Color(opt.UserName))
	}

	// Prepare directories.
	fmt.Fprintf(stdout, "Preparing the directory to store a brew formula: ")
	if err = prepareDirectory(opt.Config.Homebrew); err != nil {
		fmt.Fprintf(stdout, chalk.Red.Color("failed (%v)"), err.Error())
		return cli.NewExitError("", 2)
	}
	fmt.Fprintln(stdout, chalk.Green.Color("done"))

	// Check Makefile doesn't exist and create it.
	actor := interact.NewActor(os.Stdin, stdout)

	createMakefile := true
	if _, exist := os.Stat("Makefile"); exist == nil {
		createMakefile, err = actor.Confirm("Makefile already exists. Would you like to overwrite it?", interact.ConfirmDefaultToNo)
		if err != nil {
			fmt.Fprintf(stdout, chalk.Red.Color("failed (%v)"), err.Error())
			return cli.NewExitError("", 3)
		}
	}
	if createMakefile {
		fmt.Fprintf(stdout, "Creating Makefile: ")
		err = createResource("Makefile", &fgo.Makefile{
			Dest:     opt.Config.Package,
			UserName: opt.UserName,
		})
		if err != nil {
			fmt.Fprintf(stdout, chalk.Yellow.Color("skipped (%s)\n"), err.Error())
		} else {
			fmt.Fprintln(stdout, chalk.Green.Color("done"))
		}
	} else {
		fmt.Fprintln(stdout, "Creating Makefile:", chalk.Yellow.Color("skipped"))
	}

	// Check brew rb file doesn't exist and create it.
	if opt.Repository == "" {
		fmt.Fprintf(stdout, "Checking git configuration to get the repository name: ")
		opt.Repository, err = gitconfig.Repository()
		if err != nil {
			fmt.Fprintf(stdout, chalk.Red.Color("skipped (%s).\n"), err.Error())
			fmt.Fprintln(stdout, chalk.Yellow.Color("You must re-run init command after setting a remote repository"))
		}
		fmt.Fprintln(stdout, chalk.Yellow.Color(opt.Repository))
	}
	if opt.Repository != "" {
		tmpfile := filepath.Join(opt.Config.Homebrew, fmt.Sprintf("%s.rb.template", opt.Repository))

		createTemplate := true
		if _, exist := os.Stat(tmpfile); exist == nil {
			createTemplate, err = actor.Confirm("brew formula template already exists. Would you like to overwrite it?", interact.ConfirmDefaultToNo)
			if err != nil {
				fmt.Fprintf(stdout, chalk.Red.Color("failed (%v)"), err.Error())
				return cli.NewExitError("", 3)
			}
		}
		if createTemplate {
			fmt.Fprintf(stdout, "Creating brew formula template: ")
			err = createResource(tmpfile, &fgo.FormulaTemplate{
				Package:     opt.Repository,
				UserName:    opt.UserName,
				Description: opt.Description,
			})
			if err != nil {
				fmt.Fprintf(stdout, chalk.Yellow.Color("skipped (%s).\n"), err.Error())
			} else {
				fmt.Fprintln(stdout, chalk.Green.Color("done"))
			}
		} else {
			fmt.Fprintln(stdout, "Creating brew formula template:", chalk.Yellow.Color("skipped"))
		}
	}

	return

}

// prepareDirectory creates a directory if necessary.
func prepareDirectory(path string) error {

	if info, exist := os.Stat(path); exist == nil && !info.IsDir() {
		return fmt.Errorf("cannot make directory %s", path)
	} else if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("cannot make directory %s (%s)", path, err.Error())
	}
	return nil

}

// createResource creates a resource from a given generator and stores it to
// a given path.
func createResource(path string, data Generater) (err error) {

	buf, err := data.Generate()
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, buf, 0644)

}
