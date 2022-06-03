package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type contextKey string

var PathSep = string(os.PathSeparator)

func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if IsError(err) {
		return
	}

	log.Println("File Deleted")
}

func RemoveDir(path string) (err error) {
	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}

func CreateDir(dirName string) (err error) {
	err = os.Mkdir(dirName, 0755)
	if err != nil {
		return err
	}
	return
}

// Copy ..
func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func IsError(err error) bool {
	if err != nil {
		log.Println(err.Error())
	}

	return (err != nil)
}

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Filepathjoin(dirPath string, pathElements ...string) (string, error) {
	p := filepath.Join(append([]string{dirPath}, pathElements...)...)
	p = filepath.FromSlash(p)

	if !strings.HasPrefix(p, dirPath) {
		err := fmt.Errorf("path = %q, should be relative to %q", p, dirPath)
		return p, err
	}
	return p, nil
}

func IsFolderEmpty(dirname string) (bool, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return err == io.EOF, nil
}

func IsFileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func isFolderExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreateBackupFolder(folderName string) (string, error) {
	isExist := isFolderExist(folderName)
	if !isExist {
		err := os.Mkdir(folderName, 0755)
		if err != nil {
			return "", err
		}
	} else {
		parts := strings.Split(folderName, "_")
		if len(parts) == 2 {
			i := 1
			backupFolder := fmt.Sprintf("%s_%d", folderName, i)
			for isFolderExist(backupFolder) {
				i += 1
				backupFolder = fmt.Sprintf("%s_%d", folderName, i)
			}
			err := os.Mkdir(backupFolder, 0755)
			if err != nil {
				return "", err
			}
			return backupFolder, nil
		}
	}
	return folderName, nil
}
