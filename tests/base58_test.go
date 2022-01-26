package tests

import (
	"math/rand"
	"testing"
	"time"

	base58 "github.com/VSmert/base58-for-tinygo"
	"github.com/stretchr/testify/require"
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
		testPairs = append(testPairs, testValues{dec: data, enc: base58.Encode(data)})
	}
}

var defaultDigits = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func TestInvalidAlphabetTooShort(t *testing.T) {
	actualAlphabet := base58.NewAlphabet(defaultDigits[1:])
	require.Nil(t, actualAlphabet)
}

func TestInvalidAlphabetTooLong(t *testing.T) {
	actualAlphabet := base58.NewAlphabet("0" + defaultDigits)
	require.Nil(t, actualAlphabet)
}

func TestInvalidAlphabetNon127(t *testing.T) {
	actualAlphabet := base58.NewAlphabet("\xFF" + defaultDigits[1:])
	require.Nil(t, actualAlphabet)
}

func TestInvalidAlphabetDup(t *testing.T) {
	actualAlphabet := base58.NewAlphabet("z" + defaultDigits[1:])
	require.Nil(t, actualAlphabet)
}

func BenchmarkFastBase58Encoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.Encode(testPairs[i].dec)
	}
}

func BenchmarkFastBase58Decoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		base58.Decode(testPairs[i].enc)
	}
}
