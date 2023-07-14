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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=22
func (e *ESCcommand) SetRightSpace(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x20, n}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=23
func (e *ESCcommand) SetPrintSetting(font, bold, heightDouble, widthDouble, underlined bool) *ESCcommand {
	e.chainCount++

	var c byte = 0
	if font { // 0bit
		c += 1 << 0
	}
	if bold { // 3bit
		c += 1 << 3
	}
	if heightDouble { // 4bit
		c += 1 << 4
	}
	if widthDouble { // 5bit
		c += 1 << 5
	}
	if underlined { // 7bit
		c += 1 << 7
	}

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x21, c}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=54
func (e *ESCcommand) SetAbsolutePostion(leftL, leftH byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x24, leftL, leftH}...)
	return e
}

type BitImageMode byte

const (
	ImageModeHeight8NormalWidth  BitImageMode = 0x00
	ImageModeHeight8HalfWidth    BitImageMode = 0x01
	ImageModeHeight24NormalWidth BitImageMode = 0x20
	ImageModeHeight24HalfWidth   BitImageMode = 0x21
)

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=88
func (e *ESCcommand) PrintBitImage(imageMode BitImageMode, image []byte) *ESCcommand {
	e.chainCount++

	if ((imageMode == ImageModeHeight24NormalWidth) || (imageMode == ImageModeHeight24HalfWidth)) && len(image)/3 != 0 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, ImageModeHeight24 is image []bytes requestments 3n length", e.chainCount))
	}
	if len(image) > 0xffff {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, image []bytes requestments less than 0xffff", e.chainCount))
	}
	// Get Wigth Size
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(image)))

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x2A, byte(imageMode)}...)
	e.Cmd = append(e.Cmd, length[:2]...)
	e.Cmd = append(e.Cmd, image...)
	return e
}

type UnderlineMode byte

const (
	UnderlineNone   UnderlineMode = 0x00
	UnderlineOneDot UnderlineMode = 0x01
	UnderlineTwoDot UnderlineMode = 0x02
)

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=24
func (e *ESCcommand) Underlined(mode UnderlineMode) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x2D, byte(mode)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=19
func (e *ESCcommand) ResetNewLineHeight() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x32}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=20
func (e *ESCcommand) SetNewLineHeight(height byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x33, height}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=66
func (e *ESCcommand) ReturnHome() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x3C}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=192
func (e *ESCcommand) ResetSetting() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{0x1B, 0x40}...)
	return e
}

// 次:ESC D	水平タブ位置の設定	印字位置
