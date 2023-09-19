package project

import (
	"bufio"
	"fmt"
	"os"

	"github.com/vision-cli/vision/plugins/plugin/plugin"

	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/tmpl"
)

func main() {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input += scanner.Text()
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	e := execute.NewOsExecutor()
	t := tmpl.NewOsTmpWriter()
	fmt.Fprint(os.Stdout, plugin.Handle(input, e, t))
}
