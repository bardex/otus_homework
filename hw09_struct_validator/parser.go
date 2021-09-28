package hw09structvalidator

import (
	"fmt"
	"strings"
	"sync"
)

type validationRule struct {
	Rule   string
	Params []string
}

type validationRules []validationRule

// use lru cache in production.
var (
	parserCache      = map[string]validationRules{}
	parserCacheMutex sync.Mutex
)

func parseValidateTag(tag string) (validationRules, error) {
	if tag == "" {
		return validationRules{}, nil
	}

	parserCacheMutex.Lock()
	defer parserCacheMutex.Unlock()

	if cache, exists := parserCache[tag]; exists {
		return cache, nil
	}

	rules := strings.Split(tag, "|")
	results := make(validationRules, 0, len(rules))
	for j := range rules {
		parts := strings.SplitN(strings.TrimSpace(rules[j]), ":", 2)
		rule := strings.TrimSpace(parts[0])
		if rule == "" {
			return nil, fmt.Errorf("invalid validate tag `%s`", tag)
		}
		params := make([]string, 0)
		if len(parts) == 2 {
			parts[1] = strings.TrimSpace(parts[1])
			if parts[1] == "" {
				return nil, fmt.Errorf("invalid validate tag `%s`", tag)
			}
			params = strings.Split(parts[1], ",")
			for i := range params {
				params[i] = strings.TrimSpace(params[i])
			}
		}
		results = append(results, validationRule{Rule: rule, Params: params})
	}
	parserCache[tag] = results
	return results, nil
}
