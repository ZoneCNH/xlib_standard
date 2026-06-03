package validation

import (
	"fmt"
	"strings"
)

func RequireNonEmpty(field string, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

type runtimeFileOwner struct {
	path           string
	owner          string
	reviewRequired string
	rationale      string
}

// ValidateRuntimeFileOwnership checks the local .agent file-ownership index for
// the control-plane classifications goalcli depends on before allowing the
// runtime-file-ownership gate to pass.
func ValidateRuntimeFileOwnership(path string, content string) []string {
	var gaps []string
	if strings.TrimSpace(content) == "" {
		return []string{path + " must not be empty"}
	}
	if !containsYAMLKey(content, "schema_version") {
		gaps = append(gaps, path+" missing schema_version")
	}
	if !containsYAMLKey(content, "owners") {
		gaps = append(gaps, path+" missing owners")
	}

	owners := parseRuntimeFileOwners(content)
	if len(owners) == 0 {
		gaps = append(gaps, path+" owners must include at least one path classification")
	}

	seen := map[string]bool{}
	byPath := map[string]runtimeFileOwner{}
	for _, owner := range owners {
		if seen[owner.path] {
			gaps = append(gaps, path+" duplicate owner entry "+owner.path)
			continue
		}
		seen[owner.path] = true
		byPath[owner.path] = owner
		if owner.owner == "" {
			gaps = append(gaps, path+" "+owner.path+" missing owner")
		}
		if owner.reviewRequired == "" {
			gaps = append(gaps, path+" "+owner.path+" missing review_required")
		} else if owner.reviewRequired != "true" && owner.reviewRequired != "false" {
			gaps = append(gaps, path+" "+owner.path+" review_required must be true or false")
		}
		if owner.rationale == "" {
			gaps = append(gaps, path+" "+owner.path+" missing rationale")
		}
	}

	requireRuntimeOwner(path, byPath, ".agent/", "governance", true, &gaps)
	requireRuntimeOwner(path, byPath, "cmd/goalcli/", "gate-runtime", true, &gaps)
	requireRuntimeOwner(path, byPath, "contracts/", "standard", true, &gaps)
	return gaps
}

func requireRuntimeOwner(path string, owners map[string]runtimeFileOwner, ownerPath string, ownerName string, reviewRequired bool, gaps *[]string) {
	owner, ok := owners[ownerPath]
	if !ok {
		*gaps = append(*gaps, path+" owners must include "+ownerPath)
		return
	}
	if owner.owner != ownerName {
		*gaps = append(*gaps, path+" "+ownerPath+" owner must be "+ownerName)
	}
	wantReviewRequired := fmt.Sprintf("%t", reviewRequired)
	if owner.reviewRequired != wantReviewRequired {
		*gaps = append(*gaps, path+" "+ownerPath+" review_required must be "+wantReviewRequired)
	}
}

func containsYAMLKey(content string, key string) bool {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(stripInlineYAMLComment(line))
		if trimmed == key+":" || strings.HasPrefix(trimmed, key+": ") {
			return true
		}
	}
	return false
}

func parseRuntimeFileOwners(content string) []runtimeFileOwner {
	var owners []runtimeFileOwner
	var current *runtimeFileOwner
	inOwners := false

	flush := func() {
		if current != nil {
			owners = append(owners, *current)
			current = nil
		}
	}

	for _, rawLine := range strings.Split(content, "\n") {
		line := stripInlineYAMLComment(rawLine)
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " "))
		trimmed := strings.TrimSpace(line)
		if !inOwners {
			if trimmed == "owners:" {
				inOwners = true
			}
			continue
		}
		if indent == 0 {
			flush()
			break
		}
		if indent == 2 && strings.HasSuffix(trimmed, ":") {
			flush()
			current = &runtimeFileOwner{path: unquoteYAMLScalar(strings.TrimSuffix(trimmed, ":"))}
			continue
		}
		if current == nil || indent < 4 {
			continue
		}
		field, value, ok := strings.Cut(trimmed, ":")
		if !ok {
			continue
		}
		value = unquoteYAMLScalar(strings.TrimSpace(value))
		switch strings.TrimSpace(field) {
		case "owner":
			current.owner = value
		case "review_required":
			current.reviewRequired = value
		case "rationale":
			current.rationale = value
		}
	}
	flush()
	return owners
}

func stripInlineYAMLComment(line string) string {
	if before, _, ok := strings.Cut(line, "#"); ok {
		return before
	}
	return line
}

func unquoteYAMLScalar(value string) string {
	return strings.Trim(strings.TrimSpace(value), `"'`)
}
