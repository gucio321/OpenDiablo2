package d2dcc

import (
	"fmt"
	"os"

	"errors"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

const dccFileSignature = 0x74
const directionOffsetMultiplier = 8

// DCC represents a DCC file.
type DCC struct {
	Signature          int
	Version            int
	NumberOfDirections int
	FramesPerDirection int
	Directions         []*DCCDirection
	directionOffsets   []int
	fileData           []byte
	totalSize          int32
}

// Load loads a DCC file.
func Load(fileData []byte) (*DCC, error) {
	result := &DCC{
		fileData: fileData,
	}

	var bm = d2datautils.CreateBitMuncher(fileData, 0)

	result.Signature = int(bm.GetByte())

	if result.Signature != dccFileSignature {
		return nil, errors.New("signature expected to be 0x74 but it is not")
	}

	result.Version = int(bm.GetByte())
	result.NumberOfDirections = int(bm.GetByte())
	result.FramesPerDirection = int(bm.GetInt32())

	result.Directions = make([]*DCCDirection, result.NumberOfDirections)

	if bm.GetInt32() != 1 {
		return nil, errors.New("this value isn't 1. It has to be 1")
	}

	// generally it should be calculated
	result.totalSize = bm.GetInt32() // TotalSizeCoded

	result.directionOffsets = make([]int, result.NumberOfDirections)

	for i := 0; i < result.NumberOfDirections; i++ {
		result.directionOffsets[i] = int(bm.GetInt32())
		result.Directions[i] = result.decodeDirection(i)
	}

	ok := true
	d2, _ := result.Marshal(fileData)
	for i := range d2 {
		if d2[i] != fileData[i] {
			fmt.Printf("\nindex %d, org %d my %d", i, fileData[i], d2[i])
			ok = false
		}
	}
	fmt.Printf("\n%v\n", ok)
	os.Exit(0)

	return result, nil
}

// decodeDirection decodes and returns the given direction
func (d *DCC) decodeDirection(direction int) *DCCDirection {
	return CreateDCCDirection(d2datautils.CreateBitMuncher(d.fileData,
		d.directionOffsets[direction]*directionOffsetMultiplier), d)
}

// Clone creates a copy of the DCC
func (d *DCC) Clone() *DCC {
	clone := *d
	copy(clone.directionOffsets, d.directionOffsets)
	copy(clone.fileData, d.fileData)
	clone.Directions = make([]*DCCDirection, len(d.Directions))

	for i := range d.Directions {
		cloneDirection := *d.Directions[i]
		clone.Directions = append(clone.Directions, &cloneDirection)
	}

	return &clone
}

var mx = 12

func (d *DCC) Marshal(data []byte) ([]byte, error) {
	sw := d2datautils.CreateStreamWriter()

	if d.Signature != dccFileSignature {
		return nil, errors.New("error while encoding dcc file: invalid file signature")
	}

	sw.PushBytes(byte(d.Signature))
	sw.PushBytes(byte(d.Version))
	sw.PushBytes(byte(d.NumberOfDirections))
	sw.PushInt32(int32(d.FramesPerDirection))

	// nolint:gomnd // thats constant. in each file should be 1
	sw.PushInt32(1)

	sw.PushInt32(d.totalSize)

	for i := 0; i < d.NumberOfDirections; i++ {
		sw.PushInt32(int32(d.directionOffsets[i]))
	}

	i := 0
	//for i := 0; i < len(d.Directions); i++ {
	sw.PushUint32(uint32(d.Directions[i].OutSizeCoded))
	sw.PushBits(byte(d.Directions[i].CompressionFlags), 2)
	sw.PushBits(byte(d.Directions[i].Variable0Bits), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].WidthBits)), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].HeightBits)), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].XOffsetBits)), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].YOffsetBits)), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].OptionalDataBits)), 4)
	sw.PushBits(byte(getCBPos(d.Directions[i].CodedBytesBits)), 4)
	for j := 0; j < d.FramesPerDirection; j++ {
		sw.PushBits(byte(d.Directions[i].Frames[j].Width),
			d.Directions[i].WidthBits)
		sw.PushBits(byte(d.Directions[i].Frames[j].Height),
			d.Directions[i].HeightBits)
		sw.PushBits(byte(d.Directions[i].Frames[j].XOffset),
			d.Directions[i].XOffsetBits)
		sw.PushBits(byte(d.Directions[i].Frames[j].YOffset),
			d.Directions[i].YOffsetBits)
		sw.PushBits(byte(d.Directions[i].Frames[j].NumberOfOptionalBytes),
			d.Directions[i].OptionalDataBits)
		sw.PushBits16(uint16(d.Directions[0].Frames[j].NumberOfCodedBytes),
			d.Directions[i].CodedBytesBits)
		sw.PushBit(false)
	}

	if (d.Directions[0].CompressionFlags & 0x2) > 0 {
		sw.PushBits32(uint32(d.Directions[i].EqualCellsBitstreamSize), 20)
	}

	sw.PushBits32(uint32(d.Directions[i].PixelMaskBitstreamSize), 20)

	if (d.Directions[0].CompressionFlags & 0x1) > 0 {
		sw.PushBits32(uint32(d.Directions[i].EncodingTypeBitsreamSize), 20)
		sw.PushBits32(uint32(d.Directions[i].RawPixelCodesBitstreamSize), 20)
	}

	for j := 0; j < 256; j++ {
		// to push bool. (if 0, then false, else true)
		sw.PushBit(d.Directions[i].paletteValidate[j] != 0)
	}

	//}

	hw := sw.GetBytes()
	fmt.Printf("org: %v\nsw:  %v", data[len(hw)-mx:len(hw)+8], hw[len(hw)-mx:])

	return sw.GetBytes(), nil
}

func getCBPos(v int) int {
	var crazyBitTable = []byte{0, 1, 2, 4, 6, 8, 10, 12, 14, 16, 20, 24, 26, 28, 30, 32}
	for i := 0; i < len(crazyBitTable); i++ {
		if int(crazyBitTable[i]) == v {
			return i
		}
	}

	return 0
}
