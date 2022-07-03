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
	UseHTTPS bool   // If true, uses HTTPS, else uses SSH
	Username string // Users Name

}

var (
	prefs    = Preferances{UseHTTPS: true, Username: ""}
	version  = "dev"
	helpText = "\x1b[32mFIT - Fast Git!\x1b[0m\n\n\x1b[4m\x1b[36mOptions:\x1b[0m\n    clone   Clone a Repository with Default Owner of Specified Username\n            Requires Additional Argument of Repository Name (Overriden by --repo/-r)\n\n        *Secondary Commands*\n        --repo/-r       Specify Repository URL or SSH;\n                        Requires Additional Argument of URL or SSH\n        --user/-u       Sepcify Owner of Repository\n                        Requires Additional Argument of Owner Username\n\n    branch  Switch to Specified Branch (Overriden by any flag)\n            Requires Additional Argument of Branch Name (Overriden by --list/-l)\n\n        *Secondary Commands*\n        --list/-l       List all Branches\n        --del/-d        Delete a Branch\n                        Requires Additional Argument of Branch Name\n        --new/-n        Create New Branch; Also sets upstream\n                        Requires Additional Argument of Branch Name\n\n    push    Adds all changes, Commits With Argument Message, Pushes\n            Requires Additional Arugment of Commit Message (Overriden by --none/-n)\n\n        *Secondary Commands*\n        --none/-n       Do not Add and Commit, Just Push\n\n    pull    Pull From Upstream\n\n    set     Sets Different Defaults.\n            Requires Additional Flags of What to Set\n\n        *Secondary Commands*\n        --name/-n       Sets Git Username\n                        Rquires Additional Argument of Username\n        --email/-e      Sets Git Email\n                        Requires Additional Argument of Email.\n        --https/-h      Sets FIT to use HTTPS\n        --ssh/-s        Sets FIT to use SSH\n\n    raw     Run Argument as Git Command\n            Requires Additional Argument of Command\n\n\x1b[4m\x1b[36mExamples:\x1b[0m\n    fit clone fast-git-commands              # Clones this repo (My Username is Aadit-Ambadkar)\n    fit clone -r git@gith ... nds.git        # Clones this repo with SSH\n    fit clone linux --user torvalds          # Clones linux\n    fit branch main                          # Switches to the main branch\n    fit branch -n dev                        # Switches to a new branch, dev\n    fit branch --list                        # Lists all the branches\n    fit push \"My First Commit\"               # Adds, Commits, and Pushes all changes\n    fit push --none                          # Pushes all commits\n    fit pull                                 # Pulls from Upstream\n    fit set --name Aadit-Ambadkar            # Sets name to Aadit Ambadkar\n    fit set -s                               # Sets FIT to Use SSH\n    fit raw commit -a -m \"Github\"            # Commits all changes under message \"Github\"\n    fit raw rebase main                      # Rebases with main\n"
)

func main() {
	if ok, _ := ArgsHaveOption("help", "notacommand"); ok {
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
			RunCommandInteractive(exec.Command("git", "clone", arg))
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
			if prefs.UseHTTPS {
				repo_string = "https://github.com/" + arg + "/" + repo + ".git"
			} else {
				repo_string = "git@github.com:" + arg + "/" + repo + ".git"
			}

			RunCommandInteractive(exec.Command("git", "clone", repo_string))
			return
		}

		repo_string := ""
		if prefs.UseHTTPS {
			repo_string = "https://github.com/" + prefs.Username + "/" + repo + ".git"
		} else {
			repo_string = "git@github.com:" + prefs.Username + "/" + repo + ".git"
		}

		RunCommandInteractive(exec.Command("git", "clone", repo_string))
		return
	}

	if command == "branch" {
		if ok, _ := ArgsHaveOption("list", "l"); ok {
			RunCommandInteractive(exec.Command("git", "branch -a"))
			return
		}
		if ok, i := ArgsHaveOption("del", "d"); ok {
			if len(os.Args) < i+2 {
				fmt.Println("option --del requires an argument")
				return
			}
			arg := os.Args[i+1]
			RunCommandInteractive(exec.Command("git", "branch -D", arg))
			return
		}
		if ok, i := ArgsHaveOption("new", "n"); ok {
			if len(os.Args) < i+2 {
				fmt.Println("option --new requires an argument")
				return
			}
			arg := os.Args[i+1]
			RunCommandInteractive(exec.Command("git", "fetch"))
			branch_in_upstream := RunCommand(exec.Command("bash", "-c git branch -a | egrep 'remotes/origin/", arg, "'"))
			if branch_in_upstream {
				RunCommandInteractive(exec.Command("git", "checkout --track origin/", arg))
			} else {
				RunCommandInteractive(exec.Command("git", "checkout -b", arg))
				RunCommandInteractive(exec.Command("git", "push -u origin ", arg))
			}
			return
		}
		if len(os.Args) < 3 {
			fmt.Println("command branch requires branch name")
			return
		}
		branch := os.Args[2]
		RunCommandInteractive(exec.Command("git", "checkout", branch))
		return
	}

	if command == "push" {
		if ok, _ := ArgsHaveOption("none", "n"); ok {
			RunCommandInteractive(exec.Command("git", "push"))
			return
		}

		if len(os.Args) < 3 {
			fmt.Println("command push requires commit message")
			return
		}
		msg := os.Args[2]
		RunCommandInteractive(exec.Command("git", "add --all"))
		RunCommandInteractive(exec.Command("git", "commit -a -m \"", msg, "\""))
		RunCommandInteractive(exec.Command("git", "push"))
		return
	}

	if command == "pull" {
		RunCommandInteractive(exec.Command("git", "pull"))
		return
	}

	if command == "set" {
		if ok, i := ArgsHaveOption("name", "n"); ok {
			if len(os.Args) < i+2 {
				fmt.Println("option --name requires an argument")
				return
			}
			arg := os.Args[i+1]
			prefs.Username = arg
			RunCommandInteractive(exec.Command("git", "config --global user.name \"", arg, "\""))
		}

		if ok, i := ArgsHaveOption("email", "e"); ok {
			if len(os.Args) < i+2 {
				fmt.Println("option --email requires an argument")
				return
			}
			arg := os.Args[i+1]
			RunCommandInteractive(exec.Command("git", "config --global user.email \"", arg, "\""))
		}

		if ok, _ := ArgsHaveOption("https", "h"); ok {
			prefs.UseHTTPS = true
		}

		if ok, _ := ArgsHaveOption("ssh", "s"); ok {
			prefs.UseHTTPS = false
		}

		writeToFileAsJSON(prefs, prefsFile)
		return
	}

	if command == "raw" {
		if len(os.Args) < 3 {
			fmt.Println("command raw requires command string (in quotation marks)")
		}
		arg := os.Args[2]
		RunCommandInteractive(exec.Command("git", arg))
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
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, b, 0644)
}

func readFromFileAsJSON(fileName string) (Preferances, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return Preferances{UseHTTPS: true}, err
	}
	var dataRead Preferances
	err = json.Unmarshal(b, &dataRead)
	if err != nil {
		return Preferances{UseHTTPS: true}, err
	}
	return dataRead, nil
}

func RunCommandInteractive(cmd *exec.Cmd) {
	cmd = exec.Command("bash", "-c", cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
}

func RunCommand(cmd *exec.Cmd) bool {
	cmd = exec.Command("bash", "-c", cmd.String())
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}

// test comment
