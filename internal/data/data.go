package data

import (
	"bytes"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var (
	IsInilialized bool = false
)

const (
	TARGETFILE string = "hato.yml"
	HEADER     string = "# https://github.com/midry3/hato"
)

type Config struct {
	CheckList []string
	Actions   []string
}

func (c *Config) Save() {
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

func LoadCheckList() Config {
	f, err := os.ReadFile(TARGETFILE)
	if err != nil {
		return Config{[]string{}, []string{}}
	}
	var res Config
	err = yaml.Unmarshal(f, &res)
	if err != nil {
		fmt.Fprintln(os.Stdout, yaml.FormatError(err, true, true))
		os.Exit(1)
	}
	return res
}

func init() {
	_, err := os.Stat(TARGETFILE)
	if os.IsNotExist(err) {
		os.WriteFile(TARGETFILE, []byte(HEADER+`

checklist: []
actions: []`), 0655)
		fmt.Println("\033[32mInitialized\033[0m: \"hato.yml\"")
		IsInilialized = true
	}
}
