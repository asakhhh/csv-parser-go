package csvparser

func join(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	res := s[0]
	for _, str := range s[1:] {
		res += sep + str
	}
	return res
}

func contains(text, key string) bool {
	if len(key) == 0 {
		return true
	}
	for i := 0; i+len(key) <= len(text); i++ {
		if text[i] == key[0] {
			if text[i:i+len(key)] == key {
				return true
			}
		}
	}
	return false
}
