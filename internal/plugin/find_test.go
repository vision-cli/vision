package plugin_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vision-cli/vision/internal/plugin"
)

func TestFind(t *testing.T) {
	tmp1, tmp2 := os.TempDir(), os.TempDir()
	os.Create(filepath.Join(tmp1, "vision-plugin-tmp1-v1"))
	os.Create(filepath.Join(tmp2, "vision-plugin-tmp2-v1"))
	err := os.Setenv("PATH", fmt.Sprintf("%s:%s:%s", tmp1, tmp1, tmp2))
	if err != nil {
		t.Fatal(err)
	}
	plugins, err := plugin.Find()
	assert.Nil(t, err)
	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})
	assert.Equal(t, "tmp1", plugins[0].Name)
	assert.Equal(t, "tmp2", plugins[1].Name)
}
