# `csv-image`

```
Usage of ./csv-image:
  -csv string
    	Path to CSV to import (default "./test.csv")
  -output string
    	Directory to write images to (default "./output")
```

This program parses a CSV file containing base-64 encoded image data, and writes those images to files.

It assumes the CSV has two fields, a unique identifier and a base-64 string:

```
<identifier>,<base-64 data>
```

## Installation

Assuming you have [Go installed](https://golang.org/doc/install), install this program with `go get`:

```
go get github.com/qsymmachus/csv-image
```

## Example Usage

To parse the CSV file `my-image-data.csv` and output the images to a directory name `images`:

```
csv-image -csv my-image-data.csv -output images
```

It will attempt to

1. Parse each row of the `-csv` file, extracting the identifier and data.
1. Parse the base-64 data.
1. Encode the data as either a PNG or JPEG image.
1. Write the image to the specified `-output` directory, using the unique identifier as the file name, plus a file extension.

If an error is encountered attempting to parse the data, it will dump the base-64 string to a '.txt' file instead to help with debugging.
