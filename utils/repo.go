package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/scnon/md-doc/internal"
)

var (
	DataPath   = "./data/"
	RepoPrefix = "repo"
	GitPrefix  = "git"
)

func GetRepoBase() string {
	return fmt.Sprint(DataPath, RepoPrefix, "/")
}

func GetGitBase() string {
	return fmt.Sprint(DataPath, GitPrefix, "/")
}

func GetGitUrl(name string) string {
	return fmt.Sprint("http://localhost/repo/", name)
}

func GetRepoPath(name string) string {
	return fmt.Sprint(GetRepoBase(), name, "/")
}

func GetGitPath(name string) string {
	return fmt.Sprint(GetGitBase(), name, "/")
}

func CreateRepo(name string) error {
	path := GetRepoPath(name)

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
	path := GetRepoPath(name)
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
	cmd.Dir = GetRepoPath(repo)
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
	cmd.Dir = GetRepoPath(repo)
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
	cmd.Dir = GetRepoPath(repo)
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

func GetFile(repo, file string) ([]byte, error) {
	path := GetGitPath(repo)
	filePath := fmt.Sprint(path, "/", file)

	return ioutil.ReadFile(filePath)
}

func UpdateGit(path string) error {
	_, err := os.Stat(GetGitPath(path))

	if os.IsNotExist(err) {
		_, err := git.PlainClone(GetGitPath(path), false, &git.CloneOptions{
			URL:      GetGitUrl(path),
			Progress: os.Stdout,
		})
		return err
	} else {
		repo, err := git.PlainOpen(GetGitPath(path))
		if err != nil {
			return err
		}

		tree, err := repo.Worktree()
		if err != nil {
			return err
		}

		err = tree.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil {
			return err
		}

		ref, err := repo.Head()
		if err != nil {
			return err
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		file, err := commit.Files()
		if err != nil {
			return err
		}

		file.ForEach(func(f *object.File) error {
			log.Println(f.Name)
			return nil
		})
		return err
	}
}
