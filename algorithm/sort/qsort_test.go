package sort

import "testing"

func TestQuickSort1(t *testing.T) {
	data := []int{8, 6, 9, 1, 22}
	QuickSort(data)

	if data[0] != 1 || data[1] != 6 || data[2] != 8 || data[3] != 9 ||
		data[4] != 22 {
		t.Error("QuickSort failed, Got: ", data, "expected: 1, 6, 8, 9, 22")
	}
}

func TestQuickSort2(t *testing.T) {
	data := []int{4, 6, 1, 8, 6}
	QuickSort(data)

	if data[0] != 1 || data[1] != 4 || data[2] != 6 || data[3] != 6 ||
		data[4] != 8 {
		t.Error("QuickSort failed, Got:", data, "expected: 1, 4, 6, 6, 8")
	}
}

func TestQuickSort3(t *testing.T) {
	data := []int{87}
	QuickSort(data)

	if data[0] != 87 {
		t.Error("QuickSort failed, Got:", data, "expected: 87")
	}
}
