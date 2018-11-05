# v1.2.0 (2018-11-05)

New features:
- Users can manually overwrite the OpenStack environment variables by providing them as command-line flags.

Changes:
- For the `--cluster` flag, the domain/project must be identified by ID. Specifiying a domain/project name will not work.

Bugfixes:
- `--cluster` flag now works as expected.


# v1.1.0 (2018-10-29)

New features:
- Human friendly values: users can give the `--human-readable` flag to get the different quota/usage values in a more human friendly unit. Only valid for
	table/CSV output and can be combined with other output flags: `--names` or `--long`.


# v1.0.0 (2018-10-24)

Initial release.
