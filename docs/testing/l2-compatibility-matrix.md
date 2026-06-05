# L2 Compatibility Matrix

xlib-standard is the standards source. It defines shape, evidence, and release semantics only; testkitx executes contract libraries and xlibgate adjudicates gates. No provider credentials, provider endpoints, or Contract Runner implementation belong in this repository.

Compatibility should be recorded against adapter family, selected contract packs, test profiles, and target release level. The golden sample registry `.agent/registry/l2-golden-samples.yaml` lists initial standards-only expectations for redisx, postgresx, kafkax, natsx, ossx, clickhousex, and taosx.
