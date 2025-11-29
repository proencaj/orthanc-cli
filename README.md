# Orthanc CLI

A command-line interface for managing and querying resources on Orthanc DICOM servers.

## Features

- **Configuration Management**: Easy setup and management of Orthanc server credentials
- **Environment Variable Support**: Override configuration with environment variables
- **Single Configuration**: Simple, single-server configuration approach
- **Secure by Default**: Support for HTTPS and basic authentication
- **Cross-Platform**: Builds for Linux and Mac OS

## Installation

Not ready

### Build from Source

#### Using Make

Build and install:

```bash
make install
```

This installs the `orthanc` binary to `/usr/local/bin`.

#### Manual Build

```bash
go build -o orthanc ./cmd/orthanc
```

#### Build for Multiple Platforms

```bash
make build-all
```

This creates binaries in `bin/` for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)

## Quick Start

### 1. Initialize Configuration

```bash
orthanc config init
```

This creates `~/.orthanc-cli.yaml` with default settings.

### 2. Configure Your Server

```bash
orthanc config set orthanc.url http://localhost:8042
orthanc config set orthanc.username orthanc
orthanc config set orthanc.password orthanc
```

### 3. Verify Configuration

```bash
orthanc config list
```

## Usage

### Configuration Commands

```bash
# Initialize config
orthanc config init

# Set configuration values
orthanc config set orthanc.url http://localhost:8042
orthanc config set orthanc.username admin
orthanc config set orthanc.password secret
orthanc config set orthanc.insecure false

# Get configuration values
orthanc config get orthanc.url

# List all configuration
orthanc config list
orthanc config list --show-password  # Show password in plain text
```

### Environment Variables

Override configuration values with environment variables:

```bash
export ORTHANC_URL="http://localhost:8042"
export ORTHANC_USERNAME="admin"
export ORTHANC_PASSWORD="secret"
export ORTHANC_INSECURE="false"
```

Environment variables take precedence over the config file.

## Configuration

The CLI uses a single configuration file located at `~/.orthanc-cli.yaml`:

```yaml
orthanc:
  url: "http://localhost:8042"
  username: "orthanc"
  password: "orthanc"
  insecure: false
```

### Configuration Priority

1. **Environment Variables** (highest priority)
2. **Configuration File** (fallback)

### Custom Config File

Use a custom config file location:

```bash
orthanc --config /path/to/custom-config.yaml config list

## Development

### Project Structure

```

orthanc-cli/
├── cmd/orthanc/ # Main application entry point
├── internal/
│ ├── client/ # Orthanc client wrapper
│ ├── commands/ # CLI commands
│ │ ├── config/ # Config management commands
│ │ └── studies/ # Studies commands (coming soon)
│ └── config/ # Configuration management
├── examples.md # Usage examples
├── Makefile # Build automation
├── go.mod # Go module definition
└── README.md # This file

```

### Adding New Commands

See the [Developer Guide](examples.md#developer-guide) in `examples.md` for information on:
- Using the Orthanc client
- Implementing new commands
- Available client methods

```

## Requirements

- Go 1.25 or later
- Access to an Orthanc DICOM server

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [gorthanc](https://github.com/proencaj/gorthanc) - Orthanc Go client

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See LICENSE file.

## Support

For issues and questions, please open an issue on GitHub.
