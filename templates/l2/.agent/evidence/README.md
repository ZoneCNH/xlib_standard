# L2 Evidence Directory

Downstream L2 repositories store generated contract, integration, chaos, benchmark, adoption, compliance, and release-readiness evidence here.

This template intentionally contains no provider credentials, no provider endpoints, and no active release evidence. Runtime execution belongs to downstream repositories using testkitx for executable suites and xlibgate for release adjudication.

## Release-readiness evidence

`make l2-release-readiness-check` is fail-closed: it requires `.agent/evidence/l2/release-readiness.json` and validates that file against the xlib-standard release-readiness schema. The template ships only `.gitkeep` files so a new repository cannot pass a release gate by inheriting placeholder evidence.

Generate `.agent/evidence/l2/release-readiness.json` only from a downstream release or adjudication workflow that has real profile reports to reference. Temporary fixtures may be used by standards verification, but they must be schema-valid, provider-neutral, and cleaned up before committing.

Do not commit fabricated active release-readiness evidence to satisfy the gate. A release claim should include real machine-readable reports under `.agent/evidence/l2` and should remain provider-neutral unless the downstream repository documents provider-specific execution outside xlib-standard.
