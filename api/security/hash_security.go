package security

import (
	"belchi/logger"

	"golang.org/x/crypto/bcrypt"
)

type StringHasher struct{}

func NewStringHasher() *StringHasher {
	return &StringHasher{}
}

func (h *StringHasher) MakeHash(plain string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword(
		[]byte(plain),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	logger.Log("DEBUG", "String hashed successfully", []logger.LogDetail{
		{Key: "String", Value: plain},
		{Key: "Hash", Value: string(hashedString)},
	})

	plain = string(hashedString)

	return plain, nil
}

func (h *StringHasher) CompareHash(hashed, plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashed),
		[]byte(plain),
	)

	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			logger.Log("DEBUG", "Hash comparison failed", []logger.LogDetail{
				{Key: "HashedString", Value: hashed},
				{Key: "PlainString", Value: plain},
			})
			return false, nil
		}

		logger.Log("ERROR", "Hash comparison error", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		})
		return false, err
	}

	logger.Log("DEBUG", "Hash comparison successful", []logger.LogDetail{
		{Key: "HashedString", Value: hashed},
		{Key: "PlainString", Value: plain},
	})

	return true, nil
}
