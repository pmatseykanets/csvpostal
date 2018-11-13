package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"

	parser "github.com/openvenues/gopostal/parser"
)

func buildIndex(components []string) (index map[string]int) {
	index = make(map[string]int)
	for i, c := range components {
		index[c] = i
	}
	return
}

var (
	inFile  string
	outFile string
)

func main() {
	flag.StringVar(&inFile, "in", "", "Name of the input file")
	flag.StringVar(&outFile, "out", "", "Name of the output file")
	flag.Parse()

	var reader io.Reader
	if inFile == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(inFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	var writer io.Writer
	if outFile == "" {
		writer = os.Stdout
	} else {
		file, err := os.Create(outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		writer = bufio.NewWriter(file)
	}

	if err := do(reader, writer); err != nil {
		log.Fatal(err)
	}
}

func do(reader io.Reader, writer io.Writer) error {
	components := []string{"house", "category", "near", "house_number", "road", "unit", "level", "staircase", "entrance", "po_box", "postcode", "suburb", "city_district", "city", "island", "state_district", "state", "country_region", "country", "world_region"}
	index := buildIndex(components)

	csvReader := csv.NewReader(reader)

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Read the header
	header, err := csvReader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	// Find the address column
	addressIndex := -1
	for i, name := range header {
		if strings.EqualFold(name, "address") {
			addressIndex = i
			break
		}
	}

	if addressIndex < 0 {
		return errors.New("No address field found")
	}

	// Write header
	if err := csvWriter.Write(append(header, components...)); err != nil {
		return err
	}

	for {
		inRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		addr := strings.TrimSpace(inRecord[addressIndex])

		outRecord := make([]string, len(components))

		if addr == "" {
			if err := csvWriter.Write(append(inRecord, outRecord...)); err != nil {
				return err
			}
			continue
		}

		parsed := parser.ParseAddress(addr)

		for _, c := range parsed {
			i, ok := index[c.Label]
			if !ok {
				log.Fatalf("component %s not found", c.Label)
			}
			outRecord[i] = c.Value
		}
		if err := csvWriter.Write(append(inRecord, outRecord...)); err != nil {
			return err
		}
	}

	return nil
}
