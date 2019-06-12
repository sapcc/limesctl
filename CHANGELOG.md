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
