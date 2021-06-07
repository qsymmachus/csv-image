package main

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Attempts to parse a CSV file containing base-64 encoded image data.
// Assumes the CSV has two fields, a unique identifier and a base-64 string:
//
//     <identifier>,<base-64 image string>
//
// This will attempt to parse the base-64 string and encode it as a JPEG image
// and write it in the './output' directory, using the unique identifier as the
// file name.
//
// If an error is encountered attempting to parse the data, it will dump the
// base-64 string to a '.txt' file instead to help with debugging.
//
// Usage:
//
//     csv-jpeg -filepath path/to/csv-file.csv
//
func main() {
	filepath := flag.String("filepath", "./test.csv", "Path to CSV to import")
	flag.Parse()

	reader, err := parseCSV(*filepath)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		filename, data := record[0], record[1]
		err = base64toJpg(data, filename)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// Creates a CSV reader from a CSV file at a specified filepath.
func parseCSV(filepath string) (*csv.Reader, error) {
	fmt.Printf("Importing file '%s'...\n", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(strings.NewReader(string(bytes)))
	return reader, nil
}

// Parses base-64 encoded `data` and writes it to './output/<filename.jpeg>'.
func base64toJpg(data, filename string) error {
	jpegFileName := fmt.Sprintf("./output/%s.jpeg", filename)
	fmt.Printf("Attempting to write to '%s'...\n", jpegFileName)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	m, formatString, err := image.Decode(reader)
	fmt.Printf("Format: %s\n", formatString)
	if err != nil {
		dumpData(data, filename)
		return err
	}

	f, err := os.OpenFile(jpegFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		dumpData(data, filename)
		return err
	}
	defer f.Close()

	err = jpeg.Encode(f, m, &jpeg.Options{Quality: 100})
	if err != nil {
		dumpData(data, filename)
		return err
	}
	fmt.Printf("Created '%s'\n", jpegFileName)

	return nil
}

// Writes `data` to './output/<filename>.txt'.
func dumpData(data, filename string) {
	dumpFileName := fmt.Sprintf("./output/%s.txt", filename)
	fmt.Printf("Dumping data to '%s' for debugging...\n", dumpFileName)

	f, err := os.OpenFile(dumpFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.WriteString(data + "\n")
}
