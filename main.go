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

type glide struct {
	Imports []struct {
		Package string `yaml:"package"`
		Version string `yaml:"version"`
	} `yaml:"import"`
}

func resetRepo(gopath string, path string, rev string, hold bool) error {
	mpath := gopath + "/src/" + path
	_, err := os.Stat(mpath)
	if err != nil {
		cmd := exec.Command("go", "get", path)
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("Error go get: ", path)
			return err
		}
	}

	_, err = os.Stat(mpath)
	if err != nil {
		fmt.Println("NOT FOUND: ", path)
		return err
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = mpath
	out, err := cmd.Output()
	fmt.Println(path, string(out))
	if err != nil {
		fmt.Println("Error git pull: ", path)
		return err
	}

	if strings.Contains(path, "tractrix") {
		fmt.Println("SKIP: ", path)
		return nil
	}
	if hold {
		fmt.Println("HOLD: ", path)
		return nil
	}

	cmd = exec.Command("git", "reset", "--hard", rev)
	cmd.Dir = mpath
	out, err = cmd.Output()
	fmt.Println(path, string(out))

	return nil
}

func resetWithVendor(gopath string, buf []byte) {
	var v vendor
	err := yaml.Unmarshal(buf, &v)
	if err != nil {
		panic(err)
	}

	fmt.Println(gopath)

	for _, m := range v.Vendors {
		if err = resetRepo(gopath, m.Path, m.Rev, m.Hold); err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
}

func resetWithGlide(gopath string, buf []byte) {
	var g glide
	err := yaml.Unmarshal(buf, &g)
	if err != nil {
		panic(err)
	}

	fmt.Println(gopath)

	for _, m := range g.Imports {
		if err = resetRepo(gopath, m.Package, m.Version, false); err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
}

func main() {
	gopath := os.Getenv("GOPATH")

	buf, err := ioutil.ReadFile("glide.yaml")
	if err == nil {
		resetWithGlide(gopath, buf)
		return
	}

	buf, err = ioutil.ReadFile("vendor.yml")
	if err == nil {
		resetWithVendor(gopath, buf)
		return
	}

	fmt.Printf("Error: glide.yaml or vendor.yml is required at current working directory.")
}
