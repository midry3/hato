package manager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/midry3/hato/internal/data"
	"golang.org/x/term"
)

const (
	DEFAULT string = "_"
)

type Manager struct {
	Config *data.Config
}

func (m *Manager) Check() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	buf := make([]byte, 1)
	for n, c := range m.GetList() {
		fmt.Printf("[\033[36m%d\033[0m]: %s => ?", n+1, c)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			b := buf[0]
			if b == '\r' || b == '\n' {
				fmt.Println("\033[1D✅")
				break
			}
			if b == 0x1b {
				fmt.Println("\033[1D❌")
				return
			}
		}
	}
	term.Restore(int(os.Stdin.Fd()), oldState)
	fmt.Println("All of checklist are ok!")
	n := len(m.Config.Actions)
	if 0 < n {
		for i, c := range m.Config.Actions {
			cmd := strings.Fields(c)
			p := exec.Command(cmd[0], cmd[1:]...)
			p.Stdout = os.Stdout
			p.Stderr = os.Stderr
			fmt.Printf("\n\033[36mRunning \033[32m%d\033[0m/\033[32m%d\033[0m: `%s` ...\n", i+1, n, c)
			if p.Run() != nil {
				fmt.Fprintf(os.Stderr, "\033[31mFaild action\033[0m: `%s`\n", c)
				return
			}
		}
		fmt.Println("\n✅All actions have been \033[33mcompleted\033[0m!")
	}
}

func (m *Manager) Add(content string) {
	m.Config.CheckList = append(m.Config.CheckList, content)
	m.Config.Save()
}

func (m *Manager) Remove(n int) {
	m.Config.CheckList = append(m.Config.CheckList[:n], m.Config.CheckList[n+1:]...)
}

func (m *Manager) GetList() []string {
	return m.Config.CheckList
}

func CreateManager() Manager {
	cfg := data.LoadCheckList()
	return Manager{
		&cfg,
	}
}
