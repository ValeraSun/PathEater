package geometry

import "math"

//модуль числа
// func abs(f float64) float64 {
// 	if f < 0 {
// 		return -f
// 	}
// 	return f
// }

//минимальное из двух чисел
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

//максимальное из двух чисел
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

//число в диапазоне от lo до hi
func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

//находит синус с помощью косинуса
func cosToSin(cos float64) float64 {
	return math.Sqrt(1 - cos*cos)
}

func average(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	var sum float64 = 0
	for _, num := range nums {
		sum += num
	}

	return sum / float64(len(nums))
}
