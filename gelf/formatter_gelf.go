package gelf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/kardianos/osext"
	log "github.com/kdar/factorlog"
)

//  Extra is a type you wrap around data to be added
// as extra fields to JSON.
type Extra map[string]interface{}

// GLEFFormatter provides a formatter that outputs GELF v1.1 data.
// It will also encode extra information such as an http.Request.
type GELFFormatter struct {
	hostname string
	appname  string
	tmp      []byte
}

// NewGELFFormatter returns a new GELFFormatter
func NewGELFFormatter() *GELFFormatter {
	f := &GELFFormatter{tmp: make([]byte, 64)}

	var err error
	if f.hostname, err = os.Hostname(); err != nil {
		f.hostname = "unknown"
	}

	if f.appname, err = osext.Executable(); err == nil {
		f.appname = filepath.Base(f.appname)
	} else {
		f.appname = "app"
	}

	return f
}

// ShouldRuntimeCaller will always return true.
func (f *GELFFormatter) ShouldRuntimeCaller() bool {
	return true
}

// gelf is the strucuture of GELF 1.1 messages.
type gelf struct {
	Version  string                 `json:"version"`
	Host     string                 `json:"host"`
	Short    string                 `json:"short_message"`
	Full     string                 `json:"full_message,omitempty"`
	TimeUnix int64                  `json:"timestamp"`
	Level    int32                  `json:"level,ompitempty"`
	Extra    map[string]interface{} `json:"-"`
}

type gelfBase gelf

// MarshalJSON is a special marshaller for marshalling the gelf
// structure and the extra data.
func (g *gelf) MarshalJSON() ([]byte, error) {
	var err error
	var b, eb []byte

	extra := g.Extra
	b, err = json.Marshal((*gelfBase)(g))
	g.Extra = extra
	if err != nil {
		return nil, err
	}

	if len(extra) == 0 {
		return b, nil
	}

	if eb, err = json.Marshal(extra); err != nil {
		return nil, err
	}

	b[len(b)-1] = ','
	return append(b, eb[1:len(eb)]...), nil
}

// Format formats a LogContext into the GELF format.
func (f *GELFFormatter) Format(context log.LogContext) []byte {
	file := context.File
	if len(file) == 0 {
		file = "???"
	} else {
		slash := len(file) - 1
		for ; slash >= 0; slash-- {
			if file[slash] == filepath.Separator {
				break
			}
		}
		if slash >= 0 {
			file = file[slash+1:]
		}
	}

	e := Extra{
		"_file":     file,
		"_line":     context.Line,
		"_severity": log.LcSeverityStrings[log.SeverityToIndex(context.Severity)],
		"_facility": f.appname + "." + context.Function,
	}

	var args []interface{}
	for _, v := range context.Args {
		switch t := v.(type) {
		case Extra:
			for key, value := range t {
				e[key] = value
			}
		case *http.Request:
			header := make(map[string]interface{})
			// Copy header
			for key, value := range t.Header {
				// Don't copy auth. Don't want this in logs.
				// Note: If you add more auth schemes then you have to update
				// this code to remove those too.
				if strings.EqualFold(key, "authorization") {
					continue
				}

				header[key] = value
			}

			e["_request_method"] = t.Method
			// We have to use RequestURI because URL may be
			// modified by routes.
			e["_request_url"] = t.RequestURI //t.URL.String(),
			e["_request_host"] = t.Host
			e["_request_remote_addr"] = t.RemoteAddr
			e["_request_header"] = header

		default:
			args = append(args, v)
		}
	}

	message := ""
	if context.Format != nil {
		message = fmt.Sprintf(*context.Format, args...)
	} else {
		message = fmt.Sprint(args...)
	}

	gf := gelf{
		Version: "1.1",
		Host:    f.hostname,
		Short:   message,
		//Full:     message,
		TimeUnix: time.Now().Unix(),
		Level:    6, // info
		Extra:    e,
	}

	buf, _ := json.Marshal(&gf)
	buf = append(buf, '\n')

	return buf
}
