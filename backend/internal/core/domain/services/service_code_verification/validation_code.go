package service_code_verification

import "fmt"

type CodeValidator struct{}

func (c CodeValidator) Validate(code interface{}) error {
	codeStr, ok := code.(string)
	if !ok {
		return fmt.Errorf("code must be a string")
	}

	if codeStr == "" {
		return fmt.Errorf("code cannot be empty")
	}

	if len(codeStr) < 6 {
		return fmt.Errorf("code cannot have less than 6 digits")
	}

	if len(codeStr) > 6 {
		return fmt.Errorf("code cannot have more than 6 digits")
	}

	if codeStr[0] == '-' {
		return fmt.Errorf("code cannot be negative")
	}

	for _, char := range codeStr {
		if char < '0' || char > '9' {
			return fmt.Errorf("code can only contain digits")
		}
	}

	return nil
}