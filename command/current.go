package command

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Jabba-Team/jabba/cfg"
)

var lookPath = exec.LookPath

func Current() string {
	javaPath, err := findJavaPath()
	if err == nil {
		prefix := filepath.Join(cfg.Dir(), "jdk") + string(os.PathSeparator)
		if strings.HasPrefix(javaPath, prefix) {
			index := strings.Index(javaPath[len(prefix):], string(os.PathSeparator))
			if index != -1 {
				return javaPath[len(prefix) : len(prefix)+index]
			}
		}
	}
	return ""
}

func findJavaPath() (string, error) {
	javaPath, err := lookPath("java")
	if err == nil {
		// find java
		symLink, isSetSymLink := os.LookupEnv("JABBA_SYMLINK")
		if isSetSymLink {
			if runtime.GOOS == "windows" {
				prefix := symLink + string(os.PathSeparator)
				if strings.HasPrefix(javaPath, prefix) {
					readlink, err := os.Readlink(symLink)
					if err == nil {
						return filepath.Join(readlink, "bin", "java.exe"), nil
					}
					return "", err
				}
			}
		}
	}
	return javaPath, err
}
