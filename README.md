# teapot-wm

制作 : 多田 瑛貴 <br>
ライセンス : MIT <br>

teapot-wm は、C言語とGo言語で書かれた簡易ウインドウマネージャです。<br>

### 各ソースコードファイルの説明

`wm_main.go`  : `main()`関数の定義<br>

`wm_x.go`  : Xlibの簡易ラッパー<br>
`c_wm_x_access.h(.c)`   : `wm_x.go`にて、Go言語だけで実装しきれない機能を担う<br>

`wm_system.go`  : 基本的な機能の定義<br>
`wm_host.go`  : Hostの定義と関数<br>
`wm_event.go`  : XEventに関する機能<br>
`wm_event_loop.go`  : イベントループの定義<br>

`wm_debug.go`  : logによるデバッグ機能に関する機能<br>

### Makefileの説明

`make` : ビルド<br>
`make init_xinitrc` : `~/.xinitrc`の設定。**勝手に書き換えてしまうので気をつけて！**<br>
`make init_log_file` : `log`で出力したログの内容を保存するファイルを新規作成or初期化 (`./wmlog.txt`)<br>
`make show_log_file` : ログファイルの内容の表示<br>
