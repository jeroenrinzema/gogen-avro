package serializer

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestReadingLong(t *testing.T) {
	// expected result - input
	inputs := map[int64][]byte{
		20:                  []byte{40},
		30:                  []byte{60},
		9223372036854775807: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 1}, // 64 bit signed int max
		60:                  []byte{120},
		15:                  []byte{30},
	}

	r, w := io.Pipe()

	for expected, input := range inputs {
		go w.Write(input)

		result, err := ReadLong(r)
		if err != nil {
			t.Fatal(err)
		}

		if result != expected {
			t.Fatalf("bytes: %b, are interperated incorrectly expected result %d recieved %d", input, expected, result)
		}
	}
}

func TestWritingLong(t *testing.T) {
	// expected result - input
	inputs := map[int64][]byte{
		20:                  []byte{40},
		30:                  []byte{60},
		9223372036854775807: []byte{254, 255, 255, 255, 255, 255, 255, 255, 255, 1}, // 64 bit signed int max
		60:                  []byte{120},
		15:                  []byte{30},
	}

	for input, expected := range inputs {
		r, w := io.Pipe()

		go func() {
			err := WriteLong(w, input)
			if err != nil {
				t.Fatal(err)
			}

			w.Close()
		}()

		bb, _ := ioutil.ReadAll(r)

		if len(bb) != len(expected) {
			t.Fatalf("the returned byte buffer has an unexpected length: %b, %b\n", bb, expected)
		}

		for i, b := range bb {
			if b != expected[i] {
				t.Fatalf("unexpected byte encountered: %b, %b\n", b, expected[i])
			}
		}
	}
}

func BenchmarkReadingLong(b *testing.B) {
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		WriteLong(bb, 2147483648)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ReadLong(bb)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWritingLong(b *testing.B) {
	inputs := make([]int64, b.N)
	bb := bytes.NewBuffer(nil)

	for i := 0; i < b.N; i++ {
		inputs = append(inputs, 9223372036854775807)
	}

	b.ResetTimer()

	for _, input := range inputs {
		WriteLong(bb, input)
	}
}
