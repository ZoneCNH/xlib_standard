# L2 Release Gate

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

Release levels are defined in `.agent/registry/l2-release-levels.yaml` and may not be redefined downstream. L2-T3 is the first release-allowed level; L2-T4 is factory-grade.

Downstream gates should fail closed when required profiles or evidence are missing.
