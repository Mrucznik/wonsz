# Changelog

## Unreleased

- Support nested structs when binding flags and environment variables.
- Support `time.Time`, `time.Duration`, `net.IP`, `net.IPNet` as config fields and in their bindings.
- Remove logrus dependency; error logs replaced with returned errors.
- Add option to ignore flag binding errors (`IgnoreViperBindErrors`).
- Add option to ignore binding a field to a flag (tag `wonsz-flag-ignore:"true"`).
