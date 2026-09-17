package core

func FindAccessToken(value interface{}) (string, bool) {
	switch current := value.(type) {
	case map[string]interface{}:
		for key, nestedValue := range current {
			if key == "accessToken" {
				token, ok := nestedValue.(string)
				return token, ok
			}

			if token, found := FindAccessToken(nestedValue); found {
				return token, true
			}
		}

	case []interface{}:
		for _, item := range current {
			if token, found := FindAccessToken(item); found {
				return token, true
			}
		}
	}

	return "", false
}
