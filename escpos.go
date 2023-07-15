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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=22
func (e *ESCcommand) SetRightSpace(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, ' ', n}...)
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

	e.Cmd = append(e.Cmd, []byte{esc, '!', c}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=54
func (e *ESCcommand) SetAbsolutePostion(leftL, leftH byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '$', leftL, leftH}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=19
func (e *ESCcommand) ResetNewLineHeight() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '2'}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=20
func (e *ESCcommand) SetNewLineHeight(height byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '3', height}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=66
func (e *ESCcommand) ReturnHome() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '<'}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=192
func (e *ESCcommand) ResetSetting() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '@'}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=53
func (e *ESCcommand) SetTabSize(size []byte) *ESCcommand {
	e.chainCount++

	if len(size) > 32 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, TabSize len is max 32", e.chainCount))
	}
	e.Cmd = append(e.Cmd, []byte{esc, 'D'}...)
	e.Cmd = append(e.Cmd, size...)
	e.Cmd = append(e.Cmd, 0x00)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=27
type FontType byte

const (
	FontA    FontType = 0x00
	FontB    FontType = 0x01
	FontC    FontType = 0x02
	FontD    FontType = 0x03
	FontE    FontType = 0x04
	SpecialA FontType = 0x61
	SpecialB FontType = 0x62
)

func (e *ESCcommand) SetFont(font FontType) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'M', byte(font)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=29
type LangType byte

const (
	America       LangType = 0x00
	France        LangType = 0x01
	Germany       LangType = 0x02
	UnitedKingdom LangType = 0x03
	Denmark       LangType = 0x04
	Sweden        LangType = 0x05
	Italy         LangType = 0x06
	Spain         LangType = 0x07
	Japan         LangType = 0x08
	Norway        LangType = 0x09
	Korea         LangType = 0x0D
	Slovenia      LangType = 0x0E
	China         LangType = 0x0F
	Vietnam       LangType = 0x10
	Arabia        LangType = 0x11
)

func (e *ESCcommand) SetFontLang(lang LangType) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'R', byte(lang)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=67
func (e *ESCcommand) SetTurnOneDirection(oneDirection bool) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'U'}...)
	if oneDirection {
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

func (e *ESCcommand) Set90RightRotate(rotate RightRotate) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'V', byte(rotate)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=58
type Align byte

const (
	Left   Align = 0x00
	Center Align = 0x01
	Right  Align = 0x02
)

func (e *ESCcommand) SetTextAlign(align Align) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 'a', byte(align)}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=32
func (e *ESCcommand) SetTextCodeTable(code byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 't', code}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=130
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
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x03, 0x00, 0x31, 0x45, 48}...)
	// QRcode URL Header, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=143
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', length[0], length[1], 49, 80, 48}...)
	// QRcode URL
	e.Cmd = append(e.Cmd, urlByte...)
	// QRcode Print, https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=144
	e.Cmd = append(e.Cmd, []byte{gs, '(', 'k', 0x03, 0x00, 0x31, 0x51, 48}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=125
type BarcodeTextPos byte

const (
	BarCodeTextNone BarcodeTextPos = 0x00
	Top             BarcodeTextPos = 0x01
	Bottom          BarcodeTextPos = 0x02
	TopAndBottom    BarcodeTextPos = 0x03
)

func (e *ESCcommand) BarcodeText(pos BarcodeTextPos) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'H', byte(pos)}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=199
func (e *ESCcommand) SetDefaultPitch(x, y byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'P', x, y}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=126
func (e *ESCcommand) BarcodeTextFont(font FontType) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'f', byte(font)}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=127
func (e *ESCcommand) BarcodeHeight(dot byte) *ESCcommand {
	e.chainCount++

	if dot < 1 {
		e.Err = append(e.Err, fmt.Errorf("chain:%d is error, BarcodeHeight is >= 1dot", e.chainCount))
	}
	e.Cmd = append(e.Cmd, []byte{gs, 'h', dot}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=128
func (e *ESCcommand) PrintBarcode(code string) *ESCcommand {
	e.chainCount++
	e.Cmd = append(e.Cmd, []byte{gs, 'k', 0x04}...)
	e.Cmd = append(e.Cmd, []byte(code)...)
	e.Cmd = append(e.Cmd, 0x00)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=129
func (e *ESCcommand) BarcodeWidth(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'w', n}...)
	return e
}
