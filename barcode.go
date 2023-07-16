package escpos

import "fmt"

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=125
type BarcodeTextPos byte

const (
	TextNone     BarcodeTextPos = 0x00
	Top          BarcodeTextPos = 0x01
	Bottom       BarcodeTextPos = 0x02
	TopAndBottom BarcodeTextPos = 0x03
)

func (e *ESCcommand) BarcodeText(pos BarcodeTextPos) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'H', byte(pos)}...)
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

// https://reference.epson-biz.com/modules/ref_escpos_ja/index.php?content_id=129
func (e *ESCcommand) BarcodeWidth(n byte) *ESCcommand {
	e.chainCount++

	e.Cmd = append(e.Cmd, []byte{gs, 'w', n}...)
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
