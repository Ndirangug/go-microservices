package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "George",
		Price: 1,
		SKU:   "abc-defg-hijk",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
