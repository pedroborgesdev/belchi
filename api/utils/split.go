package utils

import (
	"belchi/logger"
	"strings"
)

func SplitPackageName(name string) (string, string) {

	parts := strings.Split(name, "@")

	if len(parts) != 2 {
		logger.Log("ERROR", "invalid package name format to split", []logger.LogDetail{
			{Key: "name", Value: name},
		})
	}

	packAuthor := parts[0]
	packName := parts[1]

	return packAuthor, packName
}
