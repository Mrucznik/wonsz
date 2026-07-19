# Contributing to Wonsz

Thanks for your interest in contributing!

## Development

- Requires Go 1.25+.
- Run tests with the race detector: `go test -race ./...`
- Lint with [golangci-lint](https://golangci-lint.run): `golangci-lint run`
- New behavior needs a test; bug fixes need a regression test.

## Pull requests

1. Fork and create a feature branch.
2. Keep commits focused; follow [Conventional Commits](https://www.conventionalcommits.org)
   (`feat:`, `fix:`, `docs:`, ...).
3. Update `CHANGELOG.md` under the *Unreleased* heading.
4. Make sure CI is green.

## Reporting bugs

Open a GitHub issue with a minimal reproduction (config struct + options +
expected vs. actual behavior).
