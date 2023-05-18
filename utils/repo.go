package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
