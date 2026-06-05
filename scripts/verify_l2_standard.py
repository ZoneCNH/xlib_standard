#!/usr/bin/env python3
"""Validate L2 standards artifacts without provider connections."""
from __future__ import annotations
import json
from pathlib import Path
import sys
import yaml
from jsonschema import Draft202012Validator

ROOT = Path(__file__).resolve().parents[1]
EVIDENCE = ROOT / ".agent" / "evidence" / "l2-standard"

REGISTRIES = {
    ".agent/registry/l2-contract-packs.yaml": ".agent/schemas/l2-contract-packs.schema.json",
    ".agent/registry/l2-capability-families.yaml": None,
    ".agent/registry/l2-golden-samples.yaml": None,
    ".agent/registry/l2-release-levels.yaml": None,
}
SCHEMAS = [
    ".agent/schemas/l2-capabilities.schema.json",
    ".agent/schemas/l2-contract-packs.schema.json",
    ".agent/schemas/l2-release-readiness.schema.json",
    ".agent/schemas/l2-compliance-matrix.schema.json",
]
TEMPLATE_REQUIRED = [
    "templates/l2/.agent/l2-capabilities.yaml",
    "templates/l2/.agent/gates/l2gate.yaml",
    "templates/l2/.agent/evidence/README.md",
    "templates/l2/.agent/evidence/l2/.gitkeep",
    "templates/l2/test/contract/l2_contract_test.go",
    "templates/l2/test/integration/README.md",
    "templates/l2/test/chaos/README.md",
    "templates/l2/test/benchmark/README.md",
    "templates/l2/test/adoption/README.md",
    "templates/l2/docker-compose.test.yml",
    "templates/l2/Makefile",
    "templates/l2/.github/workflows/l2-gates.yml",
]
DOC_REQUIRED = [
    "docs/testing/l2-adapter-testing-standard.md",
    "docs/testing/l2-capability-manifest.md",
    "docs/testing/l2-contract-pack-registry.md",
    "docs/testing/l2-evidence-standard.md",
    "docs/testing/l2-release-gate.md",
    "docs/testing/l2-compliance-matrix.md",
    "docs/testing/l2-rollout-playbook.md",
    "docs/testing/l2-downstream-adoption.md",
    "docs/testing/l2-compatibility-matrix.md",
]
FORBIDDEN_TERMS = ["provider_endpoint:", "provider_credentials:", "password:", "secret:", "token:"]

def load_yaml(rel: str):
    with (ROOT / rel).open() as f:
        return yaml.safe_load(f)

def load_json(rel: str):
    with (ROOT / rel).open() as f:
        return json.load(f)

def result(name, status, details):
    return {"check": name, "status": status, "details": details}

def main() -> int:
    EVIDENCE.mkdir(parents=True, exist_ok=True)
    results = []

    parsed_registries = {}
    for rel, schema_rel in REGISTRIES.items():
        data = load_yaml(rel)
        parsed_registries[rel] = data
        if schema_rel:
            Draft202012Validator(load_json(schema_rel)).validate(data)
        results.append(result(rel, "PASS", "YAML parsed" + (" and schema validated" if schema_rel else "")))

    for rel in SCHEMAS:
        schema = load_json(rel)
        Draft202012Validator.check_schema(schema)
        results.append(result(rel, "PASS", "JSON schema parsed and meta-schema checked"))

    manifest = load_yaml("templates/l2/.agent/l2-capabilities.yaml")
    Draft202012Validator(load_json(".agent/schemas/l2-capabilities.schema.json")).validate(manifest)
    results.append(result("templates/l2/.agent/l2-capabilities.yaml", "PASS", "template manifest validates against capability schema"))

    missing = [rel for rel in TEMPLATE_REQUIRED + DOC_REQUIRED if not (ROOT / rel).exists()]
    if missing:
        raise AssertionError(f"missing required artifacts: {missing}")
    results.append(result("required-artifacts", "PASS", f"{len(TEMPLATE_REQUIRED)} template files and {len(DOC_REQUIRED)} docs present"))

    offenders = []
    for rel in TEMPLATE_REQUIRED + DOC_REQUIRED + list(REGISTRIES) + SCHEMAS:
        text = (ROOT / rel).read_text()
        for term in FORBIDDEN_TERMS:
            if term in text:
                offenders.append({"path": rel, "term": term})
    if offenders:
        raise AssertionError(f"forbidden provider/secret terms found: {offenders}")
    results.append(result("provider-boundary", "PASS", "no provider credentials/endpoints or secret fields in standards artifacts"))

    packs = parsed_registries[".agent/registry/l2-contract-packs.yaml"]["packs"]
    levels = parsed_registries[".agent/registry/l2-release-levels.yaml"]["levels"]
    results.append(result("registry-summary", "PASS", {"contract_packs": sorted(packs), "release_levels": list(levels)}))

    (EVIDENCE / "schema-validate.json").write_text(json.dumps({"status": "PASS", "checks": [r for r in results if "schema" in r["details"] or "schema" in r["check"]]}, indent=2) + "\n")
    (EVIDENCE / "registry-check.json").write_text(json.dumps({"status": "PASS", "checks": [r for r in results if "registry" in r["check"]]}, indent=2) + "\n")
    (EVIDENCE / "template-check.json").write_text(json.dumps({"status": "PASS", "checks": [r for r in results if "template" in r["check"] or r["check"] in ["required-artifacts", "provider-boundary"]]}, indent=2) + "\n")
    (EVIDENCE / "verification-summary.json").write_text(json.dumps({"status": "PASS", "checks": results}, indent=2) + "\n")
    print(json.dumps({"status": "PASS", "checks": len(results), "evidence_dir": str(EVIDENCE.relative_to(ROOT))}, indent=2))
    return 0

if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except Exception as exc:
        EVIDENCE.mkdir(parents=True, exist_ok=True)
        (EVIDENCE / "verification-summary.json").write_text(json.dumps({"status": "FAIL", "error": str(exc)}, indent=2) + "\n")
        print(f"FAIL: {exc}", file=sys.stderr)
        raise
