package main

import (
	"log/slog"
	"os"

	"yaca/client"
	"yaca/models"
	"yaca/pkg/config"
	"yaca/pkg/logger"
	"yaca/pkg/utils"
)

var (
	utilsLoadEnv                = utils.LoadEnv
	utilsParseArgs              = utils.ParseArgs
	utilsValidateArgs           = utils.ValidateArgs
	utilsHandleError            = utils.HandleError
	clientGetZoneIDByName       = client.GetZoneIDByName
	clientDoesRecordExistOnZone = client.DoesRecordExistOnZone
	clientUpdateRecordOnZone    = client.UpdateRecordOnZone
	clientCreateRecordOnZone    = client.CreateRecordOnZone
	clientDeleteRecordOnZone    = client.DeleteRecordOnZone
)

func run() int {
	// Initialize configuration
	config.Load()
	
	// Initialize logger
	logger.Init()
	logger.Info("Starting Yet Another Cloudflare Action",
		slog.String("environment", config.AppConfig.Environment),
		slog.String("log_level", config.AppConfig.LogLevel))
	
	// Try to load .env file, but don't fail if it doesn't exist
	if err := utilsLoadEnv(); err != nil {
		// Only log as debug since .env is optional in production
		logger.Debug("Could not load .env file",
			slog.String("error", err.Error()))
	}

	args := utilsParseArgs()
	err := utilsValidateArgs(&args)
	utilsHandleError(err, "Failed to validate arguments")

	logger.Debug("Arguments validated",
		slog.String("record_name", args.Record),
		slog.String("zone_name", args.ZoneName),
		slog.Bool("delete", args.Delete),
		slog.String("type", args.Type))

	zoneID, err := clientGetZoneIDByName(args.ZoneName)
	utilsHandleError(err, "Failed to get zone ID",
		slog.String("zone_name", args.ZoneName))

	logger.Info("Zone retrieved",
		slog.String("zone_id", zoneID), // Will be masked automatically
		slog.String("zone_name", args.ZoneName))

	recordID, err := clientDoesRecordExistOnZone(zoneID, args.Record)
	utilsHandleError(err, "Failed to check record existence",
		slog.String("zone_id", zoneID),
		slog.String("record_name", args.Record))

	record := models.Record{
		Record: args.Record,
		Proxy:  args.Proxy,
		Target: args.Target,
		Ttl:    args.Ttl,
		Type:   args.Type,
	}

	if recordID != "" {
		logger.Info("Record exists",
			slog.String("record_name", args.Record),
			slog.String("zone_name", args.ZoneName),
			slog.String("record_id", recordID)) // Will be masked

		if !args.Delete {
			success, err := clientUpdateRecordOnZone(zoneID, recordID, record)
			utilsHandleError(err, "Failed to update record",
				slog.String("zone_id", zoneID),
				slog.String("record_id", recordID))

			if success {
				logger.Info("Record updated successfully",
					slog.String("record_name", args.Record),
					slog.String("zone_name", args.ZoneName),
					slog.String("operation", "update"))
				return 0
			}
		} else {
			success, err := clientDeleteRecordOnZone(zoneID, recordID, record)
			utilsHandleError(err, "Failed to delete record",
				slog.String("zone_id", zoneID),
				slog.String("record_id", recordID))

			if success {
				logger.Info("Record deleted successfully",
					slog.String("record_name", args.Record),
					slog.String("zone_name", args.ZoneName),
					slog.String("operation", "delete"))
				return 0
			}
		}
	} else {
		logger.Info("Record does not exist",
			slog.String("record_name", args.Record),
			slog.String("zone_name", args.ZoneName))

		if args.Delete {
			logger.Warn("Cannot delete non-existent record",
				slog.String("record_name", args.Record),
				slog.String("zone_name", args.ZoneName))
			return 1
		}

		success, err := clientCreateRecordOnZone(zoneID, record)
		utilsHandleError(err, "Failed to create record",
			slog.String("zone_id", zoneID))

		if success {
			logger.Info("Record created successfully",
				slog.String("record_name", args.Record),
				slog.String("zone_name", args.ZoneName),
				slog.String("operation", "create"))
			return 0
		}
	}

	logger.Error("Operation failed unexpectedly")
	return 1
}

func main() {
	os.Exit(run())
}
