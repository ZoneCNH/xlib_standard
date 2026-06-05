# L2 Adapter Testing Standard

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

## Required flow

1. Declare capabilities in `.agent/l2-capabilities.yaml`.
2. Map declared capabilities to xlib-standard contract packs.
3. Execute applicable testkitx profiles downstream.
4. Store evidence under `.agent/evidence/l2`.
5. Ask xlibgate to adjudicate compliance and release readiness.

## Profiles

`skeleton`, `unit`, `contract`, `integration`, `chaos`, `benchmark`, `adoption`, and `retrospective` are release-level inputs.
