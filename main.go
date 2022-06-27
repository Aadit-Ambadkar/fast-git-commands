package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type Preferances struct {
	useHTTPS bool   // If true, uses HTTPS, else uses SSH
	username string // Users Name

}

var (
	prefs    = Preferances{useHTTPS: true, username: ""}
	version  = "dev"
	helpText = "\x1b[32mFIT - Fast Git!\x1b[0m\n\n\x1b[4m\x1b[36mOptions:\x1b[0m\n    clone   Clone a Repository with Default Owner of Specified UserName\n            Requires Additional Argument of Repository Name (Overriden by --repo/-r)\n\n        *Secondary Commands*\n        --repo/-r       Specify Repository URL or SSH;\n                        Requires Additional Argument of URL or SSH\n        --user/-u       Sepcify Owner of Repository\n                        Requires Additional Argument of Owner Username\n\n    branch  Switch to Specified Branch (Overriden by any flag)\n            Requires Additional Argument of Branch Name (Overriden by --list/-l)\n\n        *Secondary Commands*\n        --list/-l       List all Branches\n        --del/-d        Delete a Branch\n                        Requires Additional Argument of Branch Name\n        --new/-n        Create New Branch; Also sets upstream\n                        Requires Additional Argument of Branch Name\n\n    push    Adds all changes, Commits With Argument Message, Pushes\n            Requires Additional Arugment of Commit Message (Overriden by --none/-n)\n\n        *Secondary Commands*\n        --none/-n       Do not Add and Commit, Just Push\n\n    pull    Pull From Upstream\n\n    set     Sets Different Defaults.\n            Requires Additional Flags of What to Set\n\n        *Secondary Commands*\n        --name/-n       Sets Git Username\n                        Rquires Additional Argument of Username\n        --email/-e      Sets Git Email\n                        Requires Additional Argument of Email.\n        --https/-h      Sets FIT to use HTTPS\n        --ssh/-s        Sets FIT to use SSH\n\n    raw     Run Argument as Git Command\n            Requires Additional Argument of Command\n\n\x1b[4m\x1b[36mExamples:\x1b[0m\n    fit clone fast-git-commands              # Clones this repo (My username is Aadit-Ambadkar)\n    fit clone -r git@gith ... nds.git        # Clones this repo with SSH\n    fit clone linux --user torvalds          # Clones linux\n    fit branch main                          # Switches to the main branch\n    fit branch -n dev                        # Switches to a new branch, dev\n    fit branch --list                        # Lists all the branches\n    fit push \"My First Commit\"             # Adds, Commits, and Pushes all changes\n    fit push --none                          # Pushes all commits\n    fit pull                                 # Pulls from Upstream\n    fit set --name Aadit-Ambadkar            # Sets name to Aadit Ambadkar\n    fit set -s                               # Sets FIT to Use SSH\n    fit raw commit -a -m \"Github\"          # Commits all changes under message \"Github\"\n    fit raw rebase main                      # Rebases with main\n"
)

func main() {
	if ok, _ := ArgsHaveOption("help", "h"); ok {
		fmt.Println(helpText)
		return
	}

	if ok, _ := ArgsHaveOption("version", "v"); ok {
		fmt.Println("FIT " + version)
		return
	}

	prefsFile, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.MkdirAll(filepath.Join(prefsFile, ".config", "fastgit"), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	prefsFile = filepath.Join(prefsFile, ".config", "fastgit", "data.json")

	prefs, err = readFromFileAsJSON(prefsFile)

	retCode := 0 // used to exit with non-zero code later on if needed
	defer func() {
		err = writeToFileAsJSON(prefs, prefsFile)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(retCode)
	}()

	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("fit needs a command argument, run \x1b[33mfit --help\x1b[0m for more information")
	}

	command := os.Args[1]

	if command == "clone" {
		if ok, i := ArgsHaveOption("repo", "r"); ok {
			if len(os.Args) < i+2 {
				fmt.Println("option --repo requires an argument")
				return
			}
			arg := os.Args[i+1]
			exec.Command("git", "clone", arg)
			return
		}
		if len(os.Args) < 3 {
			fmt.Println("command clone requires repository name if --repo is not specified")
			return
		}
		repo := os.Args[2]

		if ok, i := ArgsHaveOption("user", "u"); ok {
			arg := os.Args[i+1]
			repo_string := ""
			if prefs.useHTTPS {
				repo_string = "https://github.com/" + arg + "/" + repo + ".git"
			} else {
				repo_string = "git@github.com:" + arg + "/" + repo + ".git"
			}

			exec.Command("git clone", repo_string)
			return
		}

		repo_string := ""
		if prefs.useHTTPS {
			repo_string = "https://github.com/" + prefs.username + "/" + repo + ".git"
		} else {
			repo_string = "git@github.com:" + prefs.username + "/" + repo + ".git"
		}

		exec.Command("git clone", repo_string)
		return
	}
}

func ArgsHaveOption(long, short string) (hasOption bool, foundAt int) {
	for i, arg := range os.Args {
		if arg == "--"+long || arg == "-"+short {
			return true, i
		}
	}
	return false, 0
}

func writeToFileAsJSON(data Preferances, fileName string) error {
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, b, 0644)
}

func readFromFileAsJSON(fileName string) (Preferances, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Preferances{useHTTPS: true}, err
	}
	var dataRead Preferances
	err = json.Unmarshal(b, &dataRead)
	if err != nil {
		return Preferances{useHTTPS: true}, err
	}
	return dataRead, nil
}
