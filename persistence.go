package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appFolderName    = "todoing"
	userDataFileName = "user.json"
)

type UserData struct {
	Tasks []Task `json:"tasks"`
}

func getUserFilePath() (string, error) {
	configFolderPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appFolderPath := filepath.Join(configFolderPath, appFolderName)
	err = os.MkdirAll(appFolderPath, 0o755)
	if err != nil {
		return "", err
	}

	userFilePath := filepath.Join(appFolderPath, userDataFileName)

	return userFilePath, nil
}

func LoadUserData() (UserData, error) {
	newError := func(err error) error {
		return fmt.Errorf("failed to load your data: %s", err)
	}

	userFilePath, err := getUserFilePath()
	if err != nil {
		return UserData{}, newError(err)
	}

	userFile, err := os.Open(userFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return UserData{}, nil
		}

		return UserData{}, err
	}
	defer userFile.Close()

	decoder := json.NewDecoder(userFile)
	decoder.DisallowUnknownFields()

	userData := UserData{}
	err = decoder.Decode(&userData)
	if err != nil {
		return UserData{}, newError(err)
	}

	return userData, nil
}

func SaveUserData(userData UserData) error {
	newError := func(err error) error {
		return fmt.Errorf("failed to save your data: %s", err)
	}

	userFilePath, err := getUserFilePath()
	if err != nil {
		return newError(err)
	}

	userFile, err := os.Create(userFilePath)
	if err != nil {
		return newError(err)
	}
	defer userFile.Close()

	encoder := json.NewEncoder(userFile)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(userData)
	if err != nil {
		return newError(err)
	}

	return nil
}
