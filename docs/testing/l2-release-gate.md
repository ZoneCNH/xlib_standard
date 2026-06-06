# L2 Release Gate

xlib-standard is the standards source. It defines artifact shape, required evidence, and release semantics only. testkitx owns executable contract libraries, and xlibgate owns release adjudication. This repository must remain provider-neutral: no provider connections, no credentials loading, and no provider runner implementation live here.

Release levels are defined in `.agent/registry/l2-release-levels.yaml` and may not be redefined downstream. L2-T3 is the first release-allowed level; L2-T4 is factory-grade.

## Level expectations

- `L2-T0` proves skeleton readiness only.
- `L2-T1` requires unit and contract profiles.
- `L2-T2` adds integration evidence.
- `L2-T3` adds chaos, benchmark, and adoption evidence and allows release.
- `L2-T4` adds retrospective evidence and allows factory-grade claims.

## Gate behavior

Gates fail closed when required profiles, evidence paths, selected pack reports, compatibility notes, or the release-readiness summary are missing. xlibgate adjudicates the evidence against this registry; xlib-standard only defines the release semantics and template targets.

Release-readiness evidence is the active summary at `.agent/evidence/l2/release-readiness.json`. The template must not ship that file because inherited placeholder evidence would create a false release signal. Downstream release workflows generate it from real profile reports, then the release-readiness check validates the schema before a release claim proceeds.

## Template targets

`make l2-contract` runs the local Go manifest-shape contract test for the template. It proves the template declares L2 contract packs and remains provider-neutral; it does not connect to a provider.

`make l2-release-readiness-check` requires `.agent/evidence/l2/release-readiness.json` and validates it against `.agent/schemas/l2-release-readiness.schema.json`. The check fails when the active file is absent, so pull-request shape checks should not create fake release evidence just to pass. Release or tag workflows should invoke the check only after real downstream evidence has been generated.

`make l2-release-readiness` keeps the compatibility target name and delegates to the release-readiness check. The compatibility aliases `l2-manifest-check` and `l2-evidence-check` remain available for existing workflows.

Evidence paths for release claims should remain under `.agent/evidence/l2` unless xlib-standard changes the schema.
