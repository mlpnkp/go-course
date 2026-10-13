package main

import "strconv"

const (
	regularFrom = 10_000
	largeFrom   = 50_000
)

func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return strconv.Itoa(n)
	}
}

type Statistics struct {
	Min int64
	Max int64
	Sum int64
	Avg float64
}

func CalcStatistics(amounts []int64) (Statistics, bool) {
	if len(amounts) == 0 {
		return Statistics{}, false
	}

	s := Statistics{Min: amounts[0], Max: amounts[0]}
	for _, a := range amounts {
		if a < s.Min {
			s.Min = a
		}
		if a > s.Max {
			s.Max = a
		}
		s.Sum += a
	}
	s.Avg = float64(s.Sum) / float64(len(amounts))

	return s, true
}

func Category(amount int64) (string, bool) {
	switch {
	case amount < 0:
		return "", false
	case amount < regularFrom:
		return "small", true
	case amount < largeFrom:
		return "regular", true
	default:
		return "large", true
	}
}
