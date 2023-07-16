package escpos

import (
	"encoding/binary"
	"fmt"
)

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=23
func (e *ESCcommand) Decoration(font, bold, heightDouble, widthDouble, underlined bool) *ESCcommand {
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

	e.Cmd = append(e.Cmd, []byte{esc, '!', c}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=24
type UnderlineMode byte

const (
	UnderlineReset UnderlineMode = 0x00
	OneDot         UnderlineMode = 0x01
	TwoDot         UnderlineMode = 0x02
)

func (e *ESCcommand) Underlined(mode UnderlineMode) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '-', byte(mode)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=25
func (e *ESCcommand) Bold(bold bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'E'}...)
	if bold {
		e.Cmd = append(e.Cmd, 0x01)
	} else {
		e.Cmd = append(e.Cmd, 0x01)
	}
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=26
func (e *ESCcommand) Double(double bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'G'}...)
	if double {
		e.Cmd = append(e.Cmd, 0x01)
	} else {
		e.Cmd = append(e.Cmd, 0x01)
	}
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=33
func (e *ESCcommand) Handstand(handstand bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '{'}...)
	if handstand {
		e.Cmd = append(e.Cmd, 0x01)
	} else {
		e.Cmd = append(e.Cmd, 0x00)
	}
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=35
func (e *ESCcommand) Reverse(reverse bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'B'}...)
	if reverse {
		e.Cmd = append(e.Cmd, 0x01)
	} else {
		e.Cmd = append(e.Cmd, 0x00)
	}
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=30
type RightRotate byte

const (
	RotateReset          RightRotate = 0x00
	OneDotSpace          RightRotate = 0x01
	OnePointHalfDotSpace RightRotate = 0x02
)

func (e *ESCcommand) Turn90Rotate(rotate RightRotate) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'V', byte(rotate)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=34
func (e *ESCcommand) SetLetterSize(h, w int) *ESCcommand {
	e.chainCount++

	if h > 8 || h < 1 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, letter Height is x1-x8", e.chainCount))
	}
	h--
	if w > 8 || w < 1 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, letter Width is x1-x8", e.chainCount))
	}
	w--

	zoom := make([]byte, 8)
	binary.LittleEndian.PutUint64(zoom, uint64(w<<4+h))

	e.Cmd = append(e.Cmd, []byte{gs, '!', zoom[0]}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=36
func (e *ESCcommand) TextSmooth(smooth bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'b'}...)
	if smooth {
		e.Cmd = append(e.Cmd, 0x01)
	} else {
		e.Cmd = append(e.Cmd, 0x00)
	}
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=58
type Align byte

const (
	Left   Align = 0x00
	Center Align = 0x01
	Right  Align = 0x02
)

func (e *ESCcommand) TextAlign(align Align) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'a', byte(align)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=60
func (e *ESCcommand) SetLeftMargin(n int) *ESCcommand {
	e.chainCount++

	margin := make([]byte, 8)
	binary.LittleEndian.PutUint64(margin, uint64(n))

	e.Cmd = append(e.Cmd, []byte{gs, 'L', margin[0], margin[1]}...)
	return e
}
