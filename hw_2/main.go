package main

import "fmt"

func main() {
	statusNames := StatusNames()

	order, err := NewOrder(42)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}

	items := []Item{
		{Name: "Кофе", PriceKopecks: 19_990},
		{Name: "Фильтр", PriceKopecks: 2_500},
		{Name: "", PriceKopecks: 100},
		{Name: "Сахар", PriceKopecks: -1},
	}
	for _, item := range items {
		if err := AddItem(&order, item); err != nil {
			fmt.Println("ошибка:", err)
			continue
		}
		fmt.Printf("добавлено %q: len=%d cap=%d\n", item.Name, len(order.Items), cap(order.Items))
	}

	if err := AddItem(nil, items[0]); err != nil {
		fmt.Println("ошибка:", err)
	}

	if _, err := NewOrder(0); err != nil {
		fmt.Println("ошибка:", err)
	}

	fmt.Printf("заказ %d, статус %s, сумма %d коп.\n",
		order.ID, statusNames[order.Status], Total(order))
}
