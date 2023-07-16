package escpos

import (
	"encoding/binary"
	"fmt"
)

type ESCcommand struct {
	Cmd        []byte
	Err        []error
	chainCount int
}

const (
	esc = 0x1B
	fs  = 0x1C
	gs  = 0x1D
)

func New() *ESCcommand {
	return &ESCcommand{
		Cmd:        []byte{},
		Err:        []error{},
		chainCount: 0,
	}
}

func (e *ESCcommand) Text(s string) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte(s)...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=52
func (e *ESCcommand) Tab() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, 0x09)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=10
func (e *ESCcommand) NewLine() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, 0x0A)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=88
type BitImageMode byte

const (
	Height8NormalWidth  BitImageMode = 0x00
	Height8HalfWidth    BitImageMode = 0x01
	Height24NormalWidth BitImageMode = 0x20
	Height24HalfWidth   BitImageMode = 0x21
)

func (e *ESCcommand) PrintBitImage(imageMode BitImageMode, image []byte) *ESCcommand {
	e.chainCount++

	if ((imageMode == Height24NormalWidth) || (imageMode == Height24HalfWidth)) && len(image)/3 != 0 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, ImageModeHeight24 is image []bytes requestments 3n length", e.chainCount))
	}
	if len(image) > 0xffff {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, image []bytes requestments less than 0xffff", e.chainCount))
	}
	// Get Wigth Size
	length := make([]byte, 8)
	if imageMode == Height24NormalWidth || imageMode == Height24HalfWidth {
		binary.LittleEndian.PutUint64(length, uint64(len(image)/3))
	} else {
		binary.LittleEndian.PutUint64(length, uint64(len(image)))
	}
	e.Cmd = append(e.Cmd, []byte{esc, '*', byte(imageMode), length[0], length[1]}...)
	e.Cmd = append(e.Cmd, image...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=15
func (e *ESCcommand) PrintAndFeed(feed byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'J', feed}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=16
func (e *ESCcommand) PrintAndBackFeed(feed byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'K', feed}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=17
func (e *ESCcommand) PrintAndLineFeed(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'd', n}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=18
func (e *ESCcommand) PrintAndBackLineFeed(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'e', n}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=130
const (
	ErrorCorrectingLevelL byte = 48
	ErrorCorrectingLevelM byte = 49
	ErrorCorrectingLevelQ byte = 50
	ErrorCorrectingLevelH byte = 51
)

func (e *ESCcommand) PrintQRcode(dotSize, level byte, url string) *ESCcommand {
	e.chainCount++

	urlByte := []byte(url)
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(urlByte)+3))

	// QRcode Model, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=140
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x04, 0x00, 0x31, 0x31, 50, 0}...)
	// QRcode ModuleSize, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=141
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x03, 0x00, 0x31, 0x43, dotSize}...)
	// QRcode ErrorCorrectingLevel, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=142
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x03, 0x00, 0x31, 0x45, level}...)
	// QRcode URL, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=143
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', length[0], length[1], 49, 80, 48}...)
	e.Cmd = append(e.Cmd, urlByte...)
	// QRcode Print, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=144
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x03, 0x00, 0x31, 0x51, 48}...)
	return e
}
