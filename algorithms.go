package derk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"golang.org/x/crypto/pbkdf2"
)

func deriveSecretKey(masterPassword string, domain string, username string, counter int) []byte {
	salt := []byte(fmt.Sprintf("%s%s%x", domain, username, counter))
	return pbkdf2.Key([]byte(masterPassword), salt, 100000, 32, sha256.New)
}

func formatBasic(secretKey []byte) string {
	return base58Encode(secretKey[len(secretKey)-12:])
}

func formatHex(secretKey []byte) string {
	return hex.EncodeToString(secretKey)
}

func DeriveAndFormat(masterPassword string, spec map[string]string) (string, error) {
	domain := spec["domain"]
	username := spec["username"]
	method := spec["method"]
	counter := 1
	counterStr, hasCounter := spec["counter"]
	if hasCounter {
		var err error
		counter, err = strconv.Atoi(counterStr)
		if err != nil {
			return "", err
		}
		if counter <= 0 {
			return "", fmt.Errorf("Counter has to be positive")
		}
	}

	// Legacy counter specification.
	switch method {
	case "v1-shorter-count4":
		if hasCounter {
			return "", fmt.Errorf("Counter conflict")
		}
		counter = 4
	case "v1-count3", "v1-shorter-count3":
		if hasCounter {
			return "", fmt.Errorf("Counter conflict")
		}
		counter = 3
	case "v1-count2", "v1-shorter-count2", "v1-with-bang-count2":
		if hasCounter {
			return "", fmt.Errorf("Counter conflict")
		}
		counter = 2
	}

	secretKey := deriveSecretKey(masterPassword, domain, username, counter)

	switch method {
	case "v1", "v1-count2", "v1-count3":
		return formatBasic(secretKey) + "-", nil
	case "v1-wo-tail":
		return formatBasic(secretKey), nil
	case "v1-with-bang", "v1-with-bang-count2":
		return formatBasic(secretKey) + "!", nil
	case "v1-shorter", "v1-shorter-count2", "v1-shorter-count3", "v1-shorter-count4":
		hx := formatBasic(secretKey)
		return hx[:len(hx)-2], nil
	case "v1-shorter-with-dash":
		hx := formatBasic(secretKey)
		return hx[:len(hx)-2] + "-", nil
	case "ethereum":
		return formatHex(secretKey), nil
	case "none":
		return "", nil
	default:
		return "", fmt.Errorf("Unknown method: %s", method)
	}
}
