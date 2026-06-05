# L2 Capability Manifest

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

The manifest schema is `.agent/schemas/l2-capabilities.schema.json`; the downstream template is `templates/l2/.agent/l2-capabilities.yaml`. Required sections are `schema_version`, `layer`, `adapter`, `capabilities`, `contract_packs`, and `evidence`.

Forbidden manifest content includes credentials, provider endpoints, secrets, passwords, and tokens.
