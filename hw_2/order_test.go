package main

import (
	"errors"
	"testing"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name    string
		id      OrderID
		wantErr error
	}{
		{name: "positive id", id: 42},
		{name: "zero id", id: 0, wantErr: ErrInvalidOrderID},
		{name: "negative id", id: -1, wantErr: ErrInvalidOrderID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewOrder(%d) error = %v, want %v", tt.id, err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if order.ID != tt.id || order.Status != StatusNew || len(order.Items) != 0 {
				t.Errorf("NewOrder(%d) = %+v, want new order without items", tt.id, order)
			}
		})
	}
}

func TestTotalOfEmptyOrderIsZero(t *testing.T) {
	order := mustNewOrder(t, 1)

	if got := Total(order); got != 0 {
		t.Errorf("Total() = %d, want 0", got)
	}
}

func TestTotalOfTwoItemsIsTheirSum(t *testing.T) {
	order := mustNewOrder(t, 1)
	mustAddItem(t, &order, Item{Name: "Кофе", PriceKopecks: 19_990})
	mustAddItem(t, &order, Item{Name: "Фильтр", PriceKopecks: 2_500})

	if got, want := Total(order), int64(19_990+2_500); got != want {
		t.Errorf("Total() = %d, want %d", got, want)
	}
}

func TestAddItemRejectsInvalidInputWithoutPanic(t *testing.T) {
	valid := mustNewOrder(t, 1)
	tests := []struct {
		name    string
		order   *Order
		item    Item
		wantErr error
	}{
		{name: "nil order", order: nil, item: Item{Name: "Кофе", PriceKopecks: 100}, wantErr: ErrNilOrder},
		{name: "empty name", order: &valid, item: Item{Name: "", PriceKopecks: 100}, wantErr: ErrEmptyItemName},
		{name: "blank name", order: &valid, item: Item{Name: "   ", PriceKopecks: 100}, wantErr: ErrEmptyItemName},
		{name: "negative price", order: &valid, item: Item{Name: "Кофе", PriceKopecks: -1}, wantErr: ErrNegativePrice},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("AddItem panicked: %v", r)
				}
			}()

			err := AddItem(tt.order, tt.item)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddItem() error = %v, want %v", err, tt.wantErr)
			}
			if tt.order != nil && len(tt.order.Items) != 0 {
				t.Errorf("rejected item was added: %+v", tt.order.Items)
			}
		})
	}
}

func TestAddItemAcceptsZeroPrice(t *testing.T) {
	order := mustNewOrder(t, 1)

	mustAddItem(t, &order, Item{Name: "Подарок", PriceKopecks: 0})

	if len(order.Items) != 1 {
		t.Errorf("len(Items) = %d, want 1", len(order.Items))
	}
}

func TestChangingOrderCopyKeepsOriginal(t *testing.T) {
	original := mustNewOrder(t, 1)
	mustAddItem(t, &original, Item{Name: "Кофе", PriceKopecks: 19_990})

	cp := original
	cp.ID = 2
	cp.Status = StatusPaid

	if original.ID != 1 || original.Status != StatusNew {
		t.Errorf("original changed through copy: %+v", original)
	}
}

func TestPlainCopySharesItemsButCloneDoesNot(t *testing.T) {
	original := mustNewOrder(t, 1)
	mustAddItem(t, &original, Item{Name: "Кофе", PriceKopecks: 19_990})

	shared := original
	shared.Items[0].Name = "Чай"
	if original.Items[0].Name != "Чай" {
		t.Fatalf("plain copy should share Items backing array, got %q", original.Items[0].Name)
	}

	clone := Clone(original)
	clone.Items[0].Name = "Какао"
	mustAddItem(t, &clone, Item{Name: "Фильтр", PriceKopecks: 2_500})

	if original.Items[0].Name != "Чай" || len(original.Items) != 1 {
		t.Errorf("original changed through clone: %+v", original.Items)
	}
}

func TestAddItemGrowsSlice(t *testing.T) {
	order := mustNewOrder(t, 1)

	for i := range 5 {
		mustAddItem(t, &order, Item{Name: "Позиция", PriceKopecks: int64(i)})
		t.Logf("after %d: len=%d cap=%d", i+1, len(order.Items), cap(order.Items))
	}

	if len(order.Items) != 5 || cap(order.Items) < 5 {
		t.Errorf("len=%d cap=%d, want len 5 and cap >= 5", len(order.Items), cap(order.Items))
	}
}

func TestStatusNames(t *testing.T) {
	names := StatusNames()

	for _, status := range []string{StatusNew, StatusPaid, StatusShipped, StatusCanceled} {
		if names[status] == "" {
			t.Errorf("no human-readable name for status %q", status)
		}
	}
}

func TestWriteToNilMapPanicsAndMakeFixesIt(t *testing.T) {
	writeTo := func(m map[string]string) (panicked bool) {
		defer func() { panicked = recover() != nil }()
		m[StatusNew] = "новый"
		return false
	}

	var nilMap map[string]string
	if !writeTo(nilMap) {
		t.Fatal("write to nil map should panic")
	}

	made := make(map[string]string)
	if writeTo(made) {
		t.Fatal("write to map from make should not panic")
	}
	if made[StatusNew] != "новый" {
		t.Errorf("made[%q] = %q, want %q", StatusNew, made[StatusNew], "новый")
	}
}

func mustNewOrder(t *testing.T, id OrderID) Order {
	t.Helper()
	order, err := NewOrder(id)
	if err != nil {
		t.Fatalf("NewOrder(%d): %v", id, err)
	}
	return order
}

func mustAddItem(t *testing.T, order *Order, item Item) {
	t.Helper()
	if err := AddItem(order, item); err != nil {
		t.Fatalf("AddItem(%+v): %v", item, err)
	}
}
