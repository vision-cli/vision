package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Executor struct {
	PluginModule string
	// PluginRepo string
}

func (exe *Executor) UpgradeByGo() error {
	cmd := exec.Command("go", "get", exe.PluginModule+"@latest")
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

	// TODO(luke): currently, this assumes the module us built on github.com
	// Make it easy for developers of plugins to make their own versioning brand of choice available
	downloadUrl := exe.PluginModule

	// TODO(genevieve + luke): create test plugin repo with release so we can test rename command below
	// pass headers into curl command so we can access private plugins
	// finish upgrade function

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

	// need to trim whitespaces otherwise everything breaks :')
	binPkg := filepath.Base(binUrl)
	fileName := strings.TrimSuffix(strings.TrimSpace(binPkg), ".zip")
	dst := strings.TrimSpace(filepath.Join("/tmp", fileName))
	src := strings.TrimSpace(filepath.Join("/tmp", binPkg))
	isZip := binPkg[len(binPkg)-4:] == ".zip"

	fmt.Println("binPkg:", binPkg)
	fmt.Println("filename:", fileName)

	if isZip {
		_, err = exec.Command("unzip", "-d", dst, src).Output()
		if err != nil {
			fmt.Printf("unzipping latest binary: %v", err)
			return nil
		}
	}

	// separate "oldpath" for os.Rename if binPkg is not a zip file
	if isZip {
		err = os.Rename(dst+"/"+fileName, filepath.Join(goBin, fileName))
	} else {
		err = os.Rename(dst, filepath.Join(goBin, fileName))
	}
	if err != nil {
		fmt.Printf("moving latest binary to GOBIN: %v", err)
		return err
	}

	_, err = exec.Command("chmod", "+x", filepath.Join(goBin, fileName)).Output()
	if err != nil {
		fmt.Println("changing mode:", err)
		return err
	}

	return nil
}

// method 1: run curl command for releases different to current version
// method 2: run go get @latest
// method 3: run curl command for latest
// deciding on all three: method 2 first, if that fails, method 1, then method 3
// curl -OL https://github.com/charmbracelet/log/archive/refs/tags/v0.2.4.tar.gz > charmbracelet-v0.2.4.tar.gz
