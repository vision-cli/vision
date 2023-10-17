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
	m, err := exec.Command("go", "list", "-m").Output()
	if err != nil {
		return fmt.Errorf("getting module name: %w", err)
	}
	fmt.Println("module name:", string(m))

	// modulePath := buildinfo.BuildInfo.Path
	cmd := exec.Command("go", "get", "github.com/lstrata/vision-plugin-test-v0.0.1"+"@latest")
	_, err = cmd.Output()
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

	// TODO(luke): currently, this assumes the module us built on github.com
	// Make it easy for developers of plugins to make their own versioning brand of choice available
	downloadUrl := exe.PluginModule

	// This block of code curls the download url (set in config.yml) and finds the browser download url depending on the user's
	// system OS and arch.
	// Download the binary
	updateUser("downloading package")
	cmd := fmt.Sprintf(`curl %s | grep browser_download_url | grep %s-%s | cut -d '"' -f 4`, downloadUrl, sysOS, sysArch)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("grep output: %v", err)
		return err
	}
	binUrl := string(out)
	cmd = fmt.Sprintf(`curl --output-dir /tmp -OL %s`, binUrl)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("downloading latest version: %v", err)
		return err
	}

	// Unzip files if files are zipped
	updateUser("unzipping files")
	binPkg := filepath.Base(binUrl)
	fileName := strings.TrimSuffix(strings.TrimSpace(binPkg), ".zip")
	dst, isZip, err := unzipBin(binPkg, fileName)
	if err != nil {
		return fmt.Errorf("unzipping: %w", err)
	}

	// Moves file to ~/go/bin.
	sp := strings.Split(fileName, "-")
	updateUser(fmt.Sprintf("installing version %s", sp[len(sp)-1]))
	moveFiles(dst, goBin, fileName, isZip)

	// Make the binary executable
	updateUser("finalising install")
	err = changeMode(goBin, fileName)
	if err != nil {
		return fmt.Errorf("changing mode: %w", err)
	}
	return nil
}

func updateUser(s string) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("** %s... **\n", s)
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

func changeMode(goBin, fileName string) error {
	_, err := exec.Command("chmod", "+x", filepath.Join(goBin, fileName)).Output()
	if err != nil {
		fmt.Println("changing mode:", err)
		return err
	}
	return nil
}
