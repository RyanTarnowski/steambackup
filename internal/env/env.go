package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Replaced gogotenv with standard library implementation
func LoadEnvironmentFile(envPath string) error {
	//open the .env file
	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	//create a new scanner and loop over each line in the .env file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		//skip over lines that are emply or are commented out
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		//skip over lines that didn't find a "=", these are malformed lines
		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		//trim up any spaces on the key and value
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		//trim away double and single quotes
		value = strings.Trim(value, `"'`)

		//set the environment var
		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("error setting key: %s value: %s", key, value)
		}
	}

	//return any errors encountered by the scanner
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner failed on error: %v", err)
	}
	return nil
}
