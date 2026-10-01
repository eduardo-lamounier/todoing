package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const (
	appFolderName      = "todoing"
	userDataFileName   = "user.json"
	userConfigFileName = "config.json"
)

type UserData struct {
	Tasks []Task `json:"tasks"`
}

type UserConfig struct{}

func getUserDataFilePath() (string, error) {
	appDataFolderPath := filepath.Join(xdg.DataHome, appFolderName)

	if err := os.MkdirAll(appDataFolderPath, 0o755); err != nil {
		return "", err
	}

	userDataFilePath := filepath.Join(appDataFolderPath, userDataFileName)

	return userDataFilePath, nil
}

func getUserConfigFilePath() (string, error) {
	appConfigFolderPath := filepath.Join(xdg.ConfigHome, appFolderName)

	if err := os.MkdirAll(appConfigFolderPath, 0o755); err != nil {
		return "", err
	}

	userConfigFilePath := filepath.Join(appConfigFolderPath, userConfigFileName)
	return userConfigFilePath, nil
}

func LoadUserData() (UserData, error) {
	newError := func(err error) error {
		return fmt.Errorf("failed to load your data: %s", err)
	}

	userFilePath, err := getUserDataFilePath()
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

	userFilePath, err := getUserDataFilePath()
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

func LoadUserConfig() (UserConfig, error) {
	return UserConfig{}, fmt.Errorf("funcionality not implemented")
}

func SaveUserConfig(userConfig UserConfig) error {
	return fmt.Errorf("funcionality not implemented")
}
