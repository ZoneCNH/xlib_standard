# L2 Downstream Adoption

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

Adopting repositories consume xlib-standard templates and registries without redefining release levels. They may add provider-specific services, credentials loading, and Contract Runner wiring only in their own repository.

The included template starts with shape checks and placeholder tests to keep adoption safe.
