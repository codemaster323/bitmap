package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

type BMPheader struct {
	ID          [2]byte // BMP ID
	FileSize    uint32  // Changed to uint32 to match BMP spec
	Unused      [4]byte
	PixelOffset uint32
}

type DIBheader struct {
	HeaderSize     uint32
	Width          uint32
	Height         uint32
	ColorPlanes    uint16
	BitPerPixel    uint16
	RGB            uint32
	DataSize       uint32
	PWidth         uint32
	PHeight        uint32
	ColorsCount    uint32
	ImpColorsCount uint32
}

type BMPfile struct {
	BHdr BMPheader
	DHdr DIBheader
	Data []byte // Pixel data
}

func loadBMPfile(fname string) (*BMPfile, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var bmpReader BMPfile

	// Read BMP header
	if err := binary.Read(file, binary.LittleEndian, &bmpReader.BHdr); err != nil {
		return nil, err
	}

	// Read DIB header
	if err := binary.Read(file, binary.LittleEndian, &bmpReader.DHdr); err != nil {
		return nil, err
	}

	// Seek to the pixel data offset
	if _, err := file.Seek(int64(bmpReader.BHdr.PixelOffset), 1); err != nil {
		return nil, err
	}

	// Allocate space for pixel data
	bmpReader.Data = make([]byte, bmpReader.DHdr.DataSize)
	if _, err := file.Read(bmpReader.Data); err != nil {
		return nil, err
	}

	return &bmpReader, nil
}

func createMatrix(bmpFile BMPfile) [][]byte {
	// Calculate row size based on bits per pixel and width
	rowSize := ((int(bmpFile.DHdr.BitPerPixel)*int(bmpFile.DHdr.Width) + 31) / 32) * 4
	// Calculate the number of rows
	columnSize := bmpFile.DHdr.DataSize / uint32(rowSize)

	// Initialize the matrix with the correct dimensions
	matrix := make([][]byte, rowSize)
	for i := range matrix {
		matrix[i] = make([]byte, columnSize*3) // 3 bytes for R, G, B
	}

	// Fill the matrix with pixel data
	for i := 0; i < int(rowSize); i++ {
		for j := 0; j < int(columnSize); j++ {
			idx := (i*int(columnSize) + j) * 3
			// Ensure that idx does not exceed the bounds of bmpFile.Data
			if idx+2 < len(bmpFile.Data) {
				matrix[i][j*3] = bmpFile.Data[idx]     // R
				matrix[i][j*3+1] = bmpFile.Data[idx+1] // G
				matrix[i][j*3+2] = bmpFile.Data[idx+2] // B
			}
		}
	}
	return matrix
}

func outputNewFile(fname string, bmpFile BMPfile) {
	inputFile, err := os.Open(fname)
	if err != nil {
		return
	}
	defer inputFile.Close()
	// Read pixel data
	_, err = inputFile.Read(bmpFile.Data)
	if err != nil {
		fmt.Println("Error reading pixel data:", err)
		return
	}

	// Create a new BMP file
	outputFile, err := os.Create("new.bmp")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Write the header to the new BMP file
	err = binary.Write(outputFile, binary.LittleEndian, bmpFile.BHdr)
	if err != nil {
		fmt.Println("Error writing BMP header:", err)
		return
	}

	// Write the pixel data to the new BMP file
	_, err = outputFile.Write(bmpFile.Data)
	if err != nil {
		fmt.Println("Error writing pixel data:", err)
		return
	}
}

func printHeader(bmpFile BMPfile) {
	fmt.Println("BMP Header:")
	fmt.Println("- FileType:", string(bmpFile.BHdr.ID[:]))
	fmt.Println("- FileSizeInBytes:", bmpFile.BHdr.FileSize)
	fmt.Println("- HeaderSize:", bmpFile.BHdr.PixelOffset)
	fmt.Println("DIB Header:")
	fmt.Println("- DibHeaderSize:", bmpFile.DHdr.HeaderSize)
	fmt.Println("- WidthInPixels:", bmpFile.DHdr.Width)
	fmt.Println("- HeightInPixels:", bmpFile.DHdr.Height)
	fmt.Println("- PixelSizeInBits:", bmpFile.DHdr.BitPerPixel)
	fmt.Println("- ImageSizeInBytes:", bmpFile.DHdr.DataSize)
}

func main() {
	header := flag.Bool("header", false, "show header")
	flag.Parse()

	bmpFile, err := loadBMPfile("sample.bmp")
	if err != nil {
		fmt.Println("Error loading BMP file:", err)
		return
	}
	if *header {
		printHeader(*bmpFile)
		return
	}
	matrix := createMatrix(*bmpFile)
	fmt.Println(matrix[0][0])
	outputNewFile("sample.bmp", *bmpFile)
}
