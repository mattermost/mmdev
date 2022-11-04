package node

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattermost/mmdev/model"
	"github.com/mattermost/mmdev/utils"
)

func getNVMCommand(args string) (*exec.Cmd, error) {
	nvmShell := ""
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "NVM_DIR=") {
			nvmShell = fmt.Sprintf("%s/nvm.sh", v[len("NVM_DIR="):])
		}
	}

	if nvmShell == "" {
		return nil, errors.New("NVM_DIR environment variable not found")
	}

	return exec.Command("bash", "-c", fmt.Sprintf("source %s; nvm %s", nvmShell, args)), nil
}

func checkNVMVersion() (string, error) {
	cmd, err := getNVMCommand("--version")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func InstallNVMIfNeeded(version model.VersionConfig) error {
	fmt.Printf("Check & install NVM version %s on Linux\n", version)
	currentVersion, err := checkNVMVersion()
	if err != nil {
		fmt.Printf("Could not find nvm. err = %v\n", err)
	} else {
		fmt.Printf("Found nvm version %s\n", currentVersion)
	}

	if err != nil || !utils.IsCurrentVersionIncluded(currentVersion, version) {
		fmt.Printf("Installing nvm version %s\n", version.MinVersion)
		curlCommand := fmt.Sprintf("curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v%s/install.sh | bash", version)
		installCmd := exec.Command("bash", "-c", curlCommand)
		installCmd.Stdout = os.Stdout
		installCmd.Stdin = os.Stdin

		err = installCmd.Run()
		if err != nil {
			fmt.Printf("Error installing nvm: %v\n", err)
			return err
		}
		newVersion, err := checkNVMVersion()
		if err != nil {
			fmt.Printf("Could not find nvm after install: %v\n", err)
			return err
		} else if !utils.IsCurrentVersionIncluded(newVersion, version) {
			fmt.Printf("The correct version was not installed. Target: %s. Current: %s\n.", version, newVersion)
		} else {
			fmt.Printf("Successfully installed NVM version %s\n", newVersion)
		}
	}
	return nil
}

func checkNodeVersion() (string, error) {
	nodeCommand := exec.Command("node", "-v")
	buf := new(bytes.Buffer)
	nodeCommand.Stdout = buf
	err := nodeCommand.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func InstallNodeIfNeeded(version model.VersionConfig) error {
	fmt.Printf("Check & install NodeJS version %s on Linux\n", version)
	currentVersion, err := checkNodeVersion()
	if err != nil {
		fmt.Println("Could not find node.")
	} else {
		fmt.Printf("Found node version %s\n", currentVersion)
	}

	if err != nil || !utils.IsCurrentVersionIncluded(currentVersion, version) {
		fmt.Printf("Installing NodeJS version %s\n", version.MinVersion)

		nvmCommand, err := getNVMCommand("install " + version.MinVersion)
		if err != nil {
			return err
		}
		nvmCommand.Stdout = os.Stdout
		nvmCommand.Stdin = os.Stdin
		err = nvmCommand.Run()
		if err != nil {
			fmt.Printf("Error installing node: %v\n", err)
			return err
		}

		nvmDefaultCommand, err := getNVMCommand("alias default " + version.MinVersion)
		if err != nil {
			return err
		}
		nvmDefaultCommand.Stdout = os.Stdout
		nvmDefaultCommand.Stdin = os.Stdin
		err = nvmDefaultCommand.Run()
		if err != nil {
			fmt.Printf("Error setting the default node version: %v\n", err)
			return err
		}

		newVersion, err := checkNodeVersion()
		if err != nil {
			fmt.Printf("Could not find node after install: %v\n", err)
			return err
		} else if !utils.IsCurrentVersionIncluded(newVersion, version) {
			fmt.Printf("The correct version was not installed. Target: %v. Current: %s.\n", version, newVersion)
		} else {
			fmt.Printf("Successfully installed NodeJS version %s\n", newVersion)
		}
	}
	return nil
}
