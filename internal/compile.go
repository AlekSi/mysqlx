// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s invokes protoc protobuf compiler with right flags.\n", os.Args[0])
	}
	flag.Parse()

	files, err := filepath.Glob("*.proto")
	if err != nil {
		log.Fatal(err)
	}

	mapping := make([]string, len(files))
	commands := make([][]string, len(files))
	for i, f := range files {
		packageName := strings.TrimSuffix(f, filepath.Ext(f))
		if err = os.RemoveAll(packageName); err != nil {
			log.Fatal(err)
		}
		if err = os.MkdirAll(packageName, 0755); err != nil {
			log.Fatal(err)
		}

		mapping[i] = fmt.Sprintf("M%s=github.com/AlekSi/mysqlx/internal/%s", f, packageName)
		commands[i] = []string{"protoc", "--go_out=import_path=" + packageName + ",%s:" + packageName, f}
	}

	m := strings.Join(mapping, ",")
	for _, c := range commands {
		c[1] = fmt.Sprintf(c[1], m)
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Print(strings.Join(cmd.Args, " "))
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
