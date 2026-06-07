# Copilot Review Instructions

This repository is a governed standard-source repository, not a business
application. Treat the project authority chain as:

1. `CONSTITUTION.md`
2. `.agent/rules/`
3. `.agent/harness/`
4. `contracts/`
5. `docs/architecture/`
6. `AGENTS.md`
7. Tool-specific files and temporary notes

When reviewing pull requests, prioritize blockers over style suggestions.

Block or flag changes that:

- expose secrets, tokens, private connection strings, or production credential
  paths;
- claim completion without `DONE with evidence:` or without matching gate output;
- weaken required harness gates, branch governance, evidence requirements, or
  layer boundaries;
- place business logic in Standard, L0, L1, or L2 libraries;
- add L2-to-L2 dependencies, reverse dependencies, or contract bypasses;
- change public APIs without contracts, examples, docs, tests, and release
  impact notes;
- commit generated artifacts listed in `.agent/registries/generated-artifacts.yaml`
  or `.agent/contracts/scope-locks.yaml`;
- modify GitHub workflows without pinning third-party actions to a 40-character
  commit SHA and preserving the source tag comment.

Do not recommend bypassing protected-branch rules, force-pushing `main`, deleting
protected branches, weakening status checks, or treating AI review comments as a
replacement for project gates and evidence.

For each finding, include severity, file, line, violated project rule, and the
minimal corrective action. If no blocker is found, summarize residual
verification gaps without claiming gates passed unless the logs prove them.
