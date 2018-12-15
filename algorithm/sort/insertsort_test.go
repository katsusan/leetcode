package sort

import "testing"

func TestInsertSort1(t *testing.T) {
	numbers := []int{4, 5, 9, 1, 33}
	InsertSort(numbers)
	if numbers[0] != 1 || numbers[1] != 4 || numbers[2] != 5 ||
		numbers[3] != 9 || numbers[4] != 33 {
		t.Error("InsertSort failed, Got ", numbers, " Expected: 1, 4, 5, 9, 33")
	}
}

func TestInsertSort2(t *testing.T) {
	numbers := []int{6, 9, 6, 32, 9}
	InsertSort(numbers)
	if numbers[0] != 6 || numbers[1] != 6 || numbers[2] != 9 ||
		numbers[3] != 9 || numbers[4] != 32 {
		t.Error("InsertSort failed, Got ", numbers, " Expected: 6, 6, 9, 9, 32")
	}
}

func TestInsertSort3(t *testing.T) {
	numbers := []int{3}
	InsertSort(numbers)
	if numbers[0] != 3 {
		t.Error("InsertSort failed, Got ", numbers, " Expected: 3")
	}
}
