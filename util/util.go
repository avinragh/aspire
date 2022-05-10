package util

func Contains(inputSlice []int64, searchElem int64) bool {
	for _, input := range inputSlice {
		if input == searchElem {
			return true
		}
	}
	return false
}
