# Changelog

## Unreleased

- Add `New` returning an independent `Wonsz` instance; package-level state is gone,
  multiple configs no longer clobber each other. `Get` now returns the original
  config pointer, so it can be type-asserted back.
- Add `WatchConfig` option for config file hot-reload.
- Add `wonsz:"-"` tag to exclude a field from all bindings.
- Support pointer-to-struct config fields.
- Custom `mapstructure` tags now drive flag names and viper keys, so flag
  overrides work for retagged fields.
- Fix `net.IP`/`net.IPNet` fields: proper flag types and string decoding from
  env and config files.
- Validate the `shortcut` tag instead of panicking on multi-character values.
- `wonsz-flag-ignore` is honored only when set to `"true"`.
- Support non-ASCII field names in name conversion.
- Support nested structs when binding flags and environment variables.
- Support `time.Time`, `time.Duration`, `net.IP`, `net.IPNet` as config fields and in their bindings.
- Remove logrus dependency; error logs replaced with returned errors.
- Add option to ignore flag binding errors (`IgnoreViperBindErrors`).
- Add option to ignore binding a field to a flag (tag `wonsz-flag-ignore:"true"`).
