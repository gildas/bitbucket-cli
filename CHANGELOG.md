# Change Log

**2024/02/03** - v0.10.0

- Enhancement #5: the current profile can come from the .git/config file.
- New command: `bb profile which` to display the current profile.
- `bb issue list` now supports multiple `--state` options and defaults to `open` and `new`.

**2024/01/25** - v0.9.0

- Added Pull Request comments.
- All commands follow a `--dry-run` option to test the command before executing it.
- Subcommands like `delete`, `upload`, `download` support a list of arguments to process.
- Fix: `--output format` should supersede profile's option.

**2024/01/05** - v0.8.0

- Initial release.
