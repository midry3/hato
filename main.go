package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/midry3/hato/internal/data"
	"github.com/midry3/hato/internal/manager"
)

const (
	VERSION string = "v0.1"
)

func printHelp() {
	fmt.Println("[\033[36mhato\033[0m] " + VERSION + "\n" +
		"This is a CLI CheckList tool.\n\n" +

		"\033[35mUsage\033[0m: hato <Action> [Options]\n\n" +

		"[\033[32mActions\033[0m]\n" +
		`  a / add		Add checkList.
  l / list		Show the checklist.
  c / check		Check the list.

` +
		"[\033[32mOptions\033[0m]\n" +
		`  -h  --help	Print help information.
		
If you want more information, plese visit this: https://github.com/midry3/hato`)
	os.Exit(0)
}

func main() {
	m := manager.CreateManager()
	if slices.Contains(os.Args, "-h") || slices.Contains(os.Args, "--help") {
		printHelp()
	}
	if len(os.Args) == 1 {
		if 0 < len(m.GetList()) {
			m.Check()
			return
		} else if !data.IsInilialized {
			printHelp()
		}
		return
	}
	action := os.Args[1]
	switch action {
	case "add", "a":
		content := strings.Join(os.Args[2:], " ")
		if content == "" {
			fmt.Fprintln(os.Stderr, "`\033[33madd\033[0m` needs a content of a checklist.\n")
			printHelp()
		} else {
			m.Add(content)
			fmt.Printf("\033[36mAdded\033[0m: `%s`\n", content)
		}
	case "list", "l":
		for i, c := range m.GetList() {
			fmt.Printf("- [\033[36m%d\033[0m] %s\n", i+1, c)
		}
	case "help", "h":
		printHelp()
	case "check", "c":
		m.Check()
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: `\033[33m%s\033[0m`\n\n", action)
		printHelp()
	}
}
