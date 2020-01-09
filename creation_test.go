package reverse_geocoder

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkDBCreationFromCSV(b *testing.B) {
	// we first just read the entire csv into memory
	f, err := os.Open(rg_cities_file)
	if err != nil {
		b.Fatal(err)
	}
	data, err := ioutil.ReadAll(f)
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
