package rules

import (
	"fmt"
	"strings"
)

type EvaluationContext struct {
	Auth    map[string]any
	Data    map[string]any
	Record  map[string]any
	IsAdmin bool
}

func Evaluate(rule string, ctx EvaluationContext) (bool, error) {
	if ctx.IsAdmin {
		return true, nil // Admin bypass
	}
	if rule == "" {
		return true, nil // Empty rule means public access
	}

	// For Phase 5, we implement a simplified "String-based" evaluator.
	// In a full production version, this would be a real AST parser.
	// Here we handle the most common Vault use case: "field = value" or "field = @request.auth.id"

	// 1. Clean the rule
	rule = strings.TrimSpace(rule)

	// 2. Handle simple boolean strings
	if rule == "true" {
		return true, nil
	}
	if rule == "false" {
		return false, nil
	}

	// 3. Simple equality check (e.g., "id = @request.auth.id")
	if strings.Contains(rule, " = ") {
		parts := strings.Split(rule, " = ")
		if len(parts) == 2 {
			left := h_resolveValue(strings.TrimSpace(parts[0]), ctx)
			right := h_resolveValue(strings.TrimSpace(parts[1]), ctx)
			return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
		}
	}

	// 4. Simple inequality check (e.g., "id != @request.auth.id")
	if strings.Contains(rule, " != ") {
		parts := strings.Split(rule, " != ")
		if len(parts) == 2 {
			left := h_resolveValue(strings.TrimSpace(parts[0]), ctx)
			right := h_resolveValue(strings.TrimSpace(parts[1]), ctx)
			return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
		}
	}

	return false, fmt.Errorf("unsupported rule expression: %s", rule)
}

func h_resolveValue(key string, ctx EvaluationContext) any {
	// Handle literals
	if strings.HasPrefix(key, "'") && strings.HasSuffix(key, "'") {
		return strings.Trim(key, "'")
	}

	// Handle @request context
	if strings.HasPrefix(key, "@request.auth.") {
		field := strings.TrimPrefix(key, "@request.auth.")
		return ctx.Auth[field]
	}
	if strings.HasPrefix(key, "@request.data.") {
		field := strings.TrimPrefix(key, "@request.data.")
		return ctx.Data[field]
	}

	// Handle record context
	if strings.HasPrefix(key, "record.") {
		field := strings.TrimPrefix(key, "record.")
		return ctx.Record[field]
	}

	// Default to record field if no prefix
	return ctx.Record[key]
}
