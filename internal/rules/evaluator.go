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

	// Clean the rule
	rule = strings.TrimSpace(rule)
	if rule == "true" {
		return true, nil
	}
	if rule == "false" {
		return false, nil
	}

	l := NewLexer(rule)
	p := NewParser(l)
	node, err := p.Parse()
	if err != nil {
		return false, fmt.Errorf("parse error: %w", err)
	}

	result, err := evalNode(node, ctx)
	if err != nil {
		return false, fmt.Errorf("eval error: %w", err)
	}

	if boolResult, ok := result.(bool); ok {
		return boolResult, nil
	}
	return false, fmt.Errorf("rule did not evaluate to a boolean")
}

func evalNode(node Node, ctx EvaluationContext) (any, error) {
	switch n := node.(type) {
	case *BooleanLiteral:
		return n.Value, nil
	case *IntegerLiteral:
		return n.Value, nil
	case *StringLiteral:
		return n.Value, nil
	case *Identifier:
		return resolveValue(n.Value, ctx), nil
	case *InfixExpression:
		left, err := evalNode(n.Left, ctx)
		if err != nil {
			return nil, err
		}
		right, err := evalNode(n.Right, ctx)
		if err != nil {
			return nil, err
		}
		return applyOp(n.Operator, left, right)
	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

func resolveValue(key string, ctx EvaluationContext) any {
	// Handle @request context
	if strings.HasPrefix(key, "@request.auth.") {
		field := strings.TrimPrefix(key, "@request.auth.")
		if ctx.Auth == nil {
			return nil
		}
		return ctx.Auth[field]
	}
	if strings.HasPrefix(key, "@request.data.") {
		field := strings.TrimPrefix(key, "@request.data.")
		if ctx.Data == nil {
			return nil
		}
		return ctx.Data[field]
	}

	// Handle record context
	if strings.HasPrefix(key, "record.") {
		field := strings.TrimPrefix(key, "record.")
		if ctx.Record == nil {
			return nil
		}
		return ctx.Record[field]
	}

	// Default to record field if no prefix
	if ctx.Record != nil {
		if val, ok := ctx.Record[key]; ok {
			return val
		}
	}
	return nil
}

func applyOp(op string, left, right any) (any, error) {
	switch op {
	case "=", "==":
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil
	case "!=":
		return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right), nil
	case "&&":
		l, ok1 := left.(bool)
		r, ok2 := right.(bool)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("type mismatch for &&")
		}
		return l && r, nil
	case "||":
		l, ok1 := left.(bool)
		r, ok2 := right.(bool)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("type mismatch for ||")
		}
		return l || r, nil
	case ">":
		return compareNumbers(left, right, func(l, r float64) bool { return l > r })
	case "<":
		return compareNumbers(left, right, func(l, r float64) bool { return l < r })
	case ">=":
		return compareNumbers(left, right, func(l, r float64) bool { return l >= r })
	case "<=":
		return compareNumbers(left, right, func(l, r float64) bool { return l <= r })
	}
	return nil, fmt.Errorf("unknown operator: %s", op)
}

func compareNumbers(left, right any, cmp func(float64, float64) bool) (bool, error) {
	l, err := toFloat(left)
	if err != nil {
		return false, err
	}
	r, err := toFloat(right)
	if err != nil {
		return false, err
	}
	return cmp(l, r), nil
}

func toFloat(v any) (float64, error) {
	switch val := v.(type) {
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("cannot convert %v to number", v)
	}
}
