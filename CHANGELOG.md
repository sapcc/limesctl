<!--
SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company

SPDX-License-Identifier: Apache-2.0
-->

# Changelog

All notable changes to `limesctl` will be documented in this file.

The sections should follow the order `Added`, `Changed`, `Fixed`, `Removed`, and `Deprecated`.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/sapcc/limesctl/compare/v3.7.0...HEAD)

### Changed

## 3.7.1 - 2025-07-21

### Changed

- Sort subcapacities for capacity reports.
- Updated all dependencies to their latest versions.

## 3.7.0 - 2025-06-16

### Changed

- Fix derivation of service type in `limesctl liquid report-capacity` subcommand.
- Updated all dependencies to their latest versions.

## 3.6.0 - 2025-04-16

### Added

- Added validation of the responses from LIQUID implementations when running a `liquid` subcommand. 

### Changed

- Updated all dependencies to their latest versions.

## 3.5.0 - 2025-02-03

### Added

- Added the `liquid` subcommand family, for testing LIQUID implementations as a cloud operator.

### Changed

- Updated all dependencies to their latest versions.

### Removed

- Removed Muhammad Talal Anwar (@talal) from the list of maintainers. Thanks for all the fish!

## 3.4.0 - 2024-08-29

### Added

- Added `ops validate-quota-overrides` subcommand
  This is to automatically validate a quota-overrides.json file before deploying it into a Limes installation.

### Changed

- Use Golang 1.23 for prebuilt binaries.
- Updated all dependencies to their latest version including Gophercloud to version 2.0

### Removed

- Removed everything related to Bursting
- Removed `domain set` and `project set` subcommands
  Support for writing quotas manually has been removed in Limes.

## 3.3.2 - 2024-01-03

### Changed

- `OS_PW_CMD` is now handled by go-bits which code is mostly adopted from limesctl.
- Updated all dependencies to their latest version.

## 3.3.1 - 2023-10-23

### Fixed

- Completions generation for release archives.

## 3.3.0 - 2023-10-04

### Added

- Support for specifying shell command which will be used for retrieving user password
  using `--os-pw-cmd` flag or `OS_PW_CMD` environment variable.

## 3.2.1 - 2023-09-01

### Changed

- Use Golang 1.21 for prebuilt binaries.
- Updated all dependencies to their latest version.

## 3.2.0 - 2023-03-03

### Added

- Support for showing cluster global rate limits.

  ```
  limesctl cluster show-rates
  ```

## 3.1.3 - 2022-12-15

### Added

- Support for specifying TLS client certificate and key using flags (`--os-cert`/`--os-key`) or
  environment variables (`OS_CERT`/`OS_KEY`).

## 3.1.2 - 2022-11-24

### Fixed

- Examples for `domain set` and `project set` commands.

## 3.1.1 - 2022-11-20

### Fixed

- Removed dead code.

## 3.1.0 - 2022-11-20

### Added

- Add shell completions to Homebrew formula.
- Added detailed examples for `project set` and `domain set` commands.

### Changed

- Updated all dependencies to their latest version.
- Use [Cobra](https://github.com/spf13/cobra) for command-line parsing.
- Improved command descriptions.

## 3.0.3 - 2022-10-25

### Changed

- Updated all dependencies to their latest version.

### Fixed

- Switched to new Limes client for project's rate reports.

## 3.0.2 - 2022-08-30

### Changed

- Use Golang 1.19 in release workflow.

## 3.0.1 - 2022-08-30

### Changed

- Updated all dependencies to their latest version.

## 3.0.0 - 2022-03-03

### Changed

- Updated all dependencies to their latest version.

### Removed

- `cluster list` command; Limes has removed multi-cluster support.
- `--cluster` flag for domain and project subcommands.

## 2.0.1 - 2021-10-06

### Fixed

- Convert given quota value to resource's base unit during relative quota change.

## 2.0.0 - 2021-09-28

### Added

- Support for showing project rate limits.

  ```
  limesctl project show-rates
  limesctl project list-rates
  ```

- Support for specifying multiple area, service, resource values for the respective flags.
- Support for resources without quota.
- Support for relative quota adjustment. The following operators are supported: `+=`,
  `-=`, `*=`, `/=`.

  Example:

  ```
  limesctl project set -q compute/cores+=100
  ```

### Changed

- Switch to [Kong](https://github.com/alecthomas/kong) for command-line
  parsing.
- Domain ID is now optional when listing projects in current token scope.
- `human-readable` flag has been renamed to `humanize`.
- Use flag for quota values instead of positional arguments.

  Old style:

  ```
  limesctl project set compute/cores=250 compute/ram=20GiB
  ```

  New style option 1, use `--quotas` flag with comma-separated values:

  ```
  limesctl project set --quotas=compute/cores=250,compute/ram=20GiB
  ```

  New style option 2, use `-q` shorthand flag for each new quota value:

  ```
  limesctl project set -q compute/cores=250 -q compute/ram=20GiB
  ```

### Removed

- `cluster set` command, Limes' API no longer accepts cluster `PUT` requests.

## 1.6.2 - 2020-11-12

### Fixed

- Show error if domain ID is used for project subcommands.

## 1.6.1 - 2020-07-29

### Changed

- Migrate from gophercloud-limes to gophercloud-sapcc.
- Version flag now prints the Git commit hash and build date.

## 1.6.0 - 2019-11-18

### Added

- The new `--debug` flag logs all HTTP requests and responses, for
  troubleshooting in deployments where mitmproxy is not available for some
  reason.

### Changed

- On Windows, handle UTF-8-encoded environment variables in the same way as
  python-openstackclient. See
  [gophercloud/gophercloud#1572](https://github.com/gophercloud/gophercloud/issues/1572)
  for details.

### Fixed

- Do not crash when unknown service/resource are used with `set` subcommand.

## 1.5.3 - 2019-08-20

### Fixed

- Do not throw segmentation fault error for invalid service names while setting
  quota(s).

## 1.5.2 - 2019-07-17

### Fixed

- A typo that resulted in a previous instance of error not being properly
  recycled.

## 1.5.1 - 2019-06-26

### Changed

- Report non-existent `physical_usage` data as an empty string in the table and
  csv format.

## 1.5.0 - 2019-06-19

### Added

- Display physical usage information when `--long` output flag is given.

### Changed

- ID(s) are now optional for `show` and `set` operations. If ID(s) are not
  explicitly given then they are extracted from the current authorization
  token.

## 1.4.1 - 2019-06-12

### Changed

- Do not fail project operations when Keystone permissions for domain listing
  are missing.

## 1.4.0 - 2019-03-28

### Changed

- Avoid extra requests to Keystone to resolve a domain name into an ID, when
  the token scope already contains the correct domain ID.

## 1.3.0 - 2019-01-07

### Changed

- Display quota bursting information when `--long` output flag is given.
- Allow fractional quota values for the `set` subcommand.
- Optimize library dependencies. Binary size has been reduced by over 20%.

## 1.2.0 - 2018-11-05

### Added

- Users can manually overwrite the OpenStack environment variables by providing
  them as command-line flags.

### Changed

- For the `--cluster` flag, the domain/project must be identified by ID.
  Specifying a domain/project name will not work.

### Fixed

- `--cluster` flag now works as expected.

## 1.1.0 - 2018-10-29

### Added

- Human friendly values: users can give the `--human-readable` flag to get the
  different quota/usage values in a more human friendly unit. Only valid for
  table/CSV output and can be combined with other output flags: `--names` or
  `--long`.

## 1.0.0 - 2018-10-24

### Added

- Initial release.
