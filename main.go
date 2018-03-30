package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type vendor struct {
	Vendors []struct {
		Path string
		Rev  string
		Hold bool
	}
}

type glide struct {
	Imports []struct {
		Package string `yaml:"package"`
		Version string `yaml:"version"`
	} `yaml:"import"`
}

func vendorReset() {
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

func main() {
	buf, err := ioutil.ReadFile("glide.yaml")
	if err != nil {
		vendorReset()
		return
	}

	var g glide
	err = yaml.Unmarshal(buf, &g)
	if err != nil {
		panic(err)
	}

	gopath := os.Getenv("GOPATH")
	fmt.Println(gopath)

	for _, m := range g.Imports {
		cmd := exec.Command("git", "reset", "--hard", m.Version)
		cmd.Dir = gopath + "/src/" + m.Package
		_, err := os.Stat(cmd.Dir)
		if err != nil {
			fmt.Println("NOT FOUND: ", m.Package)
			continue
		}

		out, err := cmd.Output()
		fmt.Println(m.Package, string(out))
	}
}
