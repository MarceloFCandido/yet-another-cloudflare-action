package utils

import (
	"os"
	"testing"
	"yaca/models"
)

func TestParseArgs(t *testing.T) {
	t.Run("should parse arguments correctly", func(t *testing.T) {
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		// Set command-line arguments
		os.Args = []string{
			"yaca",
			"-r", "test.example.com",
			"-z", "example.com",
			"-t", "192.168.1.1",
			"-y", "A",
			"-p",
			"--ttl", "1800",
		}

		// Parse the arguments
		args := ParseArgs()

		// Check the values
		expected := models.Args{
			Record:   "test.example.com",
			ZoneName: "example.com",
			Target:   "192.168.1.1",
			Type:     "A",
			Proxy:    true,
			Ttl:      1800,
		}

		if args.Record != expected.Record {
			t.Errorf("Record is incorrect, got: %s, want: %s.", args.Record, expected.Record)
		}
		if args.ZoneName != expected.ZoneName {
			t.Errorf("ZoneName is incorrect, got: %s, want: %s.", args.ZoneName, expected.ZoneName)
		}
		if args.Target != expected.Target {
			t.Errorf("Target is incorrect, got: %s, want: %s.", args.Target, expected.Target)
		}
		if args.Type != expected.Type {
			t.Errorf("Type is incorrect, got: %s, want: %s.", args.Type, expected.Type)
		}
		if args.Proxy != expected.Proxy {
			t.Errorf("Proxy is incorrect, got: %t, want: %t.", args.Proxy, expected.Proxy)
		}
		if args.Ttl != expected.Ttl {
			t.Errorf("Ttl is incorrect, got: %f, want: %f.", args.Ttl, expected.Ttl)
		}
	})
}
