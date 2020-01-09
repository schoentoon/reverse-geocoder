package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	lib "gitlab.com/schoentoon/reverse-geocoder"
)

func main() {
	cities_file := flag.String("cities", "rg_cities.csv", "The csv file to actually use")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("Expected at least 2 arguments\n")
		fmt.Printf("Please call as %s <latitude> <longitude>\n", os.Args[0])
		os.Exit(1)
	}

	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Printf("Latitude should be a float\n")
		os.Exit(1)
	}

	lon, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		fmt.Printf("Longitude should be a float\n")
		os.Exit(1)
	}

	db, err := lib.CreateDBFromCSVFile(*cities_file)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("It looks like we're missing rg_cities.csv, make sure it's in your current directory. Run extract_cities.py if needed.")
			os.Exit(1)
		}
		panic(err)
	}

	out := db.Search(lat, lon)

	fmt.Printf("%#v\n", out)
}
