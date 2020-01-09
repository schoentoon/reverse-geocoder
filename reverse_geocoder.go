package reverse_geocoder

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	kdtree "github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type DB struct {
	tree *kdtree.KDTree
}

type City struct {
	Name   string
	Admin1 string
	Admin2 string
	CC     string
}

func parsePoint(in []string) (*points.Point, error) {
	lat, err := strconv.ParseFloat(in[0], 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(in[1], 64)
	if err != nil {
		return nil, err
	}
	city := &City{
		Name:   in[2],
		Admin1: in[3],
		Admin2: in[4],
		CC:     in[5],
	}

	return points.NewPoint([]float64{lat, lon}, city), nil
}

func CreateDBFromCSVFile(file string) (*DB, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return CreateDBFromCSV(csv.NewReader(f))
}

func CreateDBFromCSV(reader *csv.Reader) (*DB, error) {
	out := &DB{
		tree: kdtree.New(nil),
	}
	header := true
	for {
		l, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if header {
			header = false
			if len(l) != 6 || !(l[0] == "lat" && l[1] == "lon" && l[2] == "name" && l[3] == "admin1" && l[4] == "admin2" && l[5] == "cc") {
				return nil, errors.New("incorrect csv format")
			}
			continue
		}

		point, err := parsePoint(l)
		if err != nil {
			return nil, err
		}
		out.tree.Insert(point)
	}

	return out, nil
}

func (d *DB) Search(lat, lon float64) *City {
	out := d.tree.KNN(points.NewPoint([]float64{lat, lon}, nil), 1)
	res := out[0].(*points.Point)

	return res.Data.(*City)
}

func ExampleF_Search() {
	db, err := CreateDBFromCSVFile("rg_cities.csv")
	if err != nil {
		panic(err)
	}

	// These coordinates point to Amsterdam Dam Square
	ams := db.Search(52.3729306, 4.8917547)

	fmt.Printf("%s", ams.Name)
	// Output: Amsterdam
}
