package data

import "testing"

func TestProduct_Validate(t *testing.T) {
	p := &Product{
		Name:  "Junk",
		Price: 1.00,
		SKU:   "abc-abc-abc",
	}

	err := p.Validate()

	if err != nil {
		t.Fatalf("[ERROR] Can't validate SKU : %s", err)
	}

}
