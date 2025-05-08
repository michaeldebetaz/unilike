package writer

import (
	"fmt"
	"os"
)

func ToFile(filename string, s string) error {
	// Create the file if it doesn't exist, or truncate it if it does
	if _, err := os.Stat(filename); err == nil {
		if err := os.Truncate(filename, 0); err != nil {
			err := fmt.Errorf("Error while truncating file: %v", err)
			return err
		}
	}

	if err := os.WriteFile(filename, []byte(s), 0644); err != nil {
		err := fmt.Errorf("Error while writing to file: %v", err)
		return err
	}

	return nil
}
