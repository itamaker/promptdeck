# promptdeck

`promptdeck` keeps prompt experiments lightweight and reproducible.

## Render one prompt

```bash
go run . render -template examples/review.tmpl -vars examples/vars.json
```

## Render a matrix

```bash
go run . matrix -template examples/review.tmpl -matrix examples/matrix.json
```

## Why it is useful

- Turn small JSON files into prompt variants without spreadsheets.
- Generate controlled experiment sets for model or prompt comparison.
- Export prompt batches to files with `-out-dir`.

## Install

From source:

```bash
go install github.com/YOUR_GITHUB_USER/promptdeck@latest
```

From Homebrew after you publish a tap formula:

```bash
brew tap itamaker/tap https://github.com/itamaker/homebrew-tap
brew install itamaker/tap/promptdeck
```

## Repo-Ready Files

- `.github/workflows/ci.yml`
- `.github/workflows/release.yml`
- `.goreleaser.yaml`
- `PUBLISHING.md`
- `scripts/render-homebrew-formula.sh`

## Release

```bash
git tag v0.1.0
git push origin v0.1.0
```

The tagged release workflow publishes multi-platform binaries and `checksums.txt`, which you can feed into the Homebrew formula renderer.
The generated formula should be committed to `https://github.com/itamaker/homebrew-tap`.
