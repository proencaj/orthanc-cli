# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.0] - 2025-01-09

### Added

#### DICOMweb Server Management

- List all configured DICOMweb servers (`orthanc servers list`)
- Get detailed information about a DICOMweb server (`orthanc servers get`)
- Create new DICOMweb server configurations (`orthanc servers create`)
- Update existing DICOMweb server configurations (`orthanc servers update`)
- Remove DICOMweb server configurations with confirmation prompt (`orthanc servers remove`)
- Support for authentication, chunked transfers, and WADO-RS universal transfer syntax options

#### DICOMweb Operations

- QIDO-RS: Query for studies, series, and instances with flexible filters
  - Patient name, ID, accession number, study date filters
  - Modality and series number filters
  - Pagination support (limit/offset)
  - Fuzzy matching option
- WADO-RS: Retrieve DICOM objects via RESTful services
  - Download complete studies or series as ZIP archives
  - Extract files to a directory
  - Retrieve study/series/instance metadata as JSON
  - Retrieve rendered instances as JPEG/PNG
  - Retrieve specific frames from multi-frame instances
- WADO-URI: Legacy protocol support for retrieving DICOM instances
  - Content type negotiation (DICOM, JPEG, PNG)
  - Window center/width settings for rendered output
  - Frame selection for multi-frame instances

#### Multi-Context Configuration

- Support for multiple Orthanc server configurations using contexts (similar to kubectl)
- Create, update, and delete named contexts (`orthanc config set-context`, `orthanc config delete-context`)
- Switch between contexts (`orthanc config use-context`)
- List all available contexts (`orthanc config get-contexts`)
- Show current active context (`orthanc config current-context`)
- Rename contexts (`orthanc config rename-context`)
- Automatic migration from old single-server configuration format

#### Configuration Improvements

- Confirmation prompt when overwriting existing configuration file during init
- Check if configuration file already exists before creating

### Changed

- Upgraded gorthanc dependency to v0.4.0
- Updated README with DICOMweb examples and multi-context documentation

### Fixed

- Fixed installation verification command in README (thanks @Pandemonium1986)

## [0.2.0] - 2024-12-03

### Added

- GoReleaser automation for releases
- Multi-platform package support (DEB, RPM, APK)
- Homebrew tap integration
- GitHub Actions release workflow
- Consolidated release documentation (RELEASE_GUIDE.md)

### Changed

- Updated installation instructions in README.md
- Improved Makefile with GoReleaser targets

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

[unreleased]: https://github.com/proencaj/orthanc-cli/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/proencaj/orthanc-cli/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/proencaj/orthanc-cli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/proencaj/orthanc-cli/releases/tag/v0.1.0
