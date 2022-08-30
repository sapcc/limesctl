# Changelog

All notable changes to `limesctl` will be documented in this file.

The sections should follow the order `Added`, `Changed`, `Fixed`, `Removed`, and `Deprecated`.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/sapcc/limesctl/compare/v3.0.1...HEAD)

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
