# Release Template

## Version

- Version:
- Commit:
- Tree SHA:

## Required Evidence

- [ ] `GOWORK=off make release-check`
- [ ] `go run ./cmd/xlibgate score --min 9.8`
- [ ] `GOWORK=off make release-final-check`
- [ ] `GOWORK=off make release-preflight VERSION=<version>`
- [ ] `release/manifest/latest.json` generated and validated
- [ ] `release/manifest/latest.json.sha256` generated and validated

## Score

- Score JSON:
- Threshold:
- Status:

## Notes

- Template behavior impact:
- Generated library impact:
- Known gaps:

## Artifact

- CI artifact:
- Workflow run ID:
- Artifact URL:
- Local manifest:
