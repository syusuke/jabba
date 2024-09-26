package command

import (
	"github.com/Jabba-Team/jabba/w32"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/Jabba-Team/jabba/cfg"
)

func Use(selector string) ([]string, error) {
	aliasValue := GetAlias(selector)
	if aliasValue != "" {
		selector = aliasValue
	}
	ver, err := LsBestMatch(selector)
	if err != nil {
		return nil, err
	}
	return usePath(filepath.Join(cfg.Dir(), "jdk", ver))
}

func usePath(path string) ([]string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	if runtime.GOOS == "darwin" {
		path = filepath.Join(path, "Contents", "Home")
	}
	systemJavaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		systemJavaHome, _ = os.LookupEnv("JAVA_HOME")
	}
	globalUsePath(path, systemJavaHome)
	return []string{
		"export PATH=\"" + filepath.Join(path, "bin") + string(os.PathListSeparator) + pth + "\"",
		"export JAVA_HOME=\"" + path + "\"",
		"export JAVA_HOME_BEFORE_JABBA=\"" + systemJavaHome + "\"",
	}, nil
}

func globalUsePath(javaHome string, systemJavaHome string) bool {
	symLink, isSetSymLink := os.LookupEnv("JABBA_SYMLINK")
	if isSetSymLink {
		if runtime.GOOS == "windows" {
			sym, _ := os.Lstat(symLink)
			if sym != nil {
				_, err := w32.ElevatedRun("rmdir", filepath.Clean(symLink))
				if err != nil {
					if w32.IsAccessDenied(err) {
						return false
					}
				}
			}

			_, err := w32.ElevatedRun("mklink", "/D", filepath.Clean(symLink), javaHome)
			if err != nil {
				if w32.IsAccessDenied(err) {
					return false
				}
			}
			originHome, _ := os.LookupEnv("JAVA_HOME")
			if originHome != symLink {
				_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME", filepath.Clean(symLink))
			}
			_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME_BEFORE_JABBA", systemJavaHome)
			return true
		} /*else {
			_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME", javaHome)
		}
		_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME_BEFORE_JABBA", systemJavaHome)
		return true*/
	}
	return false
}
