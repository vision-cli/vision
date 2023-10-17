package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Executor struct {
	PluginModule string
	// PluginRepo string
}

func (exe *Executor) UpgradeByGo() error {
	// modulePath := buildinfo.BuildInfo.Path
	cmd := exec.Command("go", "get", "github.com/lstratta/vision-plugin-test-v0.0.1"+"@latest")
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error upgrading: %w", err)
	}

	return nil
}

func (exe *Executor) UpgradeByCurl() error {
	// find OS and arch to decide which compiled binary to download
	sysOS, sysArch := runtime.GOOS, runtime.GOARCH
	home := os.Getenv("HOME")
	goBin := fmt.Sprintf("%s/go/bin", home)
	var sleep = func() { time.Sleep(100 * time.Millisecond) }

	// TODO(luke): currently, this assumes the module us built on github.com
	// Make it easy for developers of plugins to make their own versioning brand of choice available
	downloadUrl := exe.PluginModule

	// Currently, the code below only accounts for downloads from github/github apis.
	// This block of code curls the download url (set in config.yml) and finds the browser download url depending on the user's
	// system OS and arch.
	sleep()
	fmt.Println("downloading package...")
	cmd := fmt.Sprintf(`curl %s | grep browser_download_url | grep %s-%s | cut -d '"' -f 4`, downloadUrl, sysOS, sysArch)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("grep output: %v", err)
		return err
	}

	// Downloads the binary obtained from the URL above.
	binUrl := string(out)
	cmd = fmt.Sprintf(`curl --output-dir /tmp -OL %s`, binUrl)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("downloading latest version: %v", err)
		return err
	}

	// Unzips files if files are zipped.
	sleep()
	fmt.Println("unzipping files...")
	binPkg := filepath.Base(binUrl)
	fileName := strings.TrimSuffix(strings.TrimSpace(binPkg), ".zip")
	sp := strings.Split(fileName, "-")
	versionNumber := sp[len(sp)-1]
	dst, isZip, err := unzipBin(binPkg, fileName)
	if err != nil {
		return fmt.Errorf("unzipping: %w", err)
	}

	// Moves files to ~/go/bin.
	sleep()
	fmt.Printf("installing version %s...\n", versionNumber)
	moveFiles(dst, goBin, fileName, isZip)

	// Makes the binary executable.
	sleep()
	fmt.Println("finalising install...")
	_, err = exec.Command("chmod", "+x", filepath.Join(goBin, fileName)).Output()
	if err != nil {
		fmt.Println("changing mode:", err)
		return err
	}

	return nil
}

func unzipBin(binPkg string, fileName string) (string, bool, error) {
	// Need to trim whitespace otherwise it does not work.
	dst := strings.TrimSpace(filepath.Join("/tmp", fileName))
	src := strings.TrimSpace(filepath.Join("/tmp", binPkg))
	isZip := binPkg[len(binPkg)-4:] == ".zip"

	if isZip {
		_, err := exec.Command("unzip", "-d", dst, src).Output()
		if err != nil {
			fmt.Printf("unzipping latest binary: %v", err)
			return "", true, err
		}
	}
	return dst, isZip, nil
}

func moveFiles(oldPath, goBin, fileName string, isZip bool) error {
	var err error
	if isZip {
		err = os.Rename(oldPath+"/"+fileName, filepath.Join(goBin, fileName))
	} else {
		err = os.Rename(oldPath, filepath.Join(goBin, fileName))
	}
	if err != nil {
		fmt.Printf("moving latest binary to GOBIN: %v", err)
		return err
	}
	return nil
}
