package main

import "fmt"

func main() {
	fmt.Println("1. FizzBuzz")
	for i := 1; i <= 30; i++ {
		fmt.Println(FizzBuzz(i))
	}

	fmt.Println("\n2. Статистика по суммам")
	for _, amounts := range [][]int64{{12000, 7500, 31000, 4500}, {}} {
		s, ok := CalcStatistics(amounts)
		if !ok {
			fmt.Printf("%v: пустой срез, решила что статистику посчитать нельзя :)\n", amounts)
			continue
		}
		fmt.Printf("%v: min=%d max=%d sum=%d avg=%g\n", amounts, s.Min, s.Max, s.Sum, s.Avg)
	}

	fmt.Println("\n3. Категория заказа")
	for _, amount := range []int64{0, 9_999, 10_000, 49_999, 50_000, -1} {
		category, ok := Category(amount)
		if !ok {
			fmt.Printf("%d: ошибка ввода, сумма не может быть отрицательной\n", amount)
			return
		}
		fmt.Printf("%d: %s\n", amount, category)
	}
}
