package qrcode

import "testing"

func TestEncoder_analysisMode(t *testing.T) {
	enc := new(encoder)
	enc.str="1234567890"
	enc.analysisMode()
	if enc.mode != numericMode {
		t.FailNow()
	}

	enc.str="0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"
	enc.analysisMode()
	if enc.mode != alphanumericMode {
		t.FailNow()
	}

	enc.str="蘭のテーマ"
	enc.analysisMode()
	if enc.mode != byteMode {
		t.FailNow()
	}

	enc.str="你好abc123"
	enc.analysisMode()
	if enc.mode != byteMode {
		t.FailNow()
	}
}
