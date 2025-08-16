package gs1sscc

import "fmt"

func Sscc(code string) (out string) {
	if len(code) != 17 {
		return "wrong length code"
	}
	// Validate that code contains only digits
	for i, ch := range code {
		if ch < '0' || ch > '9' {
			return fmt.Sprintf("invalid character '%c' at position %d", ch, i)
		}
	}
	sum := 0
	for i := range code {
		n := code[i] - '0'
		if i%2 == 0 {
			n *= 3
			sum += int(n)
		} else {
			sum += int(n)
		}
	}
	return fmt.Sprintf("%s%d", code, roundUp(sum)-sum)
}

func roundUp(val int) int {
	return 10 * ((val + 9) / 10)
}
