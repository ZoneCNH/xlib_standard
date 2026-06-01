# Evidence Protocol

Every completion claim must include `DONE with evidence:` and list:

- Goal/REQ ids covered.
- Commands run, with PASS/FAIL and short output.
- Artifacts produced: release manifest, checksum, docs, ADRs, matrix, review notes.
- `artifact_url`, `workflow_run_id`, `sha256`, commit, tree SHA and version when available.
- Known gaps with owning worker if a gate belongs outside the current slice.

Never include production secrets or `/home/k8s/secrets/env/*` contents in Evidence.
