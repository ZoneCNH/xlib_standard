# L2 Contract Pack Registry

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

The registry is `.agent/registry/l2-contract-packs.yaml` and is validated by `.agent/schemas/l2-contract-packs.schema.json`. Pack entries define family, title, required profiles, required evidence, and capability names.

Pack execution code is owned by testkitx, not xlib-standard.
