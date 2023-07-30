package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
)

var qs = &survey.MultiSelect{
	Message: "What do you want to install?",
	Options: []string{"starship config", "bashrc extension", "gitconfig"},
}

func main() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var configIndexes []int
	survey.AskOne(qs, &configIndexes)

	for _, index := range configIndexes {
		switch index {
		case 0:
			installStarship(pwd, dirname)
		case 1:
			installBashrc(pwd, dirname)
		case 2:
			installGitconfig(pwd, dirname)
		}
	}
}

func installStarship(pwd string, homedir string) {
	target := path.Join(homedir, ".config", "starship.toml")
	file := path.Join(pwd, ".config", "starship.toml")
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
	}
}

func installBashrc(pwd string, homedir string) {
	target := path.Join(homedir, ".bashrc_ext")
	file := path.Join(pwd, ".bashrc_ext")
	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(path.Join(homedir, ".bashrc"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	f.Close()
}

func installGitconfig(pwd string, homedir string) {
	target := path.Join(homedir, ".gitconfig")
	file := path.Join(pwd, ".gitconfig")
	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
	}
}
