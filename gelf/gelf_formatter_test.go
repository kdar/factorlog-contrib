package gelf

import (
	"strings"
	"testing"
)

type gelfTest struct {
	SomeVal string `json:"some_val"`
}

func TestGELFFormatter(t *testing.T) {
	// l := log.New(os.Stdout, NewGELFFormatter())
	// l.Print("hey there yeaaa!!!", GELF{Data: &gelfTest{SomeVal: "hey"}, Type: "TEST"})
	// l = log.New(os.Stdout, log.NewStdFormatter("[%{Date} %{Time}] [%{SEVERITY}:%{File}:%{Line}] %{Message}"))
	// l.Print("hey there yeaaa!!!", GELF{Data: &gelfTest{SomeVal: "hey"}, Type: "TEST"})
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func BenchmarkCopyMapDelete(b *testing.B) {
	original := map[string]interface{}{
		"Accept":         []string{"*/*"},
		"Content-Length": []string{"529"},
		"Content-Type":   []string{"application/x-www-form-urlencoded"},
		"User-Agent":     []string{"curl/7.22.0 (x86_64-pc-linux-gnu) libcurl/7.22.0 OpenSSL/1.0.1 zlib/1.2.3.4 libidn/1.23 librtmp/2.3"},
		"authorization":  "caca",
	}

	for x := 0; x < b.N; x++ {
		newmap := make(map[string]interface{})
		for k, v := range original {
			newmap[k] = v
		}
		delete(newmap, "Authorization")
		delete(newmap, "authorization")
	}
}

func BenchmarkCopyMapTestDelete(b *testing.B) {
	original := map[string]interface{}{
		"Accept":         []string{"*/*"},
		"Content-Length": []string{"529"},
		"Content-Type":   []string{"application/x-www-form-urlencoded"},
		"User-Agent":     []string{"curl/7.22.0 (x86_64-pc-linux-gnu) libcurl/7.22.0 OpenSSL/1.0.1 zlib/1.2.3.4 libidn/1.23 librtmp/2.3"},
		"authorization":  "caca",
	}

	for x := 0; x < b.N; x++ {
		newmap := make(map[string]interface{})
		for k, v := range original {
			if strings.EqualFold(k, "authorization") {
				continue
			}
			newmap[k] = v
		}
	}
}
