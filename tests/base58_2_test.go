package tests

import (
	"testing"

	base58 "github.com/VSmert/base58-for-tinygo"
)

func TestBase58_test2(t *testing.T) {

	testAddr := []string{
		"1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojq",
		"1DhRmSGnhPjUaVPAj48zgPV9e2oRhAQFUb",
		"17LN2oPYRYsXS9TdYdXCCDvF2FegshLDU2",
		"14h2bDLZSuvRFhUL45VjPHJcW667mmRAAn",
	}

	for ii, vv := range testAddr {
		// num := Base58Decode([]byte(vv))
		// chk := Base58Encode(num)
		num, err := base58.Decode(vv)
		if err != nil {
			t.Errorf("Test %d, expected success, got error %s\n", ii, err)
		}
		chk := base58.Encode(num)
		if vv != string(chk) {
			t.Errorf("Test %d, expected=%s got=%s Address did base58 encode/decode correctly.", ii, vv, chk)
		}
	}
}
