package detection

import (
	"regexp"
	"safectx/pkg/schema"
)

// BlockedPatterns contains regex patterns for detecting malicious prompts
var BlockedPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)drop\s+table`),
	regexp.MustCompile(`(?i)shutdown`),
	regexp.MustCompile(`(?i)delete\s+from`),
	regexp.MustCompile(`(?i)system\s+command`),
	regexp.MustCompile(`(?i)execute\s+shell`),
}

// CheckForInjection checks if the request contains any blocked patterns
func CheckForInjection(req *schema.MCPRequest) bool {
	if req.Params == nil {
		return false
	}

	// Convert params to string for pattern matching
	paramsStr := req.Params["prompt"]
	if prompt, ok := paramsStr.(string); ok {
		for _, pattern := range BlockedPatterns {
			if pattern.MatchString(prompt) {
				return true
			}
		}
	}

	return false
}
