# Fast Git Commands
I mean really, git commands take too long to type

## What does this do?
Instead of manually entering multiple commands to perform basic git operations, use fit to reduce the amound of typing you do.

## Help
```yaml
FIT - Fast Git!

Options:
    clone   Clone a Repository with Default Owner of Specified Username
            Requires Additional Argument of Repository Name (Overriden by --repo/-r)

        *Secondary Commands*
        --repo/-r       Specify Repository URL or SSH;
                        Requires Additional Argument of URL or SSH
        --user/-u       Sepcify Owner of Repository
                        Requires Additional Argument of Owner Username

    branch  Switch to Specified Branch (Overriden by any flag)
            Requires Additional Argument of Branch Name (Overriden by --list/-l)

        *Secondary Commands*
        --list/-l       List all Branches
        --del/-d        Delete a Branch
                        Requires Additional Argument of Branch Name
        --new/-n        Create New Branch; Also sets upstream
                        Requires Additional Argument of Branch Name

    push    Adds all changes, Commits With Argument Message, Pushes
            Requires Additional Arugment of Commit Message (Overriden by --none/-n)

        *Secondary Commands*
        --none/-n       Do not Add and Commit, Just Push

    pull    Pull From Upstream

    set     Sets Different Defaults.
            Requires Additional Flags of What to Set

        *Secondary Commands*
        --name/-n       Sets Git Username
                        Rquires Additional Argument of Username
        --email/-e      Sets Git Email
                        Requires Additional Argument of Email.
        --https/-h      Sets FIT to use HTTPS
        --ssh/-s        Sets FIT to use SSH

    raw     Run Argument as Git Command
            Requires Additional Argument of Command

Examples:
    fit clone fast-git-commands              # Clones this repo (My Username is Aadit-Ambadkar)
    fit clone -r git@gith ... nds.git        # Clones this repo with SSH
    fit clone linux --user torvalds          # Clones linux
    fit branch main                          # Switches to the main branch
    fit branch -n dev                        # Switches to a new branch, dev
    fit branch --list                        # Lists all the branches
    fit push "My First Commit"               # Adds, Commits, and Pushes all changes
    fit push --none                          # Pushes all commits
    fit pull                                 # Pulls from Upstream
    fit set --name Aadit-Ambadkar            # Sets name to Aadit Ambadkar
    fit set -s                               # Sets FIT to Use SSH
    fit raw commit -a -m "Github"            # Commits all changes under message "Github"
    fit raw rebase main                      # Rebases with main
```

## Install Directions

**All Methods:**
1. Clone this repository, or download the zip file and unzip it.
2. Ensure that GoLang is installed and available from the command line
3. Ensure that sudo permissions are available.

### Bash
1. Navigate to `/bash-zsh`.
2. Run `bash compile.sh`

### Other Methods
1. Run `go build -o fit` in this directory
2. Add this directory to PATH, or move `fit` to a directory in PATH.