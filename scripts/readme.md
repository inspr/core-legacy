

# Instructions on how to install the git hooks

Since git doesn't read the contents of the repository that one is cloning, there is the need run the sh script called `install-hooks.sh`.

The process that happens in the script is the installation of the git hooks in the `.git` folder of the project.

The requirements for this script to work is:
- [git](https://git-scm.com/)
- [statickcheck](https://staticcheck.io/docs/install)
- [golint](https://github.com/golang/lint)
- [go vet](https://golang.org/cmd/vet/)
