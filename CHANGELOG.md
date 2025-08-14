# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Modern CLI interface with Cobra framework
- Tailwind CSS integration with automatic processing
- CSS bundling and minification with esbuild
- Professional site scaffolding with opinionated styling
- Development server with configurable host/port
- Cross-platform builds and releases
- Comprehensive installation documentation
- Version management with build-time injection
- GitHub Actions for CI/CD and releases
- Makefile for build automation
- **GitHub Pages integration** with `--github-pages` flag
  - Automatic .gitignore generation
  - GitHub Actions workflow for deployment
  - Repository README template

### Changed

- Refactored monolithic architecture into modular packages
- Enhanced templates with modern Tailwind CSS styling
- Improved error handling throughout application
- Updated README with comprehensive documentation

### Fixed

- Template loading for nested directory structures
- Path handling for cross-platform compatibility
- Asset processing pipeline reliability
- Page output paths when building from within site directory (pages now correctly output to public root, not public/pages/)

## [0.1.0] - 2025-08-13

### Added

- Initial release
- Basic markdown to HTML conversion
- Go template system support
- Simple asset copying
- Directory structure conventions
