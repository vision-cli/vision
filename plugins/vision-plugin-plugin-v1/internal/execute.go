package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

type Executor struct {
	// PluginModule string
	// PluginRepo string
}

func (exe *Executor) UpdateByGo(pluginUrl string) error {
	cmd := exec.Command("go", "get", pluginUrl+"@latest")
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error updating: %w", err)
	}

	return nil
}

func (exe *Executor) UpdateByCurl() error {
	// find arch to decide which compiled binary to download
	// arch, err := findArch()
	// if err != nil {
	// 	return fmt.Errorf("finding CPU architecture %v", err)
	// }

	b, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return fmt.Errorf("finding GOPATH %v", err)
	}
	goPath := string(b)
	fmt.Printf("GOPATH: %s", goPath)
	goBin := filepath.Join(goPath, "bin")
	fmt.Printf("GOBIN: %s", goBin)

	// TODO(luke): currently, this assumes the module us built on github.com
	// Make it easy for developers of plugins to make their own versioning brand of choice available
	downloadUrl := `https://api.github.com/repos/im2nguyen/rover/releases/latest`

	cmd := fmt.Sprintf(`curl %s | grep browser_download_url | grep darwin_arm64 | cut -d '"' -f 4`, downloadUrl)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("grep output %v", err)
		return err
	}

	binaryUrl := string(out)
	cmd = fmt.Sprintf(`curl --output-dir /tmp -OL %s`, binaryUrl)
	downloadZip, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("dowloading latest version %v", err)
		return err
	}
	fmt.Println(string(downloadZip))

	_, err = exec.Command("unzip", "/tmp/rover_0.3.3_darwin_arm64.zip", "-d", "/tmp/rover_0.3.3_darwin_arm64").Output()
	if err != nil {
		return fmt.Errorf("unzipping latest binary %v", err)
	}

	_, err = exec.Command("mv", "/tmp/rover_0.3.3_darwin_arm64/rover_v0.3.3", goBin).Output()
	if err != nil {
		return fmt.Errorf("moving latest binary to GOBIN %v", err)
	}

	return nil
}

// method 1: run curl command for releases different to current version
// method 2: run go get @latest
// method 3: run curl command for latest
// deciding on all three: method 2 first, if that fails, method 1, then method 3
// curl -L https://github.com/charmbracelet/log/archive/refs/tags/v0.2.4.tar.gz > charmbracelet-v0.2.4.tar.gz

func findArch() (string, error) {
	cmd := exec.Command("uname", "-am")
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// regex of string to find type and arch
	// return eg: linux/amd64, darwin/arm64

	return string(b), nil
}

// b, err = exec.Command("curl", downloadUrl, "|", "grep", "browser_download_url", "|", "grep", "darwin_arm64", "|", "cut", "-d", `'"'`, "-f", "4").Output()
