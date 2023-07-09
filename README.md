# ESCPOS
ターミナルからコマンドを送るよう

デフォルト: /dev/usb/lp1 に対して送信

Command: `${CMD-1} {OPS0-1} ${OPS1-1},${CMD-2} {OPS0-2} ${OPS1-2},...`
* under
  * 0: なし
  * 1: あり 1dot
  * 2: あり 2dot
* bold
  * 0: なし
  * 1: あり
* double
  * 0: なし
  * 1: あり
* roll
  * 0: なし
  * 1: あり
* rotate
  * 0: なし
  * 1: あり (なんだっけ
  * 2: あり (なんだっけ
* reverse
  * 0: なし
  * 1: あり
* color
  * 0: なし
  * 1: あり 第一色
  * 2: あり 第二色
  * 3: あり 第三色
* back
  * 0: なし
  * 1: あり 第一色
  * 2: あり 第二色
  * 3: あり 第三色
* shadow 
  * 0: なし
  * 1: あり
* size
  * 0-7 縦と横 {Ops}+1の倍率
* align
  * 0: ひだり
  * 1: 真ん中
  * 2: みぎ
* head
* move
  * 0-255 横に{Ops}px
* space
  * 0-255 横に{Ops}px
* reset
* cut
