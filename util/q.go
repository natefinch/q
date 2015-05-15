// package util provides helper functions for the Q application and Q plugins.
package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	// FormatFuncs is map of helper functions for foramtting output from
	// text/template.
	FormatFuncs = map[string]interface{}{
		"local":  func(t time.Time) time.Time { return t.Local() },
		"format": func(t time.Time, format string) string { return t.Format(format) },
	}
)

// AtomicWriteFile atomically writes the given file with the given contents.
func AtomicWriteFile(filename string, r io.Reader) (err error) {
	dir, file := filepath.Split(filename)
	f, err := ioutil.TempFile(dir, file)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %v", err)
	}
	defer func() {
		if err != nil {
			// Don't leave the temp file lying around on error.
			os.Remove(f.Name())
		}
	}()
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("cannot write %q contents: %v", filename, err)
	}
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Chmod(name, info.Mode()); err != nil {
		return err
	}
	if err := ReplaceFile(name, filename); err != nil {
		return fmt.Errorf("cannot replace %q with %q: %v", name, filename, err)
	}
	return nil
}
