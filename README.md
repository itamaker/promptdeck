# promptdeck

`promptdeck` is a Go CLI for lightweight and reproducible prompt templating.

It helps teams turn small JSON files into prompt variants and controlled experiment batches without relying on spreadsheets or heavyweight prompt platforms.

## Quickstart

### Install

Install with your preferred method:

```bash
# From the custom tap
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/promptdeck
```

```bash
# Or install from source
go install github.com/itamaker/promptdeck@latest
```

<details>
<summary>You can also download binaries from <a href="https://github.com/itamaker/promptdeck/releases">GitHub Releases</a>.</summary>

Current release archives:

- macOS (Apple Silicon/arm64): `promptdeck_0.1.0_darwin_arm64.tar.gz`
- macOS (Intel/x86_64): `promptdeck_0.1.0_darwin_amd64.tar.gz`
- Linux (arm64): `promptdeck_0.1.0_linux_arm64.tar.gz`
- Linux (x86_64): `promptdeck_0.1.0_linux_amd64.tar.gz`

Each archive contains a single executable: `promptdeck`.

</details>

If the repository is still private, release-based installs require GitHub access to the repository assets.

### First Run

Run:

```bash
promptdeck matrix -template examples/review.tmpl -matrix examples/matrix.json
```

## Requirements

- Go `1.22+`

## Run

Render one prompt:

```bash
go run . render -template examples/review.tmpl -vars examples/vars.json
```

Render a matrix:

```bash
go run . matrix -template examples/review.tmpl -matrix examples/matrix.json
```

## Build From Source

```bash
make build
```

```bash
go build -o dist/promptdeck .
```

## What It Does

1. Loads Go text templates from local files.
2. Renders one prompt from a JSON variable object or many prompts from a JSON array.
3. Expands matrix inputs into Cartesian prompt combinations.
4. Prints output to stdout or writes prompt batches to files.

## Notes

- Use `-out-dir` when you want prompt variants as individual files.
- Maintainer release steps live in `PUBLISHING.md`.
