# fail-version-drift

Negative fixture for `goalcli fact audit --root <fixture>`.

The fixture keeps every canonical fact aligned with the repository defaults except
`current_release.version`, which is intentionally stale. The audit must reject it
with a `current_release.version drift` gap.
