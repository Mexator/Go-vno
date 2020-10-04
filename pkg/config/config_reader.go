package config

import (
	"encoding/json"
	"os"
)

// ReadConfig reads JSON config file to a dest, which should be a pointer to
// structure
func ReadConfig(dest interface{}, filename string) error {
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		err := decoder.Decode(dest)
		return err
	}
	return err
}
