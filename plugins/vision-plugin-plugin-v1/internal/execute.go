package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	home := os.Getenv("HOME")

	gopath := fmt.Sprintf("%s/go/bin", home)
	fmt.Println(gopath)

	// TODO(luke): currently, this assumes the module us built on github.com
	// Make it easy for developers of plugins to make their own versioning brand of choice available
	downloadUrl := `https://api.github.com/repos/im2nguyen/rover/releases/latest`

	// TODO(genevieve + luke): create test plugin repo with release so we can test rename command below
	// pass headers into curl command so we can access private plugins
	// finish upgrade function

	cmd := fmt.Sprintf(`curl %s | grep browser_download_url | grep darwin_arm64 | cut -d '"' -f 4`, downloadUrl)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("grep output: %v", err)
		return err
	}

	binaryUrl := string(out)
	cmd = fmt.Sprintf(`curl --output-dir ~/ -OL %s`, binaryUrl)
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("downloading latest version: %v", err)
		return err
	}

	// need to trim whitespaces otherwise everything breaks :')
	zipName := filepath.Base(binaryUrl)
	fileName := strings.TrimSuffix(strings.TrimSpace(zipName), ".zip")
	dst := strings.TrimSpace(filepath.Join(home, fileName))
	fmt.Println("test", filepath.Base(zipName))
	src := strings.TrimSpace(filepath.Join(home, zipName))

	_, err = exec.Command("unzip", "-d", dst, src).Output()
	if err != nil {
		fmt.Printf("unzipping latest binary: %v", err)
		return nil
	}

	err = os.Rename(dst+"/rover_v0.3.3", filepath.Join(gopath, "rover_v0.3.3"))
	if err != nil {
		fmt.Printf("moving latest binary to GOBIN: %v", err)
		return err
	}

	return nil
}

// method 1: run curl command for releases different to current version
// method 2: run go get @latest
// method 3: run curl command for latest
// deciding on all three: method 2 first, if that fails, method 1, then method 3
// curl -OL https://github.com/charmbracelet/log/archive/refs/tags/v0.2.4.tar.gz > charmbracelet-v0.2.4.tar.gz

func findArch() (string, error) {
	// switch runtime.GOOS {
	// case "windows":
	// 	dir = Getenv("LocalAppData")
	// 	if dir == "" {
	// 		return "", errors.New("%LocalAppData% is not defined")
	// 	}

	// case "darwin", "ios":
	// 	dir = Getenv("HOME")
	// 	if dir == "" {
	// 		return "", errors.New("$HOME is not defined")
	// 	}
	// 	dir += "/Library/Caches"

	// case "plan9":
	// 	dir = Getenv("home")
	// 	if dir == "" {
	// 		return "", errors.New("$home is not defined")
	// 	}
	// 	dir += "/lib/cache"

	// default: // Unix
	// 	dir = Getenv("XDG_CACHE_HOME")
	// 	if dir == "" {
	// 		dir = Getenv("HOME")
	// 		if dir == "" {
	// 			return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
	// 		}
	// 		dir += "/.cache"
	// 	}
	// }

	// return eg: linux/amd64, darwin/arm64

	return "", nil
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// b, err = exec.Command("curl", downloadUrl, "|", "grep", "browser_download_url", "|", "grep", "darwin_arm64", "|", "cut", "-d", `'"'`, "-f", "4").Output()
