package util

func ParseReviewFrequencyFormula(s string) ([]int, bool) {
	if len(s) > 10000 {
		return nil, false
	}
	arr := make([]int, 0)
	for i := 0; i < len(s); i++ {
		j := i
		num := 0
		for j < len(s) && s[j] != '_' {
			if s[j] < '0' || s[j] > '9' {
				return nil, false
			}
			num = num*10 + int(s[j]-'0')
			j++
		}
		if num <= 0 {
			return nil, false
		}
		arr = append(arr, num)
		j++
		i = j - 1
	}
	return arr, true
}
