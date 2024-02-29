package diagnostics

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Diagnostics struct {
	traceID string
	ctx     context.Context
	stack   []*DiagnosticItem
}

func New() *Diagnostics {
	traceId := uuid.New().String()
	result := Diagnostics{
		traceID: traceId,
		ctx:     context.WithValue(context.Background(), TraceID, traceId),
		stack:   []*DiagnosticItem{},
	}

	return &result
}

func FromContext(ctx context.Context) *Diagnostics {
	traceId := uuid.New().String()
	if ctx != nil {
		if ctx.Value(TraceID) != nil {
			if id, ok := ctx.Value(TraceID).(string); ok {
				traceId = id
			}
		} else {
			ctx = context.WithValue(ctx, TraceID, traceId)
		}
	}

	return &Diagnostics{
		traceID: traceId,
		ctx:     ctx,
		stack:   []*DiagnosticItem{},
	}
}

func (d *Diagnostics) AddInfo(info string) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, "") && strings.EqualFold(i.Description, info) && i.Level == Info {
			return
		}
	}

	d.stack = append(d.stack, NewInfo("", info))
}

func (d *Diagnostics) AddWarning(warning string) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, "") && strings.EqualFold(i.Description, warning) && i.Level == Warning {
			return
		}
	}

	d.stack = append(d.stack, NewWarning("", warning))
}

func (d *Diagnostics) AddError(err error) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, "") && strings.EqualFold(i.Description, err.Error()) && i.Level == Error {
			return
		}
	}

	d.stack = append(d.stack, NewError("", err.Error()))
}

func (d *Diagnostics) AddErrorWithCode(code string, err error) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, code) && strings.EqualFold(i.Description, err.Error()) && i.Level == Error {
			return
		}
	}

	d.stack = append(d.stack, NewError(code, err.Error()))
}

func (d *Diagnostics) AddTrace(trace string) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, "") && strings.EqualFold(i.Description, trace) && i.Level == Trace {
			return
		}
	}

	d.stack = append(d.stack, NewTrace("", trace))
}

func (d *Diagnostics) AddItem(diag *DiagnosticItem) {
	for _, i := range d.stack {
		if strings.EqualFold(i.Code, diag.Code) && strings.EqualFold(i.Description, diag.Description) && i.Level == diag.Level {
			return
		}
	}

	d.stack = append(d.stack, diag)
}

func (d *Diagnostics) Append(diagnostic *Diagnostics) {
	for _, i := range d.stack {
		for _, j := range diagnostic.stack {
			if strings.EqualFold(i.Code, j.Code) && strings.EqualFold(i.Description, j.Description) && i.Level == j.Level {
				return
			}

			d.stack = append(d.stack, j)
		}
	}
}

func (d *Diagnostics) GetDiagnostics() []*DiagnosticItem {
	return d.stack
}

func (d *Diagnostics) GetTraceID() string {
	return d.traceID
}

func (d *Diagnostics) Context() context.Context {
	return d.ctx
}

func (d *Diagnostics) HasErrors() bool {
	for _, i := range d.stack {
		if i.Level == Error {
			return true
		}
	}

	return false
}

func (d *Diagnostics) HasWarnings() bool {
	for _, i := range d.stack {
		if i.Level == Warning {
			return true
		}
	}

	return false
}

func (d *Diagnostics) Errors() []error {
	result := []error{}
	for _, i := range d.stack {
		if i.Level == Error {
			if i.Code == "" {
				result = append(result, fmt.Errorf("error: %v", i.Description))
			} else {
				result = append(result, fmt.Errorf("error %v: %v", i.Code, i.Description))
			}
		}
	}

	return result
}

func (d *Diagnostics) Warnings() []string {
	result := []string{}
	for _, i := range d.stack {
		if i.Level == Warning {
			if i.Code == "" {
				result = append(result, fmt.Sprintf("warning: %v", i.Description))
			} else {
				result = append(result, fmt.Sprintf("warning %v: %v", i.Code, i.Description))
			}
		}
	}

	return result
}

func (d *Diagnostics) Info() []string {
	result := []string{}
	for _, i := range d.stack {
		if i.Level == Info {
			result = append(result, fmt.Sprintf("%v", i.Description))
		}
	}

	return result
}

func (d *Diagnostics) Trace() []string {
	result := []string{}
	for _, i := range d.stack {
		if i.Level == Trace {
			result = append(result, fmt.Sprintf("trace: %v", i.Description))
		}
	}

	return result
}

func (d *Diagnostics) Stack() []*DiagnosticItem {
	return d.stack
}

func (d *Diagnostics) String() string {
	result := ""
	for _, i := range d.stack {
		if d.traceID != "" {
			result = fmt.Sprintf("%v[%v]%v\n", result, d.traceID, i.String())
		} else {
			result = fmt.Sprintf("%v%v\n", result, i.String())
		}
	}

	return result
}
