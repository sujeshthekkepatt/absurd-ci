package controller

import (
	"crypto/md5"
	"encoding/hex"
)

func ComputeMd5(input string) (string, error) {

	hash := md5.New()

	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}
