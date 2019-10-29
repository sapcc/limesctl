# v1.6.0 (TBD)

New features:
- The new `--debug` flag logs all HTTP requests and responses, for
  troubleshooting in deployments where mitmproxy is not available for some
  reason.

Bugfixes:
- On Windows, handle UTF-8-encoded environment variables in the same way as
  python-openstackclient. See
  <https://github.com/gophercloud/gophercloud/issues/1572> for details.

# v1.5.3 (2019-08-20)

Bugfixes:
- Do not throw segmentation fault error for invalid service names while setting
  quota(s).

# v1.5.2 (2019-07-17)

Bugfixes:
- A typo that resulted in a previous instance of error not being properly recycled.

# v1.5.1 (2019-06-26)

Changes:
- Report non-existent `physical_usage` data as an empty string in the table and
  csv format.

# v1.5.0 (2019-06-19)

New features:
- Display physical usage information when `--long` output flag is given.
- ID(s) are now optional for `show` and `set` operations. If ID(s) are not
  explicitly given then they are extracted from the current authorization
  token.

Changes:
- Code clean-up.

# v1.4.1 (2019-06-12)

Bugfixes:
- Do not fail project operations when Keystone permissions for domain listing are missing

# v1.4.0 (2019-03-28)

New features:

- Avoid extra requests to Keystone to resolve a domain name into an ID, when
  the token scope already contains the correct domain ID.

# v1.3.0 (2019-01-07)

New features:
- Display quota bursting information when `--long` output flag is given.
- Allow fractional quota values for the `set` subcommand.

Changes:
- Optimize library dependencies. Binary size has been reduced by over 20%.

# v1.2.0 (2018-11-05)

New features:
- Users can manually overwrite the OpenStack environment variables by providing
  them as command-line flags.

Changes:
- For the `--cluster` flag, the domain/project must be identified by ID.
  Specifiying a domain/project name will not work.

Bugfixes:
- `--cluster` flag now works as expected.


# v1.1.0 (2018-10-29)

New features:
- Human friendly values: users can give the `--human-readable` flag to get the
  different quota/usage values in a more human friendly unit. Only valid for
  table/CSV output and can be combined with other output flags: `--names` or
  `--long`.


# v1.0.0 (2018-10-24)

Initial release.
