# L2 Rollout Playbook

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

## Sequence

1. Copy `templates/l2` into the L2 repository.
2. Fill adapter metadata and declared capabilities.
3. Select required contract packs from the registry.
4. Wire testkitx in the downstream repository.
5. Generate evidence and run xlibgate adjudication.
6. Record blockers instead of lowering release criteria.
