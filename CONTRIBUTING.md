# Contributing Guidelines

Thank you for your interest in contributing to our project! We greatly appreciate contributions from the community.

## Issues

If you find a bug, have an idea for a new feature, or want to suggest an improvement, please open an [Issue](https://github.com/Rhyanz46/go-map-validator/issues) in this repository. Be sure to clearly explain the issue and provide any relevant additional information.

## Code Contributions

If you'd like to contribute code, here are the general steps you should follow:

1. Fork this repository.
2. Create a new branch for your feature or fix: `git checkout -b new-feature`.
3. Make the necessary changes.
4. Ensure your code follows the applicable style and standards.
5. Test your changes.
6. Commit your changes: `git commit -m "Add new feature"`.
7. Push your branch to GitHub: `git push origin new-feature`.
8. Create a Pull Request (PR) to the main repository.

We will review your PR as soon as possible. Please be patient, and we'll do our best to respond to every contribution.

## Development & Release Process

This section documents how the library is developed and released so the team can
work consistently.

### Branching & versioning

- `dev` is the working branch; `main` is the release branch.
- Land work on `dev`, then open a PR `dev` → `main`.
- Versioning is incremental patch `v0.0.x` ([Keep a Changelog](https://keepachangelog.com/) format). Confirm the latest existing tag before choosing the next number.
- **CI/CD creates the version tag automatically on merge to `main`. Do not create tags manually.**

### Test-first (TDD)

- Add new behavior and bug fixes **test-first**: write a test that reproduces the
  bug and confirm it **fails** before implementing the fix. This proves the bug is
  real and locks in a regression guard.
- Fix the root cause, minimally — avoid symptom patches and unrelated refactors.

### Verification gate (run before every commit)

All three must be clean:

```bash
make tests          # full suite
go vet ./...        # static checks
go test -race ./... # race detector
```

### Docs to keep in sync

When the public API or behavior changes, update:

- `CHANGELOG.md` — entry for the new version (Added / Fixed / Internal / Migration notes).
- `README.md` — feature docs and Quick Start.
- `AI_GUIDE.md` — guidance for AI agents generating validation code.
- `llms.txt` (pointer-style) and `llms-full.txt` (self-contained) — AI-agent digests; keep the constructor catalog and critical sections current.

For pure internal bug fixes, `CHANGELOG.md` + tests are usually enough.

## License

By contributing to this project, you agree that your contributions will be licensed under the [project's license](https://github.com/Rhyanz46/go-map-validator/blob/main/LICENSE).

## Contact

If you have any questions or need further assistance, feel free to contact us at [rianariansaputra@gmail.com] or [join our telegram](https://t.me/addlist/Wi84VFNkvz85MWFl)

Thank you for your contribution!
