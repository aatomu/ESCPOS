package main

import (
	"fmt"
	"os"

	escpos "github.com/aatomu/ESCPOS"
)

var (
	printer *os.File
)

func init() {
	var err error
	printer, err = os.OpenFile("/dev/usb/lp1", os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
}

func main() {
	cmd := escpos.New().ResetSetting().PrintBarcode("0123456789")
	printer.Write(cmd.Cmd)
	fmt.Println(cmd.Err)

}
