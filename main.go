package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

var (
	todoDir string
)

const (
	NameFormat = "2006-01-02-Mon.md"
)

func init() {
	flag.StringVar(&todoDir, "todo-directory", "~/Dropbox/Editorial/ToDo", "the directory of todo files")
}

func main() {
	if todoDir == "" {
		log.Fatalf("-todo-directory is a required parameter")
	}

	parts := strings.Split(filepath.ToSlash(todoDir), "/")
	if parts[0] == "~" {
		u, err := user.Current()
		if err != nil {
			log.Fatalf("failed to get current user: %v", err)
		}

		parts[0] = u.HomeDir
	}
	todoDir = filepath.Join(parts...)

	dir, err := ioutil.ReadDir(todoDir)
	if err != nil {
		log.Fatalf("failed to read %s: %v", todoDir, err)
	}

	type fileAndTime struct {
		fi os.FileInfo
		t  time.Time
	}
	var latest *fileAndTime
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}

		t, err := time.Parse(NameFormat, fi.Name())
		if err != nil {
			continue
		}

		if latest == nil || t.After(latest.t) {
			latest = &fileAndTime{fi, t}
		}
	}
	if latest == nil {
		log.Fatalf("failed to find latest file")
	}

	f, err := os.Open(filepath.Join(todoDir, latest.fi.Name()))
	if err != nil {
		log.Fatalf("failed to open %s: %v", latest.fi.Name(), err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read %s: %v", f.Name(), err)
	}

	contents := string(data)

	todo := strings.Index(contents, "## TODO")
	if todo == -1 {
		log.Fatalf("failed to find ## TODO in %s", f.Name())
	}

	done := strings.Index(contents, "## Done")
	if done == -1 {
		log.Fatalf("failed to find ## Done in %s", f.Name())
	}

	nc := contents[todo : done-1]

	nf, err := os.Create(filepath.Join(todoDir, time.Now().Format(NameFormat)))
	if err != nil {
		log.Fatalf("failed to create new file: %v", err)
	}

	_, err = nf.WriteString(
		fmt.Sprintf(
			`# %s

%s
## Done

* nothing
`,
			time.Now().Format("2006-01-02 (Monday)"),
			nc,
		),
	)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}

	log.Printf("created %s", nf.Name())

	cmd := exec.Command("/usr/bin/pbcopy")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("failed to create stdin pipe to pbcopy: %v", err)
	}

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, nf.Name())
		if err != nil {
			log.Printf("failed to write filename to pbcopy's stdin: %v", err)
		}
	}()

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run pbcopy: %v", err)
	}

	log.Println("copied filename to clipboard")
}
