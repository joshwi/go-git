package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joshwi/go-git/gitscm"
	"github.com/joshwi/go-utils/utils"
)

var (
	DIRECTORY = utils.Env("DIRECTORY")
	GIT_URL   = utils.Env("GIT_URL")
	GIT_USER  = utils.Env("GIT_USER")
	GIT_TOKEN = utils.Env("GIT_TOKEN")
	GIT_EMAIL = utils.Env("GIT_EMAIL")
)

func main() {

	day := time.Now().Format("2006-01-02")

	name := "nfldb-backup"

	dir := filepath.Join(DIRECTORY + "/" + name)

	project := gitscm.Project{
		Name:      name,
		Directory: dir,
		Url:       GIT_URL,
		Token:     GIT_TOKEN,
		User:      GIT_USER,
		Email:     GIT_EMAIL,
	}

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		project, err = project.Clone(dir)
	} else {
		project, err = project.Open(dir)
	}

	branches, err := project.Branches(name)

	log.Println(branches)

	err = project.Branch(name, day)

	branches, err = project.Branches(name)

	log.Println(branches)

	err = utils.Write(dir, "nfl.txt", "chiefs\nchargers\nraiders\nbroncos\n", 0777)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Add()

	err = project.Commit(fmt.Sprintf("DB backup: %v", time.Now().Format(time.RFC3339)))

	err = project.Push()

}
