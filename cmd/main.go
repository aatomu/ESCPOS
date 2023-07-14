package main

import "os"

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
	printer.WriteString("\x1B\x28\x41\x04\x00\x30\x00\x05\x0A")
	// img := []byte{}
	// var i byte
	// for i = 0; i < 185; i++ {
	// 	img = append(img, i)
	// }
	// b := make([]byte, 8)
	// binary.LittleEndian.PutUint64(b, uint64(len(img)))
	// fmt.Printf("\\x1b\\x2a\\x00%s\n", b[:2])
	// printer.WriteString(fmt.Sprintf("\x1b\x2a\x00%s%s", b[:2], img))
	// return
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	csvReader := csv.NewReader(strings.NewReader(scanner.Text()))
	// 	arr, err := csvReader.Read()
	// 	if err != nil {
	// 		log.Println(err)
	// 		continue
	// 	}

	// 	var write string
	// loop:
	// 	for _, v := range arr {
	// 		cmd := strings.Split(v, " ")
	// 		switch cmd[0] {
	// 		case "print":
	// 			write += strings.ReplaceAll(cmd[1], "\\n", "\n")

	// 		case "under": // ESC -     アンダーラインの指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1b\x2d\x00"
	// 			case "1":
	// 				write += "\x1b\x2d\x01"
	// 			case "2":
	// 				write += "\x1b\x2d\x02"
	// 			}
	// 		case "bold": // ESC E     強調印字の指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1b\x45\x00"
	// 			case "1":
	// 				write += "\x1b\x45\x01"
	// 			}
	// 		case "double": // ESC G     二重印字の指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1b\x47\x00"
	// 			case "1":
	// 				write += "\x1b\x47\x01"
	// 			}
	// 		case "roll": // ESC {     倒立印字の指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1b\x7b\x00"
	// 			case "1":
	// 				write += "\x1b\x7b\x01"
	// 			}
	// 		case "rotate": // ESC V     文字の90度右回転の指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1b\x56\x00"
	// 			case "1":
	// 				write += "\x1b\x56\x01"
	// 			case "2":
	// 				write += "\x1b\x56\x02"
	// 			}
	// 		case "reverse": // GS B	白黒反転印字の指定・解除
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x1d\x42\x00"
	// 			case "1":
	// 				write += "\x1d\x42\x01"
	// 			}
	// 		case "color": // GS ( N <機能48>	文字色の選択
	// 			write += "\x1d\x28\x4e\x02\x00\x30"
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x48"
	// 			case "1":
	// 				write += "\x49"
	// 			case "2":
	// 				write += "\x50"
	// 			case "3":
	// 				write += "\x51"
	// 			}
	// 		case "back": // GS ( N <機能49>	背景色の選択
	// 			write += "\x1d\x28\x4e\x02\x00\x31"
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x48"
	// 			case "1":
	// 				write += "\x49"
	// 			case "2":
	// 				write += "\x50"
	// 			case "3":
	// 				write += "\x51"
	// 			}
	// 		case "shadow": // GS ( N <機能50>	影付き文字装飾の指定・解除
	// 			write += "\x1d\x28\x4e\x02\x00\x32"
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x00"
	// 			case "1":
	// 				write += "\x01"
	// 			}
	// 			switch cmd[2] {
	// 			case "0":
	// 				write += "\x48"
	// 			case "1":
	// 				write += "\x49"
	// 			case "2":
	// 				write += "\x50"
	// 			case "3":
	// 				write += "\x51"
	// 			}
	// 		case "size": // GS !	文字サイズの指定
	// 			write += "\x1d\x21"
	// 			var s int
	// 			n, _ := strconv.Atoi(cmd[1])
	// 			s += n
	// 			s = s << 4
	// 			n, _ = strconv.Atoi(cmd[2])
	// 			s += n
	// 			write += fmt.Sprint(s)
	// 		case "align": // ESC a	位置揃え
	// 			write += "\x1b\x61"
	// 			switch cmd[1] {
	// 			case "0":
	// 				write += "\x00"
	// 			case "1":
	// 				write += "\x01"
	// 			case "2":
	// 				write += "\x02"
	// 			}

	// 		case "head": // GS T	行の先頭への印字位置の移動
	// 			write += "\x1d\x54\x01"
	// 		case "move": // ESC \	相対位置の指定
	// 			write += "\x1b\x5c"
	// 			sp, _ := strconv.Atoi(cmd[1])
	// 			write += fmt.Sprint(sp)
	// 			write += "\x00"
	// 		case "space": // ESC SP    文字の右スペース量の設定
	// 			write += "\x1b\x20"
	// 			sp, _ := strconv.Atoi(cmd[1])
	// 			write += fmt.Sprint(sp)
	// 		case "reset": // ESC @	プリンターの初期化
	// 			write += "\x1b\x40"

	// 		case "cut": // FS ( L <機能66>	カット位置までの紙送り
	// 			write += "\x1c\x28\x4c\x02\x00\x42\x49"

	// 		default:
	// 			write = ""
	// 			log.Println("Failed Command")
	// 			break loop
	// 		}
	// 	}

	// 	printer.Write([]byte(write))
	// }
}

// GS ( K <機能48>	印字制御モードの選択
// GS ( K <機能49>	印字濃度の選択
// GS ( K <機能50>	印字速度の選択
