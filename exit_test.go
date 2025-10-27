package stdcli

import (
	"testing"
)

func TestExit(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		wantCode int
		wantMsg  string
	}{
		{
			name:     "exit with code 0",
			code:     0,
			wantCode: 0,
			wantMsg:  "exit 0",
		},
		{
			name:     "exit with code 1",
			code:     1,
			wantCode: 1,
			wantMsg:  "exit 1",
		},
		{
			name:     "exit with code 42",
			code:     42,
			wantCode: 42,
			wantMsg:  "exit 42",
		},
		{
			name:     "exit with negative code",
			code:     -1,
			wantCode: -1,
			wantMsg:  "exit -1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := Exit(tt.code)

			if ec.ExitCode() != tt.wantCode {
				t.Errorf("ExitCode() = %d, want %d", ec.ExitCode(), tt.wantCode)
			}

			// Cast to error to access Error() method
			if err, ok := ec.(error); ok {
				if msg := err.Error(); msg != tt.wantMsg {
					t.Errorf("Error() = %q, want %q", msg, tt.wantMsg)
				}
			} else {
				t.Errorf("Exit() should implement error interface")
			}
		})
	}
}

func TestExitCoderInterface(t *testing.T) {
	ec := Exit(0)
	var _ ExitCoder = ec

	// Verify it also implements error
	if _, ok := ec.(error); !ok {
		t.Errorf("Exit() should implement error interface")
	}
}
