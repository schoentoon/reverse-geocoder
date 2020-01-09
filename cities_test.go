package reverse_geocoder

import (
	"testing"
)

const rg_cities_file = "rg_cities.csv"

func expect(t *testing.T, provided, expected string) {
	if provided != expected {
		t.Errorf("Expected %s, got %s", expected, provided)
	}
}

func TestSomeCities(t *testing.T) {
	db, err := CreateDBFromCSVFile(rg_cities_file)
	if err != nil {
		t.Fatal(err)
	}

	// These coordinates point to Amsterdam Dam Square
	ams := db.Search(52.3729306, 4.8917547)

	expect(t, ams.Name, "Amsterdam")
	expect(t, ams.CC, "NL")
}

func BenchmarkSearch(b *testing.B) {
	db, err := CreateDBFromCSVFile(rg_cities_file)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.Search(52.3729306, 4.8917547)
	}
}
