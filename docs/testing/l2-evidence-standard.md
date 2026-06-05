# L2 Evidence Standard

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

Evidence must be reproducible, machine-readable where possible, and linked to manifest capabilities and contract packs. Minimum downstream evidence includes contract reports, compliance matrix, and release readiness.

Evidence paths should stay under `.agent/evidence/l2` in adopting L2 repositories.
