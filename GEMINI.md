# Yet Another Cloudflare Action (YACA)

This project is a Go-based GitHub Action for managing Cloudflare DNS records.

## Overview

The action allows users to create, update, and delete DNS records in Cloudflare zones. It is implemented as a Docker container action.

### Core Functionality

- **Create DNS Records:** If a record does not exist, it will be created.
- **Update DNS Records:** If a record already exists, it will be updated with the new information.
- **Delete DNS Records:** It can explicitly delete a DNS record.

## Project Structure

- `action.yaml`: Defines the GitHub Action interface, including inputs (`record`, `zone-name`, `target`, etc.) and how to run the action.
- `Dockerfile`: Describes the Docker image used to execute the action, which is a minimal Alpine image with the compiled Go binary.
- `entrypoint.sh`: The entry point script for the Docker container (currently empty, but could be used for setup).
- `go.mod`, `go.sum`: Manage the Go project dependencies.

### Go Source Code

- `cmd/yaca/main.go`: The main entry point of the application. It orchestrates the logic by parsing arguments, validating them, and calling the Cloudflare client.
- `client/cloudflare.go`: Contains all the logic for interacting with the Cloudflare API, such as fetching zone IDs and creating/updating/deleting DNS records. It uses a singleton pattern for the Cloudflare API client.
- `models/models.go`: Defines the data structures for command-line arguments (`Args`) and DNS records (`Record`).
- `pkg/utils/`: A package for utility functions:
    - `load-env.go`: Loads environment variables from a `.env` file.
    - `panic-on-error.go`: A simple helper to panic on errors.
    - `parse-args.go`: Parses command-line arguments.
    - `validate-args.go`: Validates the parsed arguments.
