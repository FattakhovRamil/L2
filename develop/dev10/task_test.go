package main

import (
	"fmt"
	"testing"
	"time"
)

func TestHandlerArgs(t *testing.T) {
	tests := []struct {
		args       []string
		expected   options
		shouldFail bool
	}{
		{[]string{"go-telnet", "localhost", "8080"}, options{timeout: 10 * time.Second, host: "localhost", port: "8080"}, false},
		{[]string{"go-telnet", "--timeout=5s", "localhost", "8080"}, options{timeout: 5 * time.Second, host: "localhost", port: "8080"}, false},
		{[]string{"go-telnet", "--timeout=invalid", "localhost", "8080"}, options{}, true},
		{[]string{"go-telnet", "localhost"}, options{}, true},
		{[]string{"go-telnet", "localhost", "8080", "extra"}, options{}, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.args), func(t *testing.T) {
			opts, err := handlerArgs(tt.args)
			if (err != nil) != tt.shouldFail {
				t.Errorf("handlerArgs() error = %v, shouldFail = %v", err, tt.shouldFail)
				return
			}
			if opts != tt.expected {
				t.Errorf("handlerArgs() = %v, expected %v", opts, tt.expected)
			}
		})
	}
}

