package command

import (
	"github.com/Jabba-Team/jabba/w32"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/Jabba-Team/jabba/cfg"
)

func Deactivate() ([]string, error) {
	pth, _ := os.LookupEnv("PATH")
	rgxp := regexp.MustCompile(regexp.QuoteMeta(filepath.Join(cfg.Dir(), "jdk")) + "[^:]+[:]")
	// strip references to ~/.jabba/jdk/*, otherwise leave unchanged
	pth = rgxp.ReplaceAllString(pth, "")
	javaHome, overrideWasSet := os.LookupEnv("JAVA_HOME_BEFORE_JABBA")
	if !overrideWasSet {
		javaHome, _ = os.LookupEnv("JAVA_HOME")
	}
	globalDeactivate(javaHome)
	return []string{
		"export PATH=\"" + pth + "\"",
		"export JAVA_HOME=\"" + javaHome + "\"",
		"unset JAVA_HOME_BEFORE_JABBA",
	}, nil
}

func globalDeactivate(javaHome string) bool {
	_, isSetSymLink := os.LookupEnv("JABBA_SYMLINK")
	if isSetSymLink {
		if runtime.GOOS == "windows" {
			_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME", javaHome)
			_, _ = w32.ElevatedRun("setx", "/M", "JAVA_HOME_BEFORE_JABBA", "")
			// _, err := w32.ElevatedRun("wmic", "ENVIRONMENT", "where", "name='JAVA_HOME_BEFORE_JABBA'", "delete")
			// if err != nil {
			// 	fmt.Println(err)
			// }
			return true
		}
	}

	return false
}
