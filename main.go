package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	option = struct {
		maxSearchPackets int
		verbose          bool
		fileName         string
	}{
		1000000,
		false,
		"",
	}
)

func parseCmdArgs() {
	// オプションの処理
	flag.IntVar(&option.maxSearchPackets, "m", option.maxSearchPackets, "Number of ts packets to search Wakasa-trap (default: 1000000)")
	flag.BoolVar(&option.verbose, "v", option.verbose, "Enable verbose output to stderr")
	flag.Parse()

	// 引数の処理
	if len(flag.Args()) > 0 {
		option.fileName = flag.Arg(0)
	}
}

func vfprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	// verbose フラグが立っていたら出力
	if option.verbose {
		return fmt.Fprintf(w, format, a...)
	}

	return 0, nil
}

func main() {
	var f *os.File
	var err error
	var pkt TsPacket
	var gap int
	var wakasaSeekOffset int64

	parseCmdArgs()

	// 標準入力 or 引数から開く
	if option.fileName != "" {
		vfprintf(os.Stderr, "Reading from \"%s\"\n", option.fileName)
		f, err = os.Open(option.fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		vfprintf(os.Stderr, "Reading from stdin\n")
		f = os.Stdin
	}

	// sync_byte の頭出し処理
	if c, _ := f.Read(pkt[:]); c == 0 {
		fmt.Fprintf(os.Stderr, "TS packet is not found.\n")
		os.Exit(1)
	}
	gap, err = pkt.Gap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	} else {
		vfprintf(os.Stderr, "Gap is %d byte\n", gap)
		// gap の分先頭からシーク
		f.Seek(int64(gap), os.SEEK_SET)
	}

	// わかさトラップを探索
	vfprintf(os.Stderr, "Searching Wakasa-trap\n")

	var firstPmtPid uint16
	var firstElemPid, firstElemPidBefore uint16
	var secondElemPid, secondElemPidBefore uint16
	var patPos int64

	for i := 0; i < option.maxSearchPackets; i++ {
		// 1 パケット読む
		if c, _ := f.Read(pkt[:]); c < 188 {
			vfprintf(os.Stderr, "Reached the end of the input\n")
			break
		}
		if i == option.maxSearchPackets-1 {
			vfprintf(os.Stderr, "Reached the max number of search packet\n")
		}

		// PAT
		if pkt.Pid() == 0 {
			// 1 番目の PMT の pid を取得
			firstPmtPid = uint16(pkt[19]&0x1F)<<8 + uint16(pkt[20])
			// 位置を記録
			patPos = int64(i)

			continue
		}

		// PMT
		if pkt.Pid() == firstPmtPid && pkt.PayloadUnitStartIndicator() {
			// 固定長部分
			index := uint16(0)

			// 番組情報長
			index += uint16(pkt[index+15]&0xF)<<8 + uint16(pkt[index+16])
			index += 17

			// 1 番目の番組の 1 つめのストリームの pid を取得
			firstElemPid = uint16(pkt[index+1]&0xF)<<8 + uint16(pkt[index+2])

			// 可変長部分
			index += uint16(pkt[index+3]&0xF)<<8 + uint16(pkt[index+4])
			index += 5

			// 1 番目の番組の 2 つめのストリームの pid を取得
			secondElemPid = uint16(pkt[index+1]&0xF)<<8 + uint16(pkt[index+2])

			// わかさトラップを見つけたら直前の PAT 位置を記録し探索終了
			if (firstElemPid != firstElemPidBefore && firstElemPidBefore != 0) ||
				(secondElemPid != secondElemPidBefore && secondElemPidBefore != 0) {
				vfprintf(os.Stderr, "DISCOVERED Wakasa-trap at packet %d, previous PAT is packet %d\n", i+1, patPos+1)
				wakasaSeekOffset = patPos * int64(len(pkt))
				break
			}

			firstElemPidBefore = firstElemPid
			secondElemPidBefore = secondElemPid
			continue
		}
	}

	// わかさトラップの直前の PAT まで先頭からシーク
	f.Seek(int64(gap)+wakasaSeekOffset, os.SEEK_SET)

	// 標準出力に書き込む
	vfprintf(os.Stderr, "Starting output\n")

	bufr := bufio.NewReaderSize(f, 16384)
	bufw := bufio.NewWriterSize(os.Stdout, 16384)
	for {
		b, err := bufr.ReadByte()
		if err == io.EOF {
			break
		}
		bufw.WriteByte(b)
	}
	bufw.Flush()

	vfprintf(os.Stderr, "Done\n")

	f.Close()
}
