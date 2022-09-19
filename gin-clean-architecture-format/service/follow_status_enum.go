package service

func GetStatusEnum() []int8 {
	return []int8{1, 2}
}

func Contains(array []int8, target int8) bool {
	flag := false
	for _, val := range array {
		if val == target {
			flag = true
		}
	}
	return flag
}
