package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/midry3/hato/internal/data"
	"github.com/midry3/hato/internal/manager"
)

type ActionType int
type ArgInfo struct {
	Action     ActionType
	Args       []string
	TargetName string
	YamlFile   string
}

const (
	VERSION string     = "v1.2.0"
	ADD     ActionType = iota
	LIST
	CHECK
	CHECKLISTS
	HELP
)

func printHelp() {
	fmt.Println("[\033[36mhato\033[0m] " + VERSION + "\n" +
		"This is a CLI CheckList tool.\n\n" +

		"\033[35mUsage\033[0m: hato <Target> [Options]\n\n" +

		"[\033[32mOptions\033[0m]\n" +
		`  -i / --init		Create hato.yml
  -a / --add		Add checkList.
  -l / --list		Show the checklist.
  -c / --checklists	Show all checklist name.
  -f / --file		Set the path of hato.yml you want to use.
  -g / --global		Use ~/hato.yml

  -h  --help	Print help information.
		
If you want more information, plese visit this: https://github.com/midry3/hato`)
	os.Exit(0)
}

func argCheck() ArgInfo {
	if slices.Contains(os.Args, "-h") || slices.Contains(os.Args, "--help") {
		printHelp()
	} else if slices.Contains(os.Args, "-i") || slices.Contains(os.Args, "--init") {
		data.Inilialize()
		os.Exit(0)
	}
	action_specified := false
	skip := 0
	res := ArgInfo{}
	res.Action = CHECK
	for i, a := range os.Args[1:] {
		if 0 < skip {
			skip--
			continue
		}
		switch a {
		case "-a", "--add":
			if action_specified {
				fmt.Fprintln(os.Stderr, "Multiple actions were specified. Action is only one.")
				os.Exit(1)
			}
			res.Action = ADD
			action_specified = true
		case "-l", "--list":
			if action_specified {
				fmt.Fprintln(os.Stderr, "Multiple actions were specified. Action is only one.")
				os.Exit(1)
			}
			res.Action = LIST
			action_specified = true
		case "-c", "--checklists":
			if action_specified {
				fmt.Fprintln(os.Stderr, "Multiple actions were specified. Action is only one.")
				os.Exit(1)
			}
			res.Action = CHECKLISTS
			action_specified = true
		case "-f", "--file":
			if len(os.Args) <= i+2 {
				fmt.Fprintln(os.Stderr, "You need to set the value of -f, --file option.")
				os.Exit(1)
			}
			f := os.Args[i+2]
			if filepath.Base(f) != "hato.yml" {
				f = filepath.Join(f, "hato.yml")
			}
			res.YamlFile = f
			skip++
		case "-g", "--global":
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			res.YamlFile = filepath.Join(home, "hato.yml")
		default:
			if res.TargetName == "" && !action_specified {
				res.TargetName = a
			} else {
				res.Args = append(res.Args, a)
			}
		}
	}
	if res.TargetName == "" {
		res.TargetName = data.DEFAULT
	}
	if res.YamlFile == "" {
		res.YamlFile = "hato.yml"
	}
	return res
}

func main() {
	if len(os.Args) == 1 {
		m, err := manager.CreateManager("hato.yml", data.DEFAULT)
		if err != nil {
			fmt.Fprintln(os.Stderr, "`\033[33mdefault\033[0m` was not found.")
			os.Exit(1)
		}
		m.Check()
		return
	}
	args := argCheck()
	m, err := manager.CreateManager(args.YamlFile, args.TargetName)
	if err != nil {
		if args.Action == ADD {
			data.NewChecklist(args.TargetName)
			fmt.Printf("\033[36mNew Checklist\033[0m: `%s`\n", args.TargetName)
			m, err = manager.CreateManager(args.YamlFile, args.TargetName)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Checklist `\033[33m%s\033[0m` was not found.\n", args.TargetName)
			os.Exit(1)
		}
	}
	m.Args = args.Args
	switch args.Action {
	case ADD:
		content := strings.Join(args.Args, " ")
		if content == "" {
			fmt.Fprintln(os.Stderr, "`\033[33m--add\033[0m` needs a content of a checklist.\n")
			printHelp()
		} else {
			m.Add(content)
			fmt.Printf("\033[36mAdded\033[0m: `%s` to `\033[33m%s\033[0m`\n", content, args.TargetName)
		}
	case LIST:
		for i, c := range m.GetList() {
			fmt.Printf("- [\033[36m%d\033[0m] %s\n", i+1, c)
		}
	case CHECKLISTS:
		m.ShowAllChecklists()
	case CHECK:
		m.Check()
	}
}
