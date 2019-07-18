package view

import (
	"fmt"
	"os"
	"strings"
)

func isExistingDir(pth string) bool {
	if fi, err := os.Stat(pth); err == nil {
		return fi.Mode().IsDir()
	}
	return false
}

func getAppRoot() string{
	var AppRoot, _ = os.Getwd()
	if path := os.Getenv("WEB_ROOT"); path != "" {
		AppRoot = path
	}
	return AppRoot
}


func GOPATH() []string {
	paths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(paths) == 0 {
		fmt.Println("GOPATH doesn't exist")
	}
	return paths
}
