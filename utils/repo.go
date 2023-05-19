package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/scnon/md-doc/internal"
	"github.com/scnon/md-doc/model"
)

func GetDataPath(name string) string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(path, "/data/git/", name)
}

func CreateRepo(name string) error {
	path := GetDataPath(name)

	if CheckRepoExist(name) {
		log.Println("create failed: repo ", name, " exist")
		return errors.New("repo exist")
	}

	log.Println("begin create repo: ", name)

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	initCmd := exec.Command("git", "init", "--bare")
	initCmd.Dir = path

	if err := initCmd.Run(); err != nil {
		log.Println(err)

		if err := os.RemoveAll(path); err != nil {
			return err
		}
		return err
	}

	return nil
}

func CheckRepoExist(name string) bool {
	path := GetDataPath(name)
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func ReaderDoc(repo string, file string, content []byte) (string, error) {
	tmpl, err := template.New("doc").Parse(model.DocTmpl)
	if err != nil {
		return "", err
	}

	var reader bytes.Buffer
	err = tmpl.Execute(&reader, map[string]string{
		"Title":   "MD-Doc - " + repo + " - " + file,
		"Content": string(internal.Render2Html(content)),
	})
	if err != nil {
		return "", err
	}

	return reader.String(), nil
}
