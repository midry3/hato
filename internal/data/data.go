package data

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

var (
	IsInilialized bool = false
)

const (
	TARGETFILE string = "hato.yml"
	HEADER     string = "# https://github.com/midry3/hato"
	DEFAULT    string = "default"
)

type Data struct {
	Aliases   []string
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
		log.Fatal(err)
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

func init() {
	_, err := os.Stat(TARGETFILE)
	if os.IsNotExist(err) {
		os.WriteFile(TARGETFILE, []byte(fmt.Sprintf(`%s

%s:
  default: true
    aliases: []
	nargs: 0
    checklist: []
    actions: []`, HEADER, DEFAULT)), 0655)
		fmt.Println("\033[32mInitialized\033[0m: \"hato.yml\"")
		IsInilialized = true
	}
}
