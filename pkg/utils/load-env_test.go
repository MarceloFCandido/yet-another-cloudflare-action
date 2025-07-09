package utils

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	t.Run("should load .env file", func(t *testing.T) {
		// Create a dummy .env file
		file, err := os.Create(".env")
		if err != nil {
			t.Fatal("Failed to create .env file:", err)
		}
		defer os.Remove(".env")

		_, err = file.WriteString("TEST_KEY=TEST_VALUE")
		if err != nil {
			t.Fatal("Failed to write to .env file:", err)
		}
		file.Close()

		// Call the function
		err = LoadEnv()
		if err != nil {
			t.Errorf("LoadEnv() returned an error: %v", err)
		}

		// Check if the environment variable is set
		if os.Getenv("TEST_KEY") != "TEST_VALUE" {
			t.Errorf("Environment variable not set correctly")
		}
	})

	t.Run("should return error if .env file does not exist", func(t *testing.T) {
		// Ensure no .env file exists
		os.Remove(".env")

		// Call the function
		err := LoadEnv()
		if err == nil {
			t.Errorf("LoadEnv() should have returned an error, but it didn't")
		}
	})
}
