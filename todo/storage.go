package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getDataFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".todo.json"), nil
}

func LoadData() (*TodoData, error) {
	path, err := getDataFilePath()
	if err != nil {
		return nil, err
	}

	data := &TodoData{Items: []TodoItem{}}

	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	if len(file) == 0 {
		return data, nil
	}

	err = json.Unmarshal(file, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func SaveData(data *TodoData) error {
	path, err := getDataFilePath()
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, file, 0644)
}
