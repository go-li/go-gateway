package main

func osbanned(code string) bool {
	for i := 0; i+2 < len(code); i++ {
		if (code[i] == '"') && (code[i+1] == 'o') && (code[i+2] == 's') {
			return true
		}

	}
	return false
}
