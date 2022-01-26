package tests

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	base58 "github.com/VSmert/base58-for-tinygo"
)

type testValues struct {
	dec []byte
	enc string
}

var n = 5000000
var testPairs = make([]testValues, 0, n)

func init() {
	// If we do not seed the prng - it will default to a seed of (1)
	// https://golang.org/pkg/math/rand/#Seed
	rand.Seed(time.Now().UTC().UnixNano())
}

func initTestPairs() {
	if len(testPairs) > 0 {
		return
	}
	// pre-make the test pairs, so it doesn't take up benchmark time...
	for i := 0; i < n; i++ {
		data := make([]byte, 32)
		rand.Read(data)
		testPairs = append(testPairs, testValues{dec: data, enc: base58.FastBase58Encoding(data)})
	}
}

func randAlphabet() *base58.Alphabet {
	// Permutes [0, 127] and returns the first 58 elements.
	var randomness [128]byte
	rand.Read(randomness[:])

	var bts [128]byte
	for i, r := range randomness {
		j := int(r) % (i + 1)
		bts[i] = bts[j]
		bts[j] = byte(i)
	}
	return base58.NewAlphabet(string(bts[:58]))
}

var btcDigits = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func TestInvalidAlphabetTooShort(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on alphabet being too short did not occur")
		}
	}()

	_ = base58.NewAlphabet(btcDigits[1:])
}

func TestInvalidAlphabetTooLong(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on alphabet being too long did not occur")
		}
	}()

	_ = base58.NewAlphabet("0" + btcDigits)
}

func TestInvalidAlphabetNon127(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on alphabet containing non-ascii chars did not occur")
		}
	}()

	_ = base58.NewAlphabet("\xFF" + btcDigits[1:])
}

func TestInvalidAlphabetDup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic on alphabet containing duplicate chars did not occur")
		}
	}()

	_ = base58.NewAlphabet("z" + btcDigits[1:])
}

func TestFastEqTrivialEncodingAndDecoding(t *testing.T) {
	for k := 0; k < 10; k++ {
		testEncDecLoop(t, randAlphabet())
	}
	testEncDecLoop(t, base58.DefaultAlphabet)
}

func testEncDecLoop(t *testing.T, alph *base58.Alphabet) {
	for j := 1; j < 256; j++ {
		var b = make([]byte, j)
		for i := 0; i < 100; i++ {
			rand.Read(b)
			fe := base58.FastBase58EncodingAlphabet(b, alph)
			te := base58.TrivialBase58EncodingAlphabet(b, alph)

			if fe != te {
				t.Errorf("encoding err: %#v", hex.EncodeToString(b))
			}

			fd, ferr := base58.FastBase58DecodingAlphabet(fe, alph)
			if ferr != nil {
				t.Errorf("fast error: %v", ferr)
			}
			td, terr := base58.TrivialBase58DecodingAlphabet(te, alph)
			if terr != nil {
				t.Errorf("trivial error: %v", terr)
			}

			if hex.EncodeToString(b) != hex.EncodeToString(td) {
				t.Errorf("decoding err: %s != %s", hex.EncodeToString(b), hex.EncodeToString(td))
			}
			if hex.EncodeToString(b) != hex.EncodeToString(fd) {
				t.Errorf("decoding err: %s != %s", hex.EncodeToString(b), hex.EncodeToString(fd))
			}
		}
	}
}

func BenchmarkTrivialBase58Encoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.TrivialBase58Encoding([]byte(testPairs[i].dec))
	}
}

func BenchmarkFastBase58Encoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.FastBase58Encoding(testPairs[i].dec)
	}
}

func BenchmarkTrivialBase58Decoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.TrivialBase58Decoding(testPairs[i].enc)
	}
}

func BenchmarkFastBase58Decoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.FastBase58Decoding(testPairs[i].enc)
	}
}
