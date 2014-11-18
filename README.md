wakasacat
=========

わかさトラップ回避コマンド

Description
-----------

MPEG-2 TS パケットの PAT/PMT を解析し、  
途中で PMT の示すストリームの PID が変化する (わかさトラップ) 場合、その PMT の直前の PAT 以降のデータを標準出力に出力する。  
ストリームの PID が変化しない場合は、そのまま標準出力に出力する。

Usage
-----

ファイルから読み込む

```
$ wakasacat hoge.ts | avconv -i pipe:0 ...
```

標準入力から読み込む

```
$ cat hoge.ts | wakasacat | avconv -i pipe:0 ...
```

Command Line Options
--------------------

### `-m <num>`

デフォルト: `1000000`

先頭から `<num>` 個の TS パケット内の PAT/PMT を解析する。

Install
-------

```
$ go get github.com/kteru/wakasacat
```

Todo
----

- 詳細出力オプション
- わかさトラップの位置を出力して終了するオプション
- プログラム ID の指定オプション
- PID の変更検知をもっと厳密に (映像と音声の PID の切り替わりを対象にする)

Licence
-------

[MIT License](LICENSE)

Author
------

[teru](https://github.com/kteru)
