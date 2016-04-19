package scan

import (
	"testing"
	"time"
)

func TestScanOptions(t *testing.T) {
	so := NewScanOptions(25, 20, 50)

	if so.Port != 25 {
		t.Error("Expected Port 25, got", so.Port)
	}

	if so.ConnectionTimeout != 20*time.Second {
		t.Error("Expected Connection Timeout 20000000000, got", so.Port)
	}

	if so.IOTimeout != 50*time.Second {
		t.Error("Expected Input-Output Timeout 50000000000, got", so.Port)
	}
}
