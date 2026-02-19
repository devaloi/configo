package configo

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Rule defines validation constraints for a config key.
type Rule struct {
	Required bool
	Min      *float64
	Max      *float64
	Regex    string
	Custom   func(value any) error
}

// Validate checks config values against the given rules.
// All errors are collected into a ValidationError.
func (c *Config) Validate(rules map[string]Rule) error {
	var errs []FieldError

	c.mu.RLock()
	data := c.data
	c.mu.RUnlock()

	for key, rule := range rules {
		val, ok := data[key]

		if rule.Required && !ok {
			errs = append(errs, FieldError{Field: key, Message: "required"})
			continue
		}
		if !ok {
			continue
		}

		if rule.Min != nil || rule.Max != nil {
			n, err := toFloat64(val)
			if err != nil {
				errs = append(errs, FieldError{Field: key, Message: fmt.Sprintf("cannot convert to number: %v", err)})
			} else {
				if rule.Min != nil && n < *rule.Min {
					errs = append(errs, FieldError{Field: key, Message: fmt.Sprintf("value %v is less than min %v", n, *rule.Min)})
				}
				if rule.Max != nil && n > *rule.Max {
					errs = append(errs, FieldError{Field: key, Message: fmt.Sprintf("value %v is greater than max %v", n, *rule.Max)})
				}
			}
		}

		if rule.Regex != "" {
			s := fmt.Sprintf("%v", val)
			re, err := regexp.Compile(rule.Regex)
			if err != nil {
				errs = append(errs, FieldError{Field: key, Message: fmt.Sprintf("invalid regex: %v", err)})
			} else if !re.MatchString(s) {
				errs = append(errs, FieldError{Field: key, Message: fmt.Sprintf("value %q does not match pattern %q", s, rule.Regex)})
			}
		}

		if rule.Custom != nil {
			if err := rule.Custom(val); err != nil {
				errs = append(errs, FieldError{Field: key, Message: err.Error()})
			}
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

// ValidateStruct validates a struct using `validate` tags.
// Supported tags: required, min=N, max=N, regex=PATTERN.
func (c *Config) ValidateStruct(target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validate: target must be a struct")
	}

	rules := make(map[string]Rule)
	buildRulesFromStruct(v.Type(), rules)

	return c.Validate(rules)
}

func buildRulesFromStruct(t reflect.Type, rules map[string]Rule) {
	for i := range t.NumField() {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			buildRulesFromStruct(field.Type, rules)
			continue
		}

		key := field.Tag.Get("config")
		validateTag := field.Tag.Get("validate")
		if key == "" || validateTag == "" {
			continue
		}

		rule := parseValidateTag(validateTag)
		rules[key] = rule
	}
}

func parseValidateTag(tag string) Rule {
	var rule Rule
	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		switch {
		case part == "required":
			rule.Required = true
		case strings.HasPrefix(part, "min="):
			if v, err := strconv.ParseFloat(strings.TrimPrefix(part, "min="), 64); err == nil {
				rule.Min = &v
			}
		case strings.HasPrefix(part, "max="):
			if v, err := strconv.ParseFloat(strings.TrimPrefix(part, "max="), 64); err == nil {
				rule.Max = &v
			}
		case strings.HasPrefix(part, "regex="):
			rule.Regex = strings.TrimPrefix(part, "regex=")
		}
	}
	return rule
}
