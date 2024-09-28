package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type BMPheader struct {
	ID          [2]byte // Changed to byte for BMP ID
	FileSize    uint32
	Unused      [4]uint32
	PixelOffset uint32
}

type DIBheader struct {
	HeaderSize     uint32
	Width          int32
	Height         int32
	ColorPlanes    uint16
	BitPerPixel    uint16
	RGB            uint32
	DataSize       uint32
	PWidth         int32
	PHeight        int32
	ColorsCount    uint32
	ImpColorsCount uint32
}

type BMPfile struct {
	BHdr BMPheader
	DHdr DIBheader
	Data []byte // Changed to byte slice for pixel data
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
	if _, err := file.Seek(int64(bmpReader.BHdr.PixelOffset), 0); err != nil {
		return nil, err
	}

	// Allocate space for pixel data
	bmpReader.Data = make([]byte, bmpReader.DHdr.DataSize)
	if _, err := file.Read(bmpReader.Data); err != nil {
		return nil, err
	}

	return &bmpReader, nil
}

func printHeader(bmpFile BMPfile) {
	fmt.Println("BMP Header:")
	fmt.Println("- FileType", string(bmpFile.BHdr.ID[:]))
	fmt.Println("- FileSizeInBytes", bmpFile.BHdr.FileSize)
	fmt.Println("- HeaderSize", string(bmpFile.BHdr.PixelOffset))
	fmt.Println("DIB Header:")
	fmt.Println("- DibHeaderSize", string(bmpFile.BHdr.ID[:]))
	fmt.Println("- WidthInPixels", bmpFile.BHdr.FileSize)
	fmt.Println("- HeightInPixels", bmpFile.BHdr.PixelOffset)
	fmt.Println("- PixelSizeInBits", string(bmpFile.BHdr.ID[:]))
	fmt.Println("- ImageSizeInBytes", bmpFile.BHdr.FileSize)
}

func main() {
	bmpFile, err := loadBMPfile("sample.bmp")
	if err != nil {
		fmt.Println("Error loading BMP file:", err)
		return
	}
	printHeader(*bmpFile)
	fmt.Println("DIB Header:", bmpFile.DHdr)
	// You can process bmpFile.Data as needed
}
