# Changelog

## v0.3.0 (2026-07-13)

### Behavior changes

- `Get` returns the original config pointer instead of the internal retagged copy,
  so it can be type-asserted back to your struct type.
- `wonsz-flag-ignore` is honored only when set exactly to `"true"`; other values
  no longer skip the flag.
- Fields with a custom `mapstructure` tag get their flag name from the tag
  (`mapstructure:"custom_name"` → `--custom-name`); previously the flag was derived
  from the field name and its value never reached the field.
- `net.IP` fields get an IP-typed flag (dotted notation) instead of a hex-bytes flag.
- A multi-character `shortcut` tag returns an error instead of panicking.

### Added

- Add `New` returning an independent `Wonsz` instance; package-level state is gone,
  multiple configs no longer clobber each other. `Get` now returns the original
  config pointer, so it can be type-asserted back.
- Add `WatchConfig` option for config file hot-reload, with an `OnConfigChange`
  callback invoked after each reload with the re-unmarshal result.
- Add `wonsz:"-"` tag to exclude a field from all bindings.
- Support pointer-to-struct config fields.

### Fixed

- `net.IP`/`net.IPNet` values now decode correctly from env variables and config files.
- Non-ASCII field names no longer produce corrupted keys in name conversion.
- `IgnoreViperBindErrors` also covers `viper.BindPFlag` errors.

## v0.2.0 (2026-01-11)

- Support nested structs when binding flags and environment variables.
- Support `time.Time`, `time.Duration`, `net.IP`, `net.IPNet` as config fields and in their bindings.
- Remove logrus dependency; error logs replaced with returned errors.
- Add option to ignore flag binding errors (`IgnoreViperBindErrors`).
- Add option to ignore binding a field to a flag (tag `wonsz-flag-ignore:"true"`).
