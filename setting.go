package escpos

import (
	"fmt"
)

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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=22
func (e *ESCcommand) SetRightSpace(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, ' ', n}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=32
func (e *ESCcommand) SetTextCodeTable(code byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, 't', code}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=54
func (e *ESCcommand) SetAbsolutePostion(leftL, leftH byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '$', leftL, leftH}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=66
func (e *ESCcommand) ReturnHome() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '<'}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=192
func (e *ESCcommand) ResetSetting() *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{esc, '@'}...)
	return e
}

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=199
func (e *ESCcommand) SetDefaultPitch(x, y byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'P', x, y}...)
	return e
}
