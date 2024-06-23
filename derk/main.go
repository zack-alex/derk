package main

import (
	"bufio"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/tekhnus/derk"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/term"
)

type Config struct {
	Salt               string `json:"salt"`
	MasterPasswordHash string `json:"master_password_hash"`
}

func configPath() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "state", "derk", "config.json")
}

func getSalt() (string, error) {
	config, err := readConfig()
	if err != nil {
		return "", err
	}
	return config.Salt, nil
}

func setSalt(salt string) error {
	config, err := readConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if config == nil {
		config = &Config{}
	}
	config.Salt = salt
	return writeConfig(config)
}

func initSalt() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getOrInitSalt() (string, error) {
	salt, err := getSalt()
	if err == nil {
		return salt, nil
	}
	salt = initSalt()
	err = setSalt(salt)
	if err != nil {
		return "", err
	}
	return salt, nil
}

func getMasterPasswordHash() ([]byte, error) {
	config, err := readConfig()
	if err != nil {
		return nil, err
	}
	if config.MasterPasswordHash == "" {
		return nil, nil
	}
	return hex.DecodeString(config.MasterPasswordHash)
}

func setMasterPasswordHash(hash []byte) error {
	config, err := readConfig()
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if config == nil {
		config = &Config{}
	}
	config.MasterPasswordHash = hex.EncodeToString(hash)
	return writeConfig(config)
}

func readConfig() (*Config, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func writeConfig(config *Config) error {
	path := configPath()
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func writeToClipboard(data string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
			cmd = exec.Command("wl-copy")
		} else {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("pbcopy")
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdinPipe for clipboard command:", err)
		os.Exit(1)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting clipboard command:", err)
		os.Exit(1)
	}
	_, err = stdin.Write([]byte(data))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to clipboard command:", err)
		os.Exit(1)
	}
	err = stdin.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error closing StdinPipe for clipboard command:", err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for clipboard command:", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "The password is copied to the clipboard")
}

func passwordHash(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 64)
}

func readPassword() (string, error) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", err
	}
	defer tty.Close()

	pass, err := term.ReadPassword(int(tty.Fd()))
	fmt.Println()
	return string(pass), err
}

func getMasterPassword() (string, error) {
	saltStr, err := getOrInitSalt()
	if err != nil {
		return "", err
	}
	salt := []byte(saltStr)
	fmt.Print("Enter the master passphrase: ")
	masterPassword, err := readPassword()
	if err != nil {
		return "", err
	}
	h, err := getMasterPasswordHash()
	if err != nil {
		return "", err
	}
	if h != nil {
		hashed, err := passwordHash(masterPassword, salt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error hashing master passphrase:", err)
			os.Exit(1)
		}
		for subtle.ConstantTimeCompare(hashed, h) != 1 {
			fmt.Print("Wrong master passphrase, try again: ")
			masterPassword, err = readPassword()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting master passphrase:", err)
				os.Exit(1)
			}
			hashed, err = passwordHash(masterPassword, salt)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error hashing master passphrase:", err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Print("Repeat the master passphrase: ")
		repeatPassword, err := readPassword()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting repeated passphrase:", err)
			os.Exit(1)
		}
		for masterPassword != repeatPassword {
			fmt.Println("Passphrases don't match. Let's do this again.")
			fmt.Print("Enter the master passphrase: ")
			masterPassword, err = readPassword()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting master passphrase:", err)
				os.Exit(1)
			}
			fmt.Print("Repeat the master passphrase: ")
			repeatPassword, err = readPassword()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error getting repeated passphrase:", err)
				os.Exit(1)
			}
		}
		hashed, err := passwordHash(masterPassword, salt)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error hashing master passphrase:", err)
			os.Exit(1)
		}
		setMasterPasswordHash(hashed)
	}

	return string(masterPassword), nil
}

func main() {
	var printFlag bool
	for _, arg := range os.Args[1:] {
		if arg == "--print" {
			printFlag = true
		}
	}

	var specs []map[string]string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &specs)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
			os.Exit(1)
		}
	}

	for _, spec := range specs {
		if err, ok := spec["error"]; ok {
			fmt.Fprintln(os.Stderr, "Error in spec:", err)
			os.Exit(1)
		}
	}

	action := writeToClipboard
	if printFlag {
		action = func(data string) {
			fmt.Print(data)
		}
	}

	masterPassword, err := getMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	for _, spec := range specs {
		password, err := derk.DeriveAndFormat(masterPassword, spec)
		if err != nil {
			log.Fatal(err)
		}
		action(password)
	}
}
