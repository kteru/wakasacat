package main

import "errors"

var (
	ErrorNoSync = errors.New("sync byte (0x47) is not found")
)

type TsPacket [188]byte

func (tp *TsPacket) IsSync() bool {
	return tp[0] == 0x47
}

func (tp *TsPacket) Gap() (int, error) {
	for i := 0; i < len(tp); i++ {
		if tp[i] == 0x47 {
			return i, nil
		}
	}

	return -1, ErrorNoSync
}

/*
	8	sync byte
	1	transport error indicator
	1	payload unit start indicator
	1	transport priority
	13	PID
	2	transport scrambling control
	2	adaptation field control
	4	continuity counter
*/

func (tp *TsPacket) PayloadUnitStartIndicator() bool {
	return tp[1]&0x40 == 0x40
}

func (tp *TsPacket) Pid() uint16 {
	return uint16(tp[1]&0x1F)<<8 + uint16(tp[2])
}
