package badm

import (
	"os"
	"path/filepath"
)

// Select selects the provided file is it exists.
func Select(configurationPath, path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	path, err = filepath.Abs(fileInfo.Name())
	if err != nil {
		return err
	}

	return updateConfiguration(configurationPath, func(c *Configuration) error {
		c.SelectedPath = path
		return nil
	})
}
