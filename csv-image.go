package main

import (
	"encoding/base64"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
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
//     csv-image -csv path/to/csv-file.csv
//
func main() {
	filepath := flag.String("csv", "./test.csv", "Path to CSV to import")
	outputDir := flag.String("output", "./output", "Directory to write images to")
	flag.Parse()

	reader, err := parseCSV(*filepath)
	if err != nil {
		log.Fatalln(err)
	}

	var wg sync.WaitGroup
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		id, data := record[0], record[1]
		wg.Add(1)
		go base64ToImage(data, id, *outputDir, &wg)
	}
	wg.Wait()

	fmt.Printf("\nDone! Check %s for image output.\n", *outputDir)
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

// Attempts to parse a base-64 `data` string and encode it into an image, and writes
// the image to a file. Currently handles JPEG and PNG encoding.
func base64ToImage(data, id, outputDir string, wg *sync.WaitGroup) {
	var output string
	defer wg.Done()
	output = output + fmt.Sprintf("Attempting to decode data with ID: %s...\n", id)

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
	image, formatString, err := image.Decode(reader)
	output = output + fmt.Sprintf("Format: %s\n", formatString)
	if err != nil {
		fmt.Printf("Parsing error: %s\n", err)
		dumpData(data, id, outputDir)
		return
	}

	switch formatString {
	case "jpeg":
		output = output + encodeToJPEG(image, data, id, outputDir)
	case "png":
		output = output + encodeToPNG(image, data, id, outputDir)
	default:
		output = output + fmt.Sprintf("Unrecognized image format: %s\n", formatString)
		dumpData(data, id, outputDir)
	}

	fmt.Printf(output)
}

// Encodes image data into a PNG and writes it to `./output/<filename>.png`
func encodeToPNG(image image.Image, data, filename, outputDir string) (output string) {
	pngFilename := fmt.Sprintf("%s/%s.png", outputDir, filename)
	output = output + fmt.Sprintf("Writing to '%s'...\n", pngFilename)

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed create output directory '%s': %s\n", outputDir, err)
		return output
	}

	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed to write file '%s': %s\n", pngFilename, err)
		return output
	}
	defer f.Close()

	err = png.Encode(f, image)
	if err != nil {
		output = output + fmt.Sprintf("Parsing error: %s\n", err)
		dumpData(data, filename, outputDir)
		return output
	}

	output = output + fmt.Sprintf("Created '%s'\n\n", pngFilename)
	return output
}

// Encodes image datainto a JPEG and writes it to './output/<filename>.jpeg'.
func encodeToJPEG(image image.Image, data, filename, outputDir string) (output string) {
	jpegFileName := fmt.Sprintf("%s/%s.jpeg", outputDir, filename)
	output = output + fmt.Sprintf("Writing to '%s'...\n", jpegFileName)

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed create output directory '%s': %s\n", outputDir, err)
		return output
	}

	f, err := os.OpenFile(jpegFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed to write file '%s': %s\n", jpegFileName, err)
		return output
	}
	defer f.Close()

	err = jpeg.Encode(f, image, &jpeg.Options{Quality: 100})
	if err != nil {
		output = output + fmt.Sprintf("Parsing error: %s\n", err)
		dumpData(data, filename, outputDir)
		return output
	}

	output = output + fmt.Sprintf("Created '%s'\n\n", jpegFileName)
	return output
}

// Writes `data` to './output/<filename>.txt'.
func dumpData(data, filename, outputDir string) (output string) {
	dumpFileName := fmt.Sprintf("%s/%s.txt", outputDir, filename)
	output = output + fmt.Sprintf("Dumping data to '%s' for debugging...\n\n", dumpFileName)

	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed create output directory '%s': %s\n", outputDir, err)
		return output
	}

	f, err := os.OpenFile(dumpFileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		output = output + fmt.Sprintf("Failed to write to dump file: %s", err)
		return output
	}
	defer f.Close()

	_, err = f.WriteString(data + "\n")
	if err != nil {
		output = output + fmt.Sprintf("Failed to write to dump file: %s", err)
		return output
	}

	return output
}
