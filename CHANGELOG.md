# Changelog

All notable changes to [toon-go](https://github.com/Dicklesworthstone/toon-go) are documented here.

This project has no tagged releases yet. Changes are tracked by commit on the `main` branch.

> **What is toon-go?** Go bindings for [TOON](https://github.com/Dicklesworthstone/toon_rust) (Token-Optimized Object Notation). Wraps the `tru` CLI binary via subprocess to provide `Encode`/`Decode` functions with options, format auto-detection, and structured error handling. Requires Go 1.21+.

---

## Unreleased

All commits below are on `main` with no release tags.

### 2026-03-02 -- Binary discovery: accept `toon` as an alternate binary name

[`d140945`](https://github.com/Dicklesworthstone/toon-go/commit/d14094544b725fc630ff11080b836f464667c2bc)

- **fix:** `findTruBinary` now searches PATH for both `tru` and `toon` (in that order), accepting whichever passes the `isToonRustBinary` verification check.
- **fix:** Removed the hard-ban on binaries named `toon`/`toon.exe` from `isToonRustBinary`. The function now relies solely on `--help` and `--version` output fingerprinting to confirm the binary is toon_rust, regardless of filename.
- Handles distributions that ship the toon_rust binary as `toon` instead of `tru`.

### 2026-02-22 -- License documentation update

[`791a9cf`](https://github.com/Dicklesworthstone/toon-go/commit/791a9cf75f21aeb7007b0a7e6e331eeac1b68210)

- **docs:** Updated README license section from "MIT License" to "MIT License (with OpenAI/Anthropic Rider)" to match the actual LICENSE file.

### 2026-02-21 -- License updated to MIT with OpenAI/Anthropic Rider

[`81d4f28`](https://github.com/Dicklesworthstone/toon-go/commit/81d4f281545a8413f312d5c3390d874c19d4a387)

- **chore:** LICENSE changed from plain MIT to "MIT License (with OpenAI/Anthropic Rider)". The rider excludes OpenAI, Anthropic, and their affiliates from all granted rights unless Jeffrey Emanuel provides express written permission.

### 2026-02-21 -- GitHub social preview image

[`30b6728`](https://github.com/Dicklesworthstone/toon-go/commit/30b672863f5dd8c9813cf7d1373fdc05992bfc56)

- **chore:** Added `gh_og_share_image.png` (1280x640) for GitHub Open Graph social card.

### 2026-01-26 -- Simplify DetectFormat fallback logic

[`3a62b15`](https://github.com/Dicklesworthstone/toon-go/commit/3a62b152d28ef50a917b5f6d3f68b8e60e12ebd1)

- **refactor:** `DetectFormat` no longer attempts TOON-specific heuristics (checking for `": "` or `"]:"`). Any input that is not valid JSON now defaults to `FormatTOON`; invalid TOON will fail at decode time.
- Removes a class of false-negative misdetections where valid TOON lacked the expected key patterns.

### 2026-01-24 -- Fix DetectFormat for JSON scalars

[`250e80b`](https://github.com/Dicklesworthstone/toon-go/commit/250e80b202fb0778e88629b21ebde0207e562379)

- **fix:** `DetectFormat` previously only recognized JSON starting with `{` or `[`. It now uses `json.Unmarshal` as the primary JSON check, correctly detecting scalar JSON values (`"hello"`, `123`, `true`, `null`).
- **test:** Added test cases for JSON scalar detection.

### 2026-01-24 -- Fix tru resolution: accept command names, verify toon_rust identity

[`cec6e36`](https://github.com/Dicklesworthstone/toon-go/commit/cec6e36c2d9794e982825fd0e4ad7ebc7d7a35f6)

- **fix:** `TOON_TRU_BIN` and `TOON_BIN` environment variables now accept both absolute paths and bare command names (resolved via `PATH`).
- **fix:** All binary candidates (env vars, PATH lookup, common paths) are now verified via `isToonRustBinary`, which probes `--help` for "reference implementation in rust" and `--version` for a `tru ` or `toon_rust ` prefix. Prevents silently using an unrelated binary that happens to be named `tru`.
- **fix:** Added `resolveTruCandidate` helper that distinguishes path-like values from command names.
- **docs:** README now includes one-liner install script and corrected `cargo install` command with `--git` and `--tag` flags. Environment variable descriptions clarified to mention "path or command name".

### 2026-01-24 -- Add MIT LICENSE

[`e044b09`](https://github.com/Dicklesworthstone/toon-go/commit/e044b09590e8527bb1992cfbad6600745cf22ce7)

- **chore:** Added initial MIT LICENSE file.

### 2026-01-24 -- Initial release: toon-go library

[`10236ac`](https://github.com/Dicklesworthstone/toon-go/commit/10236ac156a3a5f2d9bb5534dfeceaa8cb077d10)

First commit. Full Go binding library for TOON via the `tru` CLI subprocess.

#### Core API

| Function | Description |
|---|---|
| `Encode(data any) (string, error)` | Convert Go value to TOON (default options) |
| `EncodeWithOptions(data any, opts EncodeOptions) (string, error)` | Convert with custom key-folding, delimiter, indent |
| `Decode(toonStr string, v any) error` | Parse TOON into a Go value (default options) |
| `DecodeWithOptions(toonStr string, opts DecodeOptions, v any) error` | Parse with expand-paths and strict-mode control |
| `DecodeToJSON(toonStr string) (string, error)` | Parse TOON, return raw JSON string |
| `DecodeToJSONWithOptions(toonStr string, opts DecodeOptions) (string, error)` | Same, with options |
| `DecodeToValue(toonStr string) (any, error)` | Parse TOON into `any` |
| `DecodeToValueWithOptions(toonStr string, opts DecodeOptions) (any, error)` | Same, with options |
| `DetectFormat(input string) Format` | Heuristic: `FormatJSON`, `FormatTOON`, or `FormatUnknown` |
| `Convert(input string) (string, Format, error)` | Auto-detect and convert between JSON and TOON |
| `Available() bool` | Check if `tru` binary is reachable |
| `TruPath() (string, error)` | Return resolved path to `tru` |

#### Options structs

- `EncodeOptions` -- `KeyFolding` (off/safe), `FlattenDepth`, `Delimiter` (`,`/`\t`/`|`), `Indent`
- `DecodeOptions` -- `ExpandPaths`, `Strict`

#### Error handling

- All errors wrapped in `*ToonError` with `Code`, `Message`, `Cause` fields.
- Error codes: `ErrCodeEncodeFailed` (10), `ErrCodeDecodeFailed` (11), `ErrCodeTruNotFound` (13).

#### Binary resolution order

1. `TOON_TRU_BIN` env var
2. `TOON_BIN` env var
3. `tru` in `PATH`
4. Common paths: `/usr/local/bin/tru`, `/usr/bin/tru`, `/data/tmp/cargo-target/{release,debug}/tru`

#### Tests

- 20 tests covering encode, decode, format detection, error handling, and option pass-through.

#### Files introduced

- `toon.go` (371 lines) -- library implementation
- `toon_test.go` (502 lines) -- test suite
- `go.mod` -- module `github.com/Dicklesworthstone/toon-go`, requires Go 1.21
- `README.md` -- full documentation with API reference, quick start, integration patterns
