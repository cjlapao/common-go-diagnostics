package diagnostics

import "fmt"

type DiagnosticItem struct {
	Code        string
	Description string
	Level       DiagnosticLevel
}

type DiagnosticLevel int

const (
	Info DiagnosticLevel = iota
	Warning
	Error
	Trace
)

func (d DiagnosticLevel) String() string {
	return [...]string{"Info", "Warning", "Error", "Trace"}[d]
}

func (d DiagnosticLevel) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

func (d *DiagnosticLevel) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case "\"Info\"":
		*d = Info
	case "\"Warning\"":
		*d = Warning
	case "\"Error\"":
		*d = Error
	case "\"Trace\"":
		*d = Trace
	}
	return nil
}

func NewDiagnosticItem(errorCode string, errorDescription string, errorLevel DiagnosticLevel) *DiagnosticItem {
	return &DiagnosticItem{
		Code:        errorCode,
		Description: errorDescription,
		Level:       errorLevel,
	}
}

func NewError(errorCode string, errorDescription string) *DiagnosticItem {
	return NewDiagnosticItem(errorCode, errorDescription, Error)
}

func NewWarning(errorCode string, errorDescription string) *DiagnosticItem {
	return NewDiagnosticItem(errorCode, errorDescription, Warning)
}

func NewInfo(errorCode string, errorDescription string) *DiagnosticItem {
	return NewDiagnosticItem(errorCode, errorDescription, Info)
}

func NewTrace(errorCode string, errorDescription string) *DiagnosticItem {
	return NewDiagnosticItem(errorCode, errorDescription, Trace)
}

func (d *DiagnosticItem) String() string {
	msg := fmt.Sprintf("[%v] %v", d.Level, d.Description)
	if d.Code != "" {
		msg = fmt.Sprintf("[%v] %v: %v", d.Level, d.Code, d.Description)
	}

	return msg
}
