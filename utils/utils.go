package utils

import (
	"github.com/google/uuid"
	"os"
)

func GetGuid() string {
	return uuid.New().String()
}

func IsTestEnvironment() bool {
	isTestArg := "y"

	if len(os.Args) > 1 {
		isTestArg = os.Args[1]
	}

	switch isTestArg {
	case "n":
		return false
	}

	return true
}
