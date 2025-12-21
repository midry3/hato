package data

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	TARGETFILE string = "hato.yml"
	HEADER     string = "# https://github.com/midry3/hato"
	DEFAULT    string = "default"
)

type Data struct {
	Aliases   []string
	Inform    []string
	CheckList []string
	Actions   []string
	NArgs     int
}

type Checklists map[string]*Data

func (c *Checklists) Save() {
	b, err := yaml.Marshal(c)
	var buf bytes.Buffer
	if err == nil {
		buf.Write([]byte(HEADER))
		buf.Write([]byte("\n\n"))
		buf.Write(b)
		os.WriteFile(TARGETFILE, buf.Bytes(), 0655)
	} else {
		fmt.Fprintln(os.Stdout, err)
	}
}

func LoadCheckList() Checklists {
	f, err := os.ReadFile(TARGETFILE)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var res Checklists
	err = yaml.Unmarshal(f, &res)
	if err != nil {
		fmt.Fprintln(os.Stdout, yaml.FormatError(err, true, true))
		os.Exit(1)
	}
	return res
}

func NewChecklist(name string) {
	ls := LoadCheckList()
	ls[name] = &Data{}
	ls.Save()
}

func Inilialize() {
	if _, err := os.Stat(TARGETFILE); err == nil {
		fmt.Fprintln(os.Stderr, "`hato.yml` already exists.")
	} else {
		os.WriteFile(TARGETFILE, []byte(fmt.Sprintf(`%s

%s:
  aliases: []
  nargs: 0
  checklist: []
  actions: []`, HEADER, DEFAULT)), 0655)
		fmt.Println("\033[32mInitialized\033[0m: \"hato.yml\"")
	}
}
