package reverse_geocoder

import (
	"bytes"
	"encoding/csv"
	"io"
	"testing"
)

func BenchmarkDBCreationFromCSV(b *testing.B) {
	data, err := Asset("rg_cities.csv")
	if err != nil {
		b.Fatal(err)
	}
	buffer := bytes.NewReader(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = CreateDBFromCSV(csv.NewReader(buffer))
		if err != nil {
			b.Fatal(err)
		}

		// reset our memory buffer back to the start
		_, err = buffer.Seek(0, io.SeekStart)
		if err != nil {
			b.Fatal(err)
		}
	}
}
