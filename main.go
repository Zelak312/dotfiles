package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var qs = &survey.MultiSelect{
	Message: "What do you want to install?",
	Options: []string{"starship config", "bashrc extension", "gitconfig", "motd-diskspace", "alacritty config", "atuin config",
		"vim config", "tmux config"},
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
		case 5:
			installStatus = installAutin(pwd, dirname)
		case 6:
			installStatus = installVimConfig(pwd, dirname)
		case 7:
			installStatus = installTmuxConfig(pwd, dirname)
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

func CheckHandleFileExist(file string, isFolder bool) bool {
	if !isFileExist(file) {
		return true
	}

	delete := false
	prompt := &survey.Confirm{
		Message: "File: " + file + " already exists, do you want to delete it?",
	}

	if isFolder {
		prompt.Message = "Folder: " + file + " is not empty, are you sure you want to delete it?"
	}
	survey.AskOne(prompt, &delete)

	if delete {
		var err error
		if isFolder {
			err = os.RemoveAll(file)
		} else {
			err = os.Remove(file)
		}
		if err != nil {
			if strings.Contains(err.Error(), "directory not empty") {
				CheckHandleFileExist(file, true)
			} else {
				fmt.Println(err)
				return false
			}
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

	continueInstall := CheckHandleFileExist(target, false)
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
	target := path.Join(homedir, ".config", "alacritty")
	file := path.Join(pwd, ".config", "alacritty")
	err := os.MkdirAll(path.Dir(target), 0755)
	if err != nil {
		fmt.Println(err)
		return false
	}

	continueInstall := CheckHandleFileExist(target, false)
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

func installAutin(pwd string, homedir string) bool {
	target := path.Join(homedir, ".config", "atuin")
	file := path.Join(pwd, ".config", "atuin")
	err := os.MkdirAll(path.Dir(target), 0755)
	if err != nil {
		fmt.Println(err)
		return false
	}

	continueInstall := CheckHandleFileExist(target, false)
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

	continueInstall := CheckHandleFileExist(target, false)
	if !continueInstall {
		return false
	}

	err := os.Symlink(file, target)
	if err != nil {
		fmt.Println(err)
		return false
	}

	f, err := os.OpenFile(path.Join(homedir, ".bashrc"), os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()

	// Check if source is already in bashrc
	bashrcLine := "source ~/.bashrc_ext"
	scanner := bufio.NewScanner(f)
	found := false
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == strings.TrimSpace(bashrcLine) {
			found = true
			break // Found it, no need to keep scanning
		}
	}

	// Seek back to the end for appending
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking:", err)
		return false
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

	continueInstall := CheckHandleFileExist(target, false)
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

	continueInstall := CheckHandleFileExist(target, false)
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

func installVimConfig(pwd string, homedir string) bool {
	target := path.Join(homedir, ".vimrc")
	file := path.Join(pwd, ".vimrc")

	continueInstall := CheckHandleFileExist(target, false)
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

func installTmuxConfig(pwd string, homedir string) bool {
	target := path.Join(homedir, ".tmux.conf")
	file := path.Join(pwd, ".tmux.conf")

	continueInstall := CheckHandleFileExist(target, false)
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
