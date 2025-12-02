# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-12-02

### Added

#### Patient Management
- List all patients on the Orthanc server
- Get detailed patient information
- Anonymize patient data
- Remove patients with confirmation

#### Study Management
- List all studies on the Orthanc server
- Get detailed study information
- Download study archives as ZIP files
- Anonymize study data
- Remove studies with confirmation
- List all series within a study
- List all instances within a study

#### Series Management
- List all series on the Orthanc server
- Get detailed series information
- Download series archives as ZIP files
- Anonymize series data
- Remove series with confirmation
- List all instances within a series

#### Instance Management
- List all instances on the Orthanc server
- Get detailed instance information
- Upload DICOM files to the server
- Download individual DICOM instances
- Anonymize instance data
- Remove instances with confirmation

#### Modality Management
- List all configured DICOM modalities
- Get modality configuration details
- Create new modality configurations
- Update existing modality configurations
- Remove modality configurations
- Test modality connectivity (C-ECHO)
- Query modalities for studies (C-FIND)
- Move studies to modalities (C-MOVE)
- Retrieve studies from modalities (C-GET)
- Store studies to modalities (C-STORE)

#### System Tools
- Advanced resource search using `/tools/find` endpoint
- Server reset functionality
- Server shutdown commands
- Dynamic log level adjustment

#### Configuration Management
- Initialize configuration file (`~/.orthanc-cli.yaml`)
- Set configuration values interactively
- Get specific configuration values
- List all configuration settings
- Environment variable support for all configuration options
- Custom config file location support (`--config` flag)
- Secure password storage and optional masking

#### General Features
- Cross-platform support (Linux and macOS, amd64 and arm64)
- HTTPS support with optional TLS verification control
- Basic authentication support
- Comprehensive error handling and user feedback
- Consistent command structure across all resource types
- JSON output support for scripting
- Interactive confirmations for destructive operations
- Build automation via Makefile

### Developer Features
- Built on [gorthanc](https://github.com/proencaj/gorthanc) client library
- Modular command structure for easy extension
- Cobra CLI framework integration
- Viper configuration management
- Clean separation of concerns (commands, client, config)

### Documentation
- Comprehensive README with usage examples
- Configuration guide
- Quick start guide
- Development setup instructions
- MIT License

[unreleased]: https://github.com/proencaj/orthanc-cli/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/proencaj/orthanc-cli/releases/tag/v0.1.0
