# Orthanc CLI {VERSION}

A command-line interface for managing Orthanc DICOM servers.

## What's New

<!-- Summarize the key changes in this release -->

### Features
-

### Improvements
-

### Bug Fixes
-

### Documentation
-

## Installation

### Download Pre-built Binaries

Download the appropriate archive for your platform below, extract it, and add the `orthanc` binary to your PATH.

**Linux:**
- [orthanc-{VERSION}-linux-amd64.tar.gz](./orthanc-{VERSION}-linux-amd64.tar.gz) - Intel/AMD 64-bit
- [orthanc-{VERSION}-linux-arm64.tar.gz](./orthanc-{VERSION}-linux-arm64.tar.gz) - ARM 64-bit

**macOS:**
- [orthanc-{VERSION}-darwin-amd64.tar.gz](./orthanc-{VERSION}-darwin-amd64.tar.gz) - Intel Macs
- [orthanc-{VERSION}-darwin-arm64.tar.gz](./orthanc-{VERSION}-darwin-arm64.tar.gz) - Apple Silicon (M1/M2/M3)

### Quick Install

**Linux/macOS (Intel):**
```bash
curl -L https://github.com/proencaj/orthanc-cli/releases/download/{VERSION}/orthanc-{VERSION}-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64.tar.gz | tar xz
sudo mv orthanc /usr/local/bin/
```

**macOS (Apple Silicon):**
```bash
curl -L https://github.com/proencaj/orthanc-cli/releases/download/{VERSION}/orthanc-{VERSION}-darwin-arm64.tar.gz | tar xz
sudo mv orthanc /usr/local/bin/
```

### Build from Source

Requires Go 1.25 or later:

```bash
go install github.com/proencaj/orthanc-cli/cmd/orthanc@{VERSION}
```

Or clone and build:

```bash
git clone https://github.com/proencaj/orthanc-cli.git
cd orthanc-cli
git checkout {VERSION}
make install
```

## Verification

After installation, verify the version:

```bash
orthanc version
```

## Quick Start

Initialize configuration:

```bash
orthanc config init
orthanc config set orthanc.url http://your-orthanc-server:8042
orthanc config set orthanc.username your-username
orthanc config set orthanc.password your-password
```

Test connectivity:

```bash
orthanc studies list
```

## Documentation

- [README](https://github.com/proencaj/orthanc-cli/blob/{VERSION}/README.md) - Full documentation
- [CHANGELOG](https://github.com/proencaj/orthanc-cli/blob/{VERSION}/CHANGELOG.md) - Complete changelog

## Checksums

SHA256 checksums for release verification:

```
<!-- Paste SHA256SUMS.txt content here -->
```

## Full Changelog

See [CHANGELOG.md](https://github.com/proencaj/orthanc-cli/blob/{VERSION}/CHANGELOG.md) for detailed changes.

**Full Diff**: https://github.com/proencaj/orthanc-cli/compare/{PREVIOUS_VERSION}...{VERSION}

---

## Support

- **Issues**: Report bugs via [GitHub Issues](https://github.com/proencaj/orthanc-cli/issues)
- **Documentation**: Check the [README](https://github.com/proencaj/orthanc-cli#readme)

## Contributing

Contributions are welcome! Please see our [Contributing Guidelines](https://github.com/proencaj/orthanc-cli/blob/main/README.md#contributing).

## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/proencaj/orthanc-cli/blob/{VERSION}/LICENSE) file for details.
