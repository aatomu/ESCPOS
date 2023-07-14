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

// ESC D	水平タブ位置の設定	印字位置
// ESC E	強調印字の指定・解除	印字文字
// ESC G	二重印字の指定・解除	印字文字
// ESC J	印字および紙送り	印字命令
// ESC K	印字および逆方向紙送り	印字命令
// ESC L	ページモードの選択	補助機能
// ESC M	文字フォントの選択	印字文字
// ESC R	国際文字の選択	印字文字
// ESC S	スタンダードモードの選択	補助機能
// ESC T	ページモードにおける文字の印字方向の選択	印字位置
// ESC U	単方向印字の指定・解除	メカコントロール
// ESC V	文字の90度右回転の指定・解除	印字文字
// ESC W	ページモードにおける印字領域の設定	印字位置
// ESC \	相対位置の指定	印字位置
// ESC a	位置揃え	印字位置
// ESC c 3	紙なし信号出力に有効な紙なし検出器の選択	用紙の検出器
// ESC c 4	印字停止に有効な紙なし検出器の選択	用紙の検出器
// ESC c 5	パネルスイッチの有効・無効	パネルスイッチ
// ESC d	印字およびn行の紙送り	印字命令
// ESC e	印字およびn行の逆方向紙送り	印字命令
// ESC i (∗)	パーシャルカット (1点を残す)	メカコントロール
// ESC m (∗)	パーシャルカット (3点を残す)	メカコントロール
// ESC p	指定パルスの発生	補助機能
// ESC r	印字色の選択	印字文字
// ESC t	文字コードテーブルの選択	印字文字
// ESC u (∗)	周辺機器ステータスの送信	ステータス
// ESC v (∗)	用紙検出器ステータスの送信	ステータス
// ESC {	倒立印字の指定・解除	印字文字
// FS !	漢字の印字モードの一括指定	漢字制御
// FS &	漢字モードの指定	漢字制御
// FS ( A	漢字の文字装飾の指定	漢字制御
// FS ( A <機能48>	漢字フォントの選択	漢字制御
// FS ( C	コード変換方式の選択	印字文字
// FS ( C <機能48>	文字のエンコード種類の選択	印字文字
// FS ( C <機能60>	フォント優先順位設定	印字文字
// FS ( E	レシートエンハンス制御のコマンド	補助機能
// FS ( E <機能60>	トップロゴ/ボトムロゴ印字の設定値の抹消	補助機能
// FS ( E <機能61>	トップロゴ/ボトムロゴ印字の設定値の送信	補助機能
// FS ( E <機能62>	トップロゴ印字の設定	補助機能
// FS ( E <機能63>	ボトムロゴ印字の設定	補助機能
// FS ( E <機能64>	トップロゴ/ボトムロゴ印字の拡張設定	補助機能
// FS ( E <機能65>	トップロゴ/ボトムロゴ印字の有効・無効	補助機能
// FS ( L	ラベル紙／ブラックマーク紙の制御	印字用紙
// FS ( L <機能33>	用紙レイアウトの設定	印字用紙
// FS ( L <機能34>	用紙レイアウト情報の送信	印字用紙
// FS ( L <機能48>	位置情報の送信	印字用紙
// FS ( L <機能65>	剥離位置までの紙送り	印字用紙
// FS ( L <機能66>	カット位置までの紙送り	印字用紙
// FS ( L <機能67>	頭出し位置までの紙送り	印字用紙
// FS ( L <機能80>	用紙レイアウトエラーの特別なマージンの設定	印字用紙
// FS ( e	拡張機能に関する自動ステータス送信の有効・無効	ステータス
// FS -	漢字アンダーラインの指定・解除	漢字制御
// FS .	漢字モードの解除	漢字制御
// FS 2	外字の定義	漢字制御
// FS ?	外字の抹消	漢字制御
// FS C	漢字コード体系の選択	漢字制御
// FS S	漢字のスペース量の設定	漢字制御
// FS W	漢字の4倍角文字の指定・解除	漢字制御
// FS g 1 (∗)	ユーザーNVメモリーへのデータ書き込み	カスタマイズ
// FS g 2 (∗)	ユーザーNVメモリーデータの読み出し	カスタマイズ
// FS p (∗)	NVビットイメージの印字	ビットイメージ
// FS q (∗)	NVビットイメージの定義	ビットイメージ
// GS !	文字サイズの指定	印字文字
// GS $	ページモードにおける文字縦方向絶対位置の指定	印字位置
// GS ( A	テスト印字の実行	補助機能
// GS ( C	ユーザーNVメモリーの編集	カスタマイズ
// GS ( C <機能0>	指定レコードの消去	カスタマイズ
// GS ( C <機能1>	指定レコードへのデータの格納	カスタマイズ
// GS ( C <機能2>	指定レコードの格納データの送信	カスタマイズ
// GS ( C <機能3>	使用容量の送信	カスタマイズ
// GS ( C <機能4>	残容量の送信	カスタマイズ
// GS ( C <機能5>	格納レコードのキーコード一覧の送信	カスタマイズ
// GS ( C <機能6>	ユーザーNVメモリー全領域の一括消去	カスタマイズ
// GS ( D	リアルタイムコマンドの有効・無効	補助機能
// GS ( E	ユーザー設定コマンド群	カスタマイズ
// GS ( E <機能1>	ユーザー設定モードへの移行	カスタマイズ
// GS ( E <機能2>	ユーザー設定モードの終了	カスタマイズ
// GS ( E <機能3>	メモリースイッチの値の設定	カスタマイズ
// GS ( E <機能4>	メモリースイッチの値の送信	カスタマイズ
// GS ( E <機能5>	カスタマイズバリューの設定	カスタマイズ
// GS ( E <機能6>	カスタマイズバリューの送信	カスタマイズ
// GS ( E <機能7>	ユーザー定義ページのデータのコピー	カスタマイズ
// GS ( E <機能8>	作業領域の文字コードページへのデータ(カラム形式)の定義	カスタマイズ
// GS ( E <機能9>	作業領域の文字コードページへのデータ(ラスター形式)の定義	カスタマイズ
// GS ( E <機能10>	作業領域の文字コードページのデータの消去	カスタマイズ
// GS ( E <機能11>	シリアルインターフェイスの通信条件の設定	カスタマイズ
// GS ( E <機能12>	シリアルインターフェイスの通信条件の送信	カスタマイズ
// GS ( E <機能13>	Bluetoothインターフェイスの通信条件の設定	カスタマイズ
// GS ( E <機能14>	Bluetoothインターフェイスの通信条件の送信	カスタマイズ
// GS ( E <機能15>	USBインターフェイスの通信条件の設定	カスタマイズ
// GS ( E <機能16>	USBインターフェイスの通信条件の送信	カスタマイズ
// GS ( E <機能48>	用紙レイアウトの消去	カスタマイズ
// GS ( E <機能49>	用紙レイアウトの設定	カスタマイズ
// GS ( E <機能50>	用紙のレイアウト情報の送信	カスタマイズ
// GS ( E <機能99>	内蔵ブザーパターンの設定	カスタマイズ
// GS ( E <機能100>	内蔵ブザーパターンの送信	カスタマイズ
// GS ( H	レスポンス／状態通知に関するコマンド群	補助機能
// GS ( H <機能48>	プロセスID レスポンスの指定	補助機能
// GS ( H <機能49>	オフラインレスポンス送信の指定・解除	補助機能
// GS ( K	印字制御方法の選択	補助機能
// GS ( K <機能48>	印字制御モードの選択	補助機能
// GS ( K <機能49>	印字濃度の選択	補助機能
// GS ( K <機能50>	印字速度の選択	補助機能
// GS ( K <機能97>	サーマルヘッド通電の分割数の選択	補助機能
// GS ( L / GS 8 L	グラフィックスデータの指定	ビットイメージ
// GS ( L <機能48>	NVグラフィックスのメモリー容量の送信	ビットイメージ
// GS ( L <機能49>	グラフィックスの基本ドット密度の設定	ビットイメージ
// GS ( L <機能50>	プリントバッファーに格納されているグラフィックスデータの印字	ビットイメージ
// GS ( L <機能51>	NVグラフィックスメモリーの残容量の送信	ビットイメージ
// GS ( L <機能52>	ダウンロード・グラフィックスメモリーの残容量の送信	ビットイメージ
// GS ( L <機能64>	定義されているNVグラフィックスのキーコード一覧の送信	ビットイメージ
// GS ( L <機能65>	NVグラフィックスの全データの一括消去	ビットイメージ
// GS ( L <機能66>	指定されたNVグラフィックスデータの消去	ビットイメージ
// GS ( L / GS 8 L <機能67>	NVグラフィックスデータ(ラスター形式)の定義	ビットイメージ
// GS ( L / GS 8 L <機能68>	NVグラフィックスデータ(カラム形式)の定義	ビットイメージ
// GS ( L <機能69>	指定されたNVグラフィックスの印字	ビットイメージ
// GS ( L <機能80>	定義されているダウンロード・グラフィックスのキーコード一覧の送信	ビットイメージ
// GS ( L <機能81>	ダウンロード・グラフィックスの全データの一括消去	ビットイメージ
// GS ( L <機能82>	指定されたダウンロード・グラフィックスデータの消去	ビットイメージ
// GS ( L / GS 8 L <機能83>	ダウンロード・グラフィックスデータ(ラスター形式)の定義	ビットイメージ
// GS ( L / GS 8 L <機能84>	ダウンロード・グラフィックスデータ(カラム形式)の定義	ビットイメージ
// GS ( L <機能85>	指定されたダウンロード・グラフィックスの印字	ビットイメージ
// GS ( L / GS 8 L <機能112>	グラフィックスデータ(ラスター形式)のプリントバッファーへの格納	ビットイメージ
// GS ( L / GS 8 L <機能113>	グラフィックスデータ(カラム形式)のプリントバッファーへの格納	ビットイメージ
// GS ( M	プリンターのカスタマイズ	カスタマイズ
// GS ( M <機能1>	作業領域の設定値の保存領域へのセーブ	カスタマイズ
// GS ( M <機能2>	指定された設定値の作業領域へのロード	カスタマイズ
// GS ( M <機能3>	初期化処理における作業領域の設定値の選択	カスタマイズ
// GS ( N	文字装飾の指定	印字文字
// GS ( N <機能48>	文字色の選択	印字文字
// GS ( N <機能49>	背景色の選択	印字文字
// GS ( N <機能50>	影付き文字装飾の指定・解除	印字文字
// GS ( P	ページモードの制御	補助機能
// GS ( P <機能48>	ページモード選択時における印字可能領域の設定	補助機能
// GS ( Q	図形の描画	補助機能
// GS ( Q <機能48>	ページモードにおける直線の描画	補助機能
// GS ( Q <機能49>	ページモードにおける矩形の描画	補助機能
// GS ( Q <機能50>	スタンダードモードにおける横方向直線の描画	補助機能
// GS ( Q <機能51>	スタンダードモードにおける縦方向直線の描画	補助機能
// GS ( V	用紙カットの指定	メカコントロール
// GS ( V <機能48>	用紙のカット	メカコントロール
// GS ( V <機能49>	紙送りと用紙のカット	メカコントロール
// GS ( V <機能51>	用紙のカット予約	メカコントロール
// GS ( k	シンボルの設定と印字	2次元シンボル
// GS ( k <機能065>	PDF417: 桁数の設定	2次元シンボル
// GS ( k <機能066>	PDF417: 段数の設定	2次元シンボル
// GS ( k <機能067>	PDF417: モジュール幅の設定	2次元シンボル
// GS ( k <機能068>	PDF417: 段の高さの設定	2次元シンボル
// GS ( k <機能069>	PDF417: エラー訂正レベルの設定	2次元シンボル
// GS ( k <機能070>	PDF417: オプションの選択	2次元シンボル
// GS ( k <機能080>	PDF417: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能081>	PDF417: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能082>	PDF417: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能165>	QR Code: モデルの選択	2次元シンボル
// GS ( k <機能167>	QR Code: モジュールのサイズの設定	2次元シンボル
// GS ( k <機能169>	QR Code: エラー訂正レベルの選択	2次元シンボル
// GS ( k <機能180>	QR Code: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能181>	QR Code: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能182>	QR Code: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能265>	MaxiCode: モードの選択	2次元シンボル
// GS ( k <機能280>	MaxiCode: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能281>	MaxiCode: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能282>	MaxiCode: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能367>	2次元GS1 DataBar: モジュール幅の設定	2次元シンボル
// GS ( k <機能371>	2次元GS1 DataBar: GS1 DataBar Expanded Stacked の最大幅の設定	2次元シンボル
// GS ( k <機能380>	2次元GS1 DataBar: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能381>	2次元GS1 DataBar: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能382>	2次元GS1 DataBar: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能467>	Composite Symbology: モジュール幅の設定	2次元シンボル
// GS ( k <機能471>	Composite Symbology: GS1 DataBar Expanded Stacked の最大幅の設定	2次元シンボル
// GS ( k <機能472>	Composite Symbology: HRI文字のフォントの選択	2次元シンボル
// GS ( k <機能480>	Composite Symbology: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能481>	Composite Symbology: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能482>	Composite Symbology: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能566>	Aztec Code: モードタイプ、データレイヤー数の設定	2次元シンボル
// GS ( k <機能567>	Aztec Code: モジュールサイズの設定	2次元シンボル
// GS ( k <機能569>	Aztec Code: エラー訂正レベルの選択	2次元シンボル
// GS ( k <機能580>	Aztec Code: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能581>	Aztec Code: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能582>	Aztec Code: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS ( k <機能666>	DataMatrix: シンボルタイプ、行数、列数の設定	2次元シンボル
// GS ( k <機能667>	DataMatrix: モジュールサイズの設定	2次元シンボル
// GS ( k <機能680>	DataMatrix: シンボル保存領域へのデータの格納	2次元シンボル
// GS ( k <機能681>	DataMatrix: シンボル保存領域のシンボルデータの印字	2次元シンボル
// GS ( k <機能682>	DataMatrix: シンボル保存領域のシンボルデータのサイズ情報の送信	2次元シンボル
// GS * (∗)	ダウンロード・ビットイメージの定義	ビットイメージ
// GS / (∗)	ダウンロード・ビットイメージの印字	ビットイメージ
// GS :	マクロ定義の開始・終了	マクロ機能
// GS B	白黒反転印字の指定・解除	印字文字
// GS C 0 (∗)	カウンターの印字モードの設定	カウンター印字
// GS C 1 (∗)	カウントモードの設定(A)	カウンター印字
// GS C 2 (∗)	カウンター値の設定	カウンター印字
// GS C ; (∗)	カウントモードの設定(B)	カウンター印字
// GS D	Windows BMP データによるグラフィックスデータの指定	ビットイメージ
// GS D <機能67>	Windows BMP のNVグラフィックスデータの定義	ビットイメージ
// GS D <機能83>	Windows BMP のダウンロード・グラフィックスデータの定義	ビットイメージ
// GS H	HRI文字の印字位置の選択	バーコード
// GS I	プリンターID の送信	補助機能
// GS L	左マージンの設定	印字位置
// GS P	基本計算ピッチの設定	補助機能
// GS Q 0 (∗)	縦サイズ可変ビットイメージの印字	ビットイメージ
// GS T	行の先頭への印字位置の移動	印字位置
// GS V	用紙のカット	メカコントロール
// GS W	印字領域幅の設定	印字位置
// GS \	ページモードにおける文字縦方向相対位置の指定	印字位置
// GS ^	マクロの実行	マクロ機能
// GS a	自動ステータス送信の有効・無効	ステータス
// GS b	スムージングの指定・解除	印字文字
// GS c (∗)	カウンターの印字	カウンター印字
// GS f	HRI文字のフォントの選択	バーコード
// GS g 0	メンテナンスカウンターの初期化	補助機能
// GS g 2	メンテナンスカウンターの送信	補助機能
// GS h	バーコードの高さの設定	バーコード
// GS j	インクに関する自動ステータス送信の有効・無効	ステータス
// GS k	バーコードの印字	バーコード
// GS r	ステータスの送信	ステータス
// GS v 0 (∗)	ラスタービットイメージの印字	ビットイメージ
// GS w	バーコードの横サイズの設定	バーコード
// GS z 0	オンライン復帰待ち時間の設定	補助機能
