package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type vendor struct {
	Vendors []struct {
		Path string
		Rev  string
		Hold bool
	}
}

func main() {
	buf, err := ioutil.ReadFile("vendor.yml")
	if err != nil {
		panic(err)
	}

	var v vendor
	err = yaml.Unmarshal(buf, &v)
	if err != nil {
		panic(err)
	}

	gopath := os.Getenv("GOPATH")
	fmt.Println(gopath)

	for _, m := range v.Vendors {
		if strings.Contains(m.Path, "tractrix") {
			fmt.Println("SKIP: ", m.Path)
			continue
		}
		if m.Hold {
			fmt.Println("HOLD: ", m.Path)
			continue
		}
		cmd := exec.Command("git", "reset", "--hard", m.Rev)
		cmd.Dir = gopath + "/src/" + m.Path
		_, err := os.Stat(cmd.Dir)
		if err != nil {
			fmt.Println("NOT FOUND: ", m.Path)
			continue
		}

		out, err := cmd.Output()
		fmt.Println(m.Path, string(out))
	}
}
