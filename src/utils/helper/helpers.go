package helper

func StringInSlice(val string, s []string) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}

func UintInSlice(val uint, s []uint) bool {
	for _, v := range s {
		if val == v {
			return true
		}
	}
	return false
}
