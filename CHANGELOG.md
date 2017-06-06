# summon-cerberus changelog

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/) and [Keep a change log](http://keepachangelog.com/).

## [Unreleased]
### Added

## [0.1.4]
### Added
- set `X-Cerberus-Client` header in request to cerberus

## [0.1.3]
### Added
- allow for the retrival of secrets regardless of Safety Deposit Box

## [0.1.2]
### Modified
- README: fixed provider name in first paragraph of Usage
- README: example of using the provider directly, without `summon`
- README: added Limitations section

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
