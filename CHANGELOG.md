# summon-cerberus changelog

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/) and [Keep a change log](http://keepachangelog.com/).

## [Unreleased]
### Added
- README: example of using the provider directly, without `summon`
- README: added Limitations section

### Modified
- README: fixed provider name in first paragraph of Usage

### Removed
- removed unused Dockerfile

## [0.1.1]
### Added
- allow retrieving all available secrets via `product/environment/` pattern
- added an examples with custom secrets file

### Modified
- moving auth logic into authCerberus function
- build script now generates release archive

## [0.1.0]
- initial release
