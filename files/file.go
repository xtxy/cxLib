package files

import (
	"encoding/json"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func ReadConfig(path string, cfgPtr any) error {
	slice, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(slice, cfgPtr)
}
