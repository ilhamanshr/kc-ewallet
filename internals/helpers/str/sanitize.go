package strhelper

func UnstyleAlphabets(input string) string {
	// Convert the input string to a rune slice
	runes := []rune(input)

	// Remove styling by setting the character to its unstyled equivalent
	// For example, convert uppercase letters to lowercase
	for i := range runes {
		runes[i] = ToUnstyled(runes[i])
	}

	// Convert the modified rune slice back to a string
	return string(runes)
}

func ToUnstyled(char rune) rune {
	// Define your unstyle logic here
	// For example, normalize italic
	switch {
	case '𝘈' <= char && char <= '𝘡':
		return char + ('A' - '𝘈')
	case '𝘢' <= char && char <= '𝘻':
		return char + ('a' - '𝘢')
	}
	return char
}
