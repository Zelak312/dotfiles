package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var qs = &survey.MultiSelect{
	Message: "What do you want to install?",
	Options: []string{"starship config", "bashrc extension", "gitconfig", "motd-diskspace", "alacritty config"},
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
		fmt.Println("Installing: " + qs.Options[index])
		installStatus := false
		switch index {
		case 0:
			installStatus = installStarship(pwd, dirname)
		case 1:
			installStatus = installBashrc(pwd, dirname)
		case 2:
			installStatus = installGitconfig(pwd, dirname)
		case 3:
			installStatus = installMotdDiskspace(pwd, dirname)
		case 4:
			installStatus = installAlacritty(pwd, dirname)
		}

		if installStatus {
			fmt.Println("Success")
		} else {
			fmt.Println("Failed")
		}

		fmt.Println()
	}
}

func isFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	return true
}

func CheckHandleFileExist(file string) bool {
	if !isFileExist(file) {
		return true
	}

	delete := false
	prompt := &survey.Confirm{
		Message: "File: " + file + " already exists, do you want to delete it?",
	}
	survey.AskOne(prompt, &delete)

	if delete {
		err := os.Remove(file)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return delete
}

func installStarship(pwd string, homedir string) bool {
	target := path.Join(homedir, ".config", "starship.toml")
	file := path.Join(pwd, ".config", "starship.toml")
	err := os.MkdirAll(path.Dir(target), 0755)
	if err != nil {
		fmt.Println(err)
		return false
	}

	continueInstall := CheckHandleFileExist(target)
	if !continueInstall {
		return false
	}

	err = os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func installAlacritty(pwd string, homedir string) bool {
	target := path.Join(homedir, ".config", "alacritty", "alacritty.yml")
	file := path.Join(pwd, ".config", "alacritty", "alacritty.yml")
	err := os.MkdirAll(path.Dir(target), 0755)
	if err != nil {
		fmt.Println(err)
		return false
	}

	continueInstall := CheckHandleFileExist(target)
	if !continueInstall {
		return false
	}

	err = os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func installBashrc(pwd string, homedir string) bool {
	target := path.Join(homedir, ".bashrc_ext")
	file := path.Join(pwd, ".bashrc_ext")

	continueInstall := CheckHandleFileExist(target)
	if !continueInstall {
		return false
	}

	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	f, err := os.OpenFile(path.Join(homedir, ".bashrc"), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer f.Close()
	// check if source is already in bashrc
	bashrcLine := "source ~/.bashrc_ext"
	scanner := bufio.NewScanner(f)
	found := false
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(bashrcLine) {
			found = true
		}
	}

	if !found {
		_, err = f.WriteString(fmt.Sprintf("\n%s\n", bashrcLine))
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}

func installGitconfig(pwd string, homedir string) bool {
	target := path.Join(homedir, ".gitconfig")
	file := path.Join(pwd, ".gitconfig")

	continueInstall := CheckHandleFileExist(target)
	if !continueInstall {
		return false
	}

	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func installMotdDiskspace(pwd string, homedir string) bool {
	folder := "/etc/update-motd.d"
	if !isFileExist(folder) {
		fmt.Println("Folder: " + folder + " does not exist\nCan't install motd-diskspace")
		return false
	}

	target := path.Join(folder, "30-diskspace")
	file := path.Join(pwd, "30-diskspace")

	continueInstall := CheckHandleFileExist(target)
	if !continueInstall {
		return false
	}

	err := os.Symlink(file, target)
	if err != nil {
		if !os.IsPermission(err) {
			fmt.Println(err)
			return false
		}

		fmt.Println("Permission error, escalating permissions")
		err := exec.Command("sudo", "ln", "-s", file, target).Run()
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}
