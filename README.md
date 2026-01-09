# Orthanc CLI

A command-line interface for managing Orthanc DICOM servers, designed to streamline daily workflows in medical imaging environments.

## Why Orthanc CLI?

Working with Orthanc DICOM servers through web interfaces or REST APIs directly can be time-consuming and repetitive. Orthanc CLI was created to solve this problem by providing a fast, efficient command-line tool that simplifies common operations and integrates seamlessly into your daily workflow.

### The Problem

Medical imaging professionals, PACS administrators, and developers working with Orthanc often need to:

- Query and manage large volumes of DICOM studies across multiple servers
- Perform bulk operations like anonymization, archiving, or migration
- Automate repetitive tasks in shell scripts or CI/CD pipelines
- Quickly troubleshoot issues without navigating through web UIs
- Integrate Orthanc operations with other command-line tools

Doing these tasks through the web interface is slow, and writing custom scripts against the REST API requires constant API reference lookups and error handling boilerplate.

### The Solution

Orthanc CLI provides a unified, intuitive interface for all common Orthanc operations:

- **Single Command Access**: Manage patients, studies, series, and instances with simple, memorable commands
- **Script-Friendly**: Perfect for automation, batch processing, and integration into existing workflows
- **Modality Management**: Configure and interact with DICOM modalities (C-ECHO, C-FIND, C-MOVE, C-GET, C-STORE)
- **Bulk Operations**: Perform operations across multiple resources efficiently
- **System Administration**: Monitor, configure, and maintain Orthanc servers from the terminal
- **Data Privacy**: Built-in anonymization commands for protecting patient data
- **Fast Iteration**: Quickly test, debug, and explore DICOM data without leaving your terminal

Whether you're a radiologist managing studies, a developer building PACS integrations, or a system administrator maintaining medical imaging infrastructure, Orthanc CLI optimizes your day-to-day interactions with Orthanc servers.

## Features

### Resource Management

- **Patients**: List, query, retrieve, anonymize, and delete patient records
- **Studies**: Full CRUD operations, archiving, and batch processing
- **Series**: Manage series, list instances, download archives
- **Instances**: Upload, download, anonymize, and manage individual DICOM files

### DICOM Networking

- **Modality Configuration**: Create, update, and manage DICOM modalities
- **DICOM Operations**: C-ECHO, C-FIND, C-MOVE, C-GET, and C-STORE support
- **Batch Transfer**: Move or retrieve studies across modalities efficiently

### DICOMweb Integration

- **Server Management**: Configure and manage remote DICOMweb servers
- **WADO-RS/QIDO-RS**: Full DICOMweb protocol support via configured servers (STOW not implemented yet)

### System Operations

- **Server Management**: Monitor system status, adjust log levels, perform maintenance
- **Advanced Search**: Powerful query capabilities using Orthanc's `/tools/find` endpoint
- **Configuration**: Simple, secure credential management with environment variable support

### Developer-Friendly

- **Pipeline Integration**: Exit codes and JSON output for scripting
- **Cross-Platform**: Linux and macOS support (amd64 and arm64)
- **Secure**: HTTPS support, credential encryption, environment variable overrides
- **Extensible**: Built on the [gorthanc](https://github.com/proencaj/gorthanc) library

## Installation

### Homebrew (macOS and Linux)

The easiest way to install Orthanc CLI on macOS or Linux is via Homebrew:

```bash
brew tap proencaj/orthanc-cli
brew install orthanc-cli
```

To upgrade to the latest version:

```bash
brew upgrade orthanc-cli
```

### Debian/Ubuntu (DEB Package)

For Debian-based distributions (Debian, Ubuntu, Linux Mint, etc.):

```bash
# Download the .deb package (replace VERSION with actual version, e.g., 0.1.0)
wget https://github.com/proencaj/orthanc-cli/releases/download/vVERSION/orthanc-cli_VERSION_linux_amd64.deb

# Install the package
sudo dpkg -i orthanc-cli_VERSION_linux_amd64.deb

# Verify installation
orthanc version
```

For ARM64 systems:

```bash
wget https://github.com/proencaj/orthanc-cli/releases/download/vVERSION/orthanc-cli_VERSION_linux_arm64.deb
sudo dpkg -i orthanc-cli_VERSION_linux_arm64.deb
```

### CentOS/RHEL/Fedora (RPM Package)

For RPM-based distributions (CentOS, RHEL, Fedora, Rocky Linux, etc.):

```bash
# Download the .rpm package (replace VERSION with actual version, e.g., 0.1.0)
wget https://github.com/proencaj/orthanc-cli/releases/download/vVERSION/orthanc-cli_VERSION_linux_amd64.rpm

# Install the package
sudo rpm -i orthanc-cli_VERSION_linux_amd64.rpm

# Or use dnf (Fedora/RHEL 8+)
sudo dnf install orthanc-cli_VERSION_linux_amd64.rpm

# Or use yum (CentOS/RHEL 7)
sudo yum install orthanc-cli_VERSION_linux_amd64.rpm

# Verify installation
orthanc version
```

For ARM64 systems:

```bash
wget https://github.com/proencaj/orthanc-cli/releases/download/vVERSION/orthanc-cli_VERSION_linux_arm64.rpm
sudo rpm -i orthanc-cli_VERSION_linux_arm64.rpm
```

### Alpine Linux (APK Package)

For Alpine Linux:

```bash
# Download the .apk package (replace VERSION with actual version)
wget https://github.com/proencaj/orthanc-cli/releases/download/vVERSION/orthanc-cli_VERSION_linux_amd64.apk

# Install the package
sudo apk add --allow-untrusted orthanc-cli_VERSION_linux_amd64.apk

# Verify installation
orthanc version
```

### Pre-built Binaries (Manual Installation)

Download pre-built binaries for your platform from the [GitHub Releases](https://github.com/proencaj/orthanc-cli/releases) page.

#### Linux (amd64)

```bash
curl -LO https://github.com/proencaj/orthanc-cli/releases/latest/download/orthanc-cli-VERSION-linux-amd64.tar.gz
tar -xzf orthanc-cli-VERSION-linux-amd64.tar.gz
sudo mv orthanc /usr/local/bin/
```

#### Linux (arm64)

```bash
curl -LO https://github.com/proencaj/orthanc-cli/releases/latest/download/orthanc-cli-VERSION-linux-arm64.tar.gz
tar -xzf orthanc-cli-VERSION-linux-arm64.tar.gz
sudo mv orthanc /usr/local/bin/
```

#### macOS (Intel)

```bash
curl -LO https://github.com/proencaj/orthanc-cli/releases/latest/download/orthanc-cli-VERSION-darwin-amd64.tar.gz
tar -xzf orthanc-cli-VERSION-darwin-amd64.tar.gz
sudo mv orthanc /usr/local/bin/
```

#### macOS (Apple Silicon)

```bash
curl -LO https://github.com/proencaj/orthanc-cli/releases/latest/download/orthanc-cli-VERSION-darwin-arm64.tar.gz
tar -xzf orthanc-cli-VERSION-darwin-arm64.tar.gz
sudo mv orthanc /usr/local/bin/
```

> **Note**: Replace `VERSION` with the actual version number (e.g., `v0.1.0`), or use `latest` to get the most recent release.

### Building from Source

#### Prerequisites

- Go 1.25 or later
- Make (optional, for convenience)

#### Using Make

```bash
# Clone the repository
git clone https://github.com/proencaj/orthanc-cli.git
cd orthanc-cli

# Build and install to /usr/local/bin
make install
```

#### Manual Build

```bash
# Build for your current platform
go build -o orthanc ./cmd/orthanc

# Move to your PATH
sudo mv orthanc /usr/local/bin/
```

#### Cross-Platform Builds

Build binaries for multiple platforms at once:

```bash
make build-all
```

This creates binaries in the `bin/` directory for:

- Linux (amd64, arm64)
- macOS (amd64, arm64, including Apple Silicon)

## Quick Start

### 1. Initialize Configuration

```bash
orthanc config init
```

This creates `~/.orthanc-cli.yaml` with default settings.

### 2. Configure Your Orthanc Server

You can configure using the simple `set` command (updates the current context):

```bash
orthanc config set orthanc.url http://localhost:8042
orthanc config set orthanc.username orthanc
orthanc config set orthanc.password orthanc
```

Or create named contexts for multiple servers:

```bash
orthanc config set-context local --url http://localhost:8042 --username orthanc --password orthanc
orthanc config set-context production --url https://orthanc.prod.com --username admin --password secret
orthanc config use-context local
```

### 3. Start Using the CLI

```bash
# List all studies
orthanc studies list

# Get details for a specific study
orthanc studies get <study-id>

# Download a study archive
orthanc studies archive <study-id> -o study.zip

# List all configured modalities
orthanc modalities list

# Test connectivity to a modality
orthanc modalities echo <modality-name>
```

## Usage Examples

### Context Management

Manage multiple Orthanc server configurations using contexts (similar to kubectl):

```bash
# List all available contexts
orthanc config get-contexts

# Show current active context
orthanc config current-context

# Create a new context
orthanc config set-context local --url http://localhost:8042 --username orthanc --password orthanc

# Create a production context
orthanc config set-context production --url https://orthanc.prod.com --username admin --password secret

# Create a dev context with TLS verification disabled
orthanc config set-context dev --url http://dev.orthanc.com:8042 --username dev --password dev --insecure

# Switch to a different context
orthanc config use-context production

# Update a context
orthanc config set-context production --username newadmin

# Rename a context
orthanc config rename-context dev staging

# Delete a context (cannot delete current context)
orthanc config delete-context staging

# Get/set config values in the current context
orthanc config get orthanc.url
orthanc config set orthanc.username newuser
orthanc config list
```

### Patient Management

```bash
# List all patients
orthanc patients list

# Get patient details
orthanc patients get <patient-id>

# Anonymize a patient
orthanc patients anonymize <patient-id>

# Remove a patient (with confirmation)
orthanc patients remove <patient-id>
```

### Study Operations

```bash
# List all studies
orthanc studies list

# Download study as ZIP archive
orthanc studies archive <study-id> -o output.zip

# Anonymize a study
orthanc studies anonymize <study-id>

# List all series in a study
orthanc studies list-series <study-id>

# List all instances in a study
orthanc studies list-instances <study-id>
```

### Instance Management

```bash
# Upload a DICOM file
orthanc instances upload /path/to/file.dcm

# Download an instance
orthanc instances download <instance-id> -o output.dcm

# Anonymize an instance
orthanc instances anonymize <instance-id>
```

### Modality Operations

```bash
# Create a new modality
orthanc modalities create REMOTE_PACS --aet REMOTE --host 192.168.1.100 --port 11112

# Test connectivity
orthanc modalities echo REMOTE_PACS

# Find studies on a remote modality
orthanc modalities find REMOTE_PACS --patient-id "12345"

# Move a study to a modality
orthanc modalities move REMOTE_PACS <study-id>

# Retrieve study from modality (C-GET)
orthanc modalities retrieve REMOTE_PACS <study-id>

# Store study to modality (C-STORE)
orthanc modalities store REMOTE_PACS <study-id>
```

### DICOMweb Server Management

```bash
# List all configured DICOMweb servers
orthanc servers list

# List servers with full details
orthanc servers list --expand

# Get details of a specific server
orthanc servers get my-pacs

# Create a new DICOMweb server
orthanc servers create my-pacs --url https://pacs.example.com/dicom-web

# Create with authentication
orthanc servers create my-pacs \
  --url https://pacs.example.com/dicom-web \
  --username admin \
  --password secret

# Create with all options
orthanc servers create my-pacs \
  --url https://pacs.example.com/dicom-web \
  --username admin \
  --password secret \
  --has-delete \
  --chunked-transfers \
  --has-wado-rs-universal-transfer-syntax

# Update an existing server
orthanc servers update my-pacs --url https://new-pacs.example.com/dicom-web

# Enable DELETE support on a server
orthanc servers update my-pacs --has-delete

# Disable chunked transfers (for Orthanc <= 1.5.6)
orthanc servers update my-pacs --chunked-transfers=false

# Remove a server (with confirmation prompt)
orthanc servers remove my-pacs

# Remove without confirmation
orthanc servers remove my-pacs --force
```

### DICOMweb Operations

```bash
# QIDO-RS: Query for studies
orthanc dicomweb qido --level studies

# QIDO-RS: Search studies by patient name (supports wildcards)
orthanc dicomweb qido --level studies --patient-name "Smith*"

# QIDO-RS: Search studies by date range
orthanc dicomweb qido --level studies --study-date 20230101-20231231

# QIDO-RS: Search for all CT series
orthanc dicomweb qido --level series --modality CT

# QIDO-RS: Search series within a specific study
orthanc dicomweb qido --level series --study-uid 1.2.840.113619.2.55.3.123456

# QIDO-RS: Search instances with pagination
orthanc dicomweb qido --level instances --limit 10 --offset 20

# WADO-RS: Retrieve a complete study as a zip archive
orthanc dicomweb wado-rs --study-uid 1.2.840.113619.2.55.3.123456 --output study.zip

# WADO-RS: Retrieve a study and extract files to a directory
orthanc dicomweb wado-rs --study-uid 1.2.840.113619.2.55.3.123456 --output-dir ./study_files/

# WADO-RS: Retrieve a series
orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --output series.zip

# WADO-RS: Retrieve a single instance
orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 --output instance.dcm

# WADO-RS: Retrieve study metadata as JSON
orthanc dicomweb wado-rs --study-uid 1.2.3 --metadata

# WADO-RS: Retrieve rendered instance as JPEG
orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 \
  --rendered --accept image/jpeg --output image.jpg

# WADO-RS: Retrieve specific frames
orthanc dicomweb wado-rs --study-uid 1.2.3 --series-uid 1.2.3.4 --instance-uid 1.2.3.4.5 \
  --frames 1,2,3 --output-dir ./frames/

# WADO-URI: Retrieve a DICOM instance (legacy protocol)
orthanc dicomweb wado --study-uid 1.2.3 --series-uid 1.2.3.4 --object-uid 1.2.3.4.5 --output instance.dcm

# WADO-URI: Retrieve as JPEG with window settings
orthanc dicomweb wado --study-uid 1.2.3 --series-uid 1.2.3.4 --object-uid 1.2.3.4.5 \
  --content-type image/jpeg --window-center 40 --window-width 400 --output image.jpg
```

### System Administration

```bash
# Find resources using advanced queries
orthanc tools find --level Study --query '{"PatientName":"DOE*"}'

# Change log level
orthanc tools log-level default

# Reset the server (careful!)
orthanc tools reset --force

# Shutdown the server
orthanc tools shutdown
```

## Configuration

### Configuration File

The CLI uses a configuration file at `~/.orthanc-cli.yaml` with support for multiple contexts:

```yaml
# Multiple server configurations
contexts:
  local:
    orthanc:
      url: "http://localhost:8042"
      username: "orthanc"
      password: "orthanc"
      insecure: false
  production:
    orthanc:
      url: "https://orthanc.prod.com"
      username: "admin"
      password: "secret"
      insecure: false

# The currently active context
current-context: local

# Global output configuration
output:
  json: false
```

The CLI automatically migrates old single-server configurations to the new multi-context format.

### Environment Variables

Override the current context's configuration with environment variables (useful for CI/CD):

```bash
export ORTHANC_URL="http://localhost:8042"
export ORTHANC_USERNAME="admin"
export ORTHANC_PASSWORD="secret"
export ORTHANC_INSECURE="false"
```

Environment variables take precedence over the current context's values, allowing you to temporarily override settings without modifying the config file.

Example:

```bash
# Use production context but override the URL
orthanc config use-context production
ORTHANC_URL=http://localhost:8042 orthanc studies list
```

### Custom Config File

You can also use a different config file (though contexts are the recommended approach):

```bash
orthanc --config /path/to/alternate-config.yaml studies list
```

**Note**: Using contexts within a single config file is generally more convenient than managing multiple config files.

## Roadmap

The roadmap for future features and improvements is currently being defined. Stay tuned for updates!

If you have feature requests or ideas, please open an issue on GitHub to discuss them.

## Requirements

- Go 1.25 or later (for building from source)
- Access to an Orthanc DICOM server (v1.9.0 or later recommended)

## Dependencies

This project is built on top of excellent open-source libraries:

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management
- [gorthanc](https://github.com/proencaj/gorthanc) - Orthanc Go client library

## Contributing

Contributions are welcome! Whether it's bug reports, feature requests, documentation improvements, or code contributions, all help is appreciated.

Please feel free to:

- Open issues for bugs or feature requests
- Submit pull requests with improvements
- Improve documentation
- Share your use cases and workflows

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/proencaj/orthanc-cli/issues)
- **Documentation**: Check this README and command help (`orthanc --help`, `orthanc <command> --help`)
- **Community**: Discussions and questions are welcome in GitHub Issues

## Acknowledgments

This project builds upon the excellent work of the [Orthanc project](https://www.orthanc-server.com/) by Sébastien Jodogne and contributors. Orthanc is a lightweight DICOM server that has made medical imaging infrastructure more accessible to developers and institutions worldwide.
