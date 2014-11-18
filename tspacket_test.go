package main

import "testing"

func Test_IsSync(t *testing.T) {
	var pair = []struct {
		in  TsPacket
		out bool
	}{
		{testPkt01, true},
		{testPkt03, false},
		{testPkt04, false},
		{testPkt05, false},
	}

	for _, p := range pair {
		result := p.in.IsSync()
		if result != p.out {
			t.Errorf("IsSync() = %v | want %v", result, p.out)
		}
	}
}

func Test_Gap(t *testing.T) {
	var pair = []struct {
		in   TsPacket
		out1 int
		out2 error
	}{
		{testPkt01, 0, nil},
		{testPkt03, 1, nil},
		{testPkt04, 187, nil},
		{testPkt05, -1, ErrorNoSync},
	}

	for _, p := range pair {
		result1, result2 := p.in.Gap()
		if result1 != p.out1 || result2 != p.out2 {
			t.Errorf("Gap() = %v, %v | want %v, %v", result1, result2, p.out1, p.out2)
		}
	}
}

func Test_PayloadUnitStartIndicator(t *testing.T) {
	var pair = []struct {
		in  TsPacket
		out bool
	}{
		{testPkt01, true},
		{testPkt02, false},
	}

	for _, p := range pair {
		result := p.in.PayloadUnitStartIndicator()
		if result != p.out {
			t.Errorf("PayloadUnitStartIndicator() = %v | want %v", result, p.out)
		}
	}
}

func Test_Pid(t *testing.T) {
	var pair = []struct {
		in  TsPacket
		out uint16
	}{
		{testPkt01, 123},
		{testPkt02, 234},
	}

	for _, p := range pair {
		result := p.in.Pid()
		if result != p.out {
			t.Errorf("Pid() = %v | want %v", result, p.out)
		}
	}
}
