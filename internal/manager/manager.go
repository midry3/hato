package manager

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/midry3/hato/internal/data"
	"golang.org/x/term"
)

type Manager struct {
	Data data.Checklists
	Name string
	Args []string
}

func (m *Manager) applyFormat(s string) string {
	for n, i := range m.Args {
		reg1 := regexp2.MustCompile(fmt.Sprintf(`(?<!%%)%%%d`, n+1), 0)
		reg2 := regexp2.MustCompile(fmt.Sprintf(`(?<!%%)%%\(%d\)`, n+1), 0)
		s, _ = reg1.Replace(s, "\""+i+"\"", 0, -1)
		s, _ = reg2.Replace(s, "\""+i+"\"", 0, -1)
	}
	return s
}

func (m *Manager) Check() {
	if len(m.Args) != m.Data[m.Name].NArgs {
		fmt.Fprintf(os.Stderr, "This checklist needs just \033[33m%d\033[0m arguments.\n", m.Data[m.Name].NArgs)
		os.Exit(1)
	}
	if 0 < len(m.Data[m.Name].Inform) {
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width = 80
		}
		fmt.Printf("%sInformation%s\r\n", strings.Repeat("―", (width-11)/2), strings.Repeat("―", (width-11)/2))
		for _, c := range m.Data[m.Name].Inform {
			RunCmd(m.applyFormat(c))
		}
		fmt.Println(strings.Repeat("―", width))
	}
	list := m.GetList()
	if 0 < len(list) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)
		buf := make([]byte, 1)
		for n, c := range list {
			fmt.Printf("[\033[36m%d\033[0m]: %s => ?", n+1, m.applyFormat(c))
			for {
				_, err := os.Stdin.Read(buf)
				if err != nil {
					log.Fatal(err)
				}
				b := buf[0]
				if b == '\r' || b == '\n' {
					fmt.Print("\033[1D✅\r\n")
					break
				}
				if b == 0x1b {
					fmt.Print("\033[1D❌\r\n")
					return
				}
			}
		}
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Println("All of checklist are ok!")
	}
	n := len(m.Data[m.Name].Actions)
	if 0 < n {
		if 0 < len(list) {
			fmt.Println()
		}
		for i, c := range m.Data[m.Name].Actions {
			fmt.Printf("\033[36mRunning \033[32m%d\033[0m/\033[32m%d\033[0m: `%s` ...\n", i+1, n, m.applyFormat(c))
			if RunCmd(m.applyFormat(c)) != nil {
				fmt.Fprintf(os.Stderr, "\033[31mFaild action\033[0m: `%s`\n", c)
				return
			}
		}
		fmt.Println("\n✅All actions have been \033[33mcompleted\033[0m!")
	}
}

func (m *Manager) Add(content string) {
	m.Data[m.Name].CheckList = append(m.Data[m.Name].CheckList, content)
	m.Data.Save()
}

func (m *Manager) Remove(n int) {
	m.Data[m.Name].CheckList = append(m.Data[m.Name].CheckList[:n], m.Data[m.Name].CheckList[n+1:]...)
}

func (m *Manager) GetList() []string {
	return m.Data[m.Name].CheckList
}

func CreateManager(name string) (Manager, error) {
	ls := data.LoadCheckList()
	_, ok := ls[name]
	if ok {
		return Manager{
			ls,
			name,
			[]string{},
		}, nil
	} else {
		for k, v := range ls {
			if slices.Contains(v.Aliases, name) {
				return Manager{
					ls,
					k,
					[]string{},
				}, nil
			}
		}
		return Manager{}, fmt.Errorf("Not found.\n")
	}
}
