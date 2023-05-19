package utils

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/scnon/md-doc/internal"
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

func ReaderDoc(repo, file, author, created, updated string, content []byte) (string, error) {
	tmpl, err := template.ParseFiles("./static/doc.html")
	if err != nil {
		return "", err
	}

	var reader bytes.Buffer
	err = tmpl.Execute(&reader, map[string]string{
		"Title":   file,
		"Repo":    repo,
		"Author":  author,
		"Created": created,
		"Updated": updated,
		"Content": string(internal.Render2Html(content)),
	})
	if err != nil {
		return "", err
	}

	return reader.String(), nil
}

func GetFileAuthor(repo, file string) (string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%an", "HEAD", "--", file)
	cmd.Dir = GetDataPath(repo)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	str := string(out)
	strs := strings.Split(str, "\n")

	return strs[0], nil
}

func GetFileCreated(repo, file string) (string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%ad", "--date=format-local:'%Y/%m/%d %H:%M:%S'", "--diff-filter=A", "HEAD", "--", file)
	cmd.Dir = GetDataPath(repo)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	str := string(out)
	strs := strings.Split(str, "\n")

	return strings.ReplaceAll(strs[0], "'", ""), nil
}

func GetFileUpdated(repo, file string) (string, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%ad", "--date=format-local:'%Y/%m/%d %H:%M:%S'", "--diff-filter=M", "HEAD", "--", file)
	cmd.Dir = GetDataPath(repo)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	str := string(out)
	if strings.Trim(str, " ") == "" {
		return GetFileCreated(repo, file)
	}
	strs := strings.Split(str, "\n")
	log.Println(len(strs))

	return strings.ReplaceAll(strs[0], "'", ""), nil
}

func GetFileInfo(repo, file string) (string, string, string) {
	author, err := GetFileAuthor(repo, file)
	if err != nil {
		author = "unknown"
	}
	created, err := GetFileCreated(repo, file)
	if err != nil {
		created = "unknown"
	}
	updated, err := GetFileUpdated(repo, file)
	if err != nil {
		updated = "unknown"
	}

	return author, created, updated
}
