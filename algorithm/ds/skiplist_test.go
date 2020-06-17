package ds

import (
	"errors"
	"testing"
)

func TestCRUD(t *testing.T) {
	sk := NewSkipList()
	sk.Insert(5, "val5")
	sk.Insert(8, "val8")
	sk.Insert(13, "val13")
	sk.Insert(18, "val18")
	sk.Insert(33, "val33")
	sk.Insert(21, "val21")
	sk.Insert(4, "val4")

	//sk.PrintSkipList()

	if v, err := sk.Search(13); err != nil || v.(string) != "val13" {
		t.Errorf("searching for 13 should be val13, got:%v", v)
	}

	if v, err := sk.Search(88); !errors.Is(err, ErrKeyNotFound) || v != nil {
		t.Errorf("searching for 88 should get nothing, got val=%v, error=%v", v, err)
	}

	sk.Insert(33, "val33u")

	if v, err := sk.Search(33); err != nil || v.(string) != "val33u" {
		t.Errorf("val33 should be updated to val33u, got val=%v, error=%v", v, err)
	}

	if err := sk.Delete(21); err != nil {
		t.Errorf("delete 21 should get nil as return value, got:%v", err)
	}

	if err := sk.Delete(110); !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("delete 110 should get ErrKeyNotFound, got:%v", ErrKeyNotFound)
	}
}
