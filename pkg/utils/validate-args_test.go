package utils

import (
	"testing"
	"yaca/models"
)

func TestValidateArgs(t *testing.T) {
	t.Run("should return error when record is empty", func(t *testing.T) {
		args := &models.Args{Record: ""}
		err := ValidateArgs(args)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return error when zone name is empty", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: ""}
		err := ValidateArgs(args)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return error when delete is true and other args are not empty", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: "zone", Delete: true, Target: "target"}
		err := ValidateArgs(args)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return nil when delete is true and other args are empty", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: "zone", Delete: true}
		err := ValidateArgs(args)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})

	t.Run("should return error when target is empty", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: "zone", Target: ""}
		err := ValidateArgs(args)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return error when type is empty", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: "zone", Target: "target", Type: ""}
		err := ValidateArgs(args)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("should return nil when all args are valid", func(t *testing.T) {
		args := &models.Args{Record: "record", ZoneName: "zone", Target: "target", Type: "A"}
		err := ValidateArgs(args)
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}
	})
}
