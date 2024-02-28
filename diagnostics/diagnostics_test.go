package diagnostics

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	d := New()

	// Verify that traceID is not empty
	if d.traceID == "" {
		t.Errorf("Expected traceID to be non-empty, got empty")
	}

	// Verify that ctx is not nil
	if d.ctx == nil {
		t.Errorf("Expected ctx to be non-nil, got nil")
	}

	// Verify that stack is empty
	if len(d.stack) != 0 {
		t.Errorf("Expected stack to be empty, got %d items", len(d.stack))
	}
}

func TestFromContext(t *testing.T) {
	t.Run("With Context and existing Trace Id", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceID, "existing-trace-id")
		d := FromContext(ctx)

		assert.Equalf(t, "existing-trace-id", d.traceID, "Expected traceID to be 'existing-trace-id', got '%s'", d.traceID)
		assert.Equalf(t, ctx, d.ctx, "Expected ctx to be the same as the input context, got different")
		assert.Equalf(t, 0, len(d.stack), "Expected stack to be empty, got %d items", len(d.stack))
	})

	t.Run("With Context no TraceId", func(t *testing.T) {
		d := FromContext(context.Background())

		assert.NotEmptyf(t, d.traceID, "Expected traceID to be non-empty, got empty")
		contextTraceId := d.ctx.Value(TraceID).(string)
		assert.Equalf(t, d.traceID, contextTraceId, "Expected traceID to be the same as the context value, got different")
		assert.Equalf(t, 0, len(d.stack), "Expected stack to be empty, got %d items", len(d.stack))
	})
}

func TestAddInfo(t *testing.T) {
	d := New()

	// Add an info to an empty stack
	d.AddInfo("info1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after adding an info")
	assert.Equal(t, "", d.stack[0].Code, "Expected code to be empty")
	assert.Equal(t, "info1", d.stack[0].Description, "Expected description to be 'info1'")
	assert.Equal(t, Info, d.stack[0].Level, "Expected level to be Info")

	// Add an info with the same description
	d.AddInfo("info1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to still have 1 item after adding an info with the same description")

	// Add an info with a different description
	d.AddInfo("info2")
	assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 items after adding an info with a different description")
	assert.Equal(t, "", d.stack[1].Code, "Expected code to be empty")
	assert.Equal(t, "info2", d.stack[1].Description, "Expected description to be 'info2'")
	assert.Equal(t, Info, d.stack[1].Level, "Expected level to be Info")
}

func TestAddWarning(t *testing.T) {
	d := New()

	// Add a warning to an empty stack
	d.AddWarning("warning1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after adding a warning")
	assert.Equal(t, "", d.stack[0].Code, "Expected code to be empty")
	assert.Equal(t, "warning1", d.stack[0].Description, "Expected description to be 'warning1'")
	assert.Equal(t, Warning, d.stack[0].Level, "Expected level to be Warning")

	// Add a warning with the same description
	d.AddWarning("warning1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to still have 1 item after adding a warning with the same description")

	// Add a warning with a different description
	d.AddWarning("warning2")
	assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 items after adding a warning with a different description")
	assert.Equal(t, "", d.stack[1].Code, "Expected code to be empty")
	assert.Equal(t, "warning2", d.stack[1].Description, "Expected description to be 'warning2'")
	assert.Equal(t, Warning, d.stack[1].Level, "Expected level to be Warning")
}

func TestAddError(t *testing.T) {
	d := New()

	// Add an error to an empty stack
	err := errors.New("error1")
	d.AddError(err)
	assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after adding an error")
	assert.Equal(t, "", d.stack[0].Code, "Expected code to be empty")
	assert.Equal(t, err.Error(), d.stack[0].Description, "Expected description to be 'error1'")
	assert.Equal(t, Error, d.stack[0].Level, "Expected level to be Error")

	// Add an error with the same description
	d.AddError(err)
	assert.Equal(t, 1, len(d.stack), "Expected stack to still have 1 item after adding an error with the same description")

	// Add an error with a different description
	err = errors.New("error2")
	d.AddError(err)
	assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 items after adding an error with a different description")
	assert.Equal(t, "", d.stack[1].Code, "Expected code to be empty")
	assert.Equal(t, err.Error(), d.stack[1].Description, "Expected description to be 'error2'")
	assert.Equal(t, Error, d.stack[1].Level, "Expected level to be Error")
}

func TestAddErrorWithCode(t *testing.T) {
	d := New()

	// Add an error with code "code1" and description "error1"
	err1 := errors.New("error1")
	d.AddErrorWithCode("code1", err1)
	assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after adding an error with code and description")
	assert.Equal(t, "code1", d.stack[0].Code, "Expected code to be 'code1'")
	assert.Equal(t, err1.Error(), d.stack[0].Description, "Expected description to be 'error1'")
	assert.Equal(t, Error, d.stack[0].Level, "Expected level to be Error")

	// Add an error with the same code and description
	d.AddErrorWithCode("code1", err1)
	assert.Equal(t, 1, len(d.stack), "Expected stack to still have 1 item after adding an error with the same code and description")

	// Add an error with a different code and description
	err2 := errors.New("error2")
	d.AddErrorWithCode("code2", err2)
	assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 items after adding an error with a different code and description")
	assert.Equal(t, "code2", d.stack[1].Code, "Expected code to be 'code2'")
	assert.Equal(t, err2.Error(), d.stack[1].Description, "Expected description to be 'error2'")
	assert.Equal(t, Error, d.stack[1].Level, "Expected level to be Error")
}

func TestAddTrace(t *testing.T) {
	d := New()

	// Add a trace to an empty stack
	d.AddTrace("trace1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after adding a trace")
	assert.Equal(t, "", d.stack[0].Code, "Expected code to be empty")
	assert.Equal(t, "trace1", d.stack[0].Description, "Expected description to be 'trace1'")
	assert.Equal(t, Trace, d.stack[0].Level, "Expected level to be Trace")

	// Add a trace with the same description
	d.AddTrace("trace1")
	assert.Equal(t, 1, len(d.stack), "Expected stack to still have 1 item after adding a trace with the same description")

	// Add a trace with a different description
	d.AddTrace("trace2")
	assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 items after adding a trace with a different description")
	assert.Equal(t, "", d.stack[1].Code, "Expected code to be empty")
	assert.Equal(t, "trace2", d.stack[1].Description, "Expected description to be 'trace2'")
	assert.Equal(t, Trace, d.stack[1].Level, "Expected level to be Trace")
}

func TestAppend(t *testing.T) {
	t.Run("Append with new Item", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
		}

		// Create a diagnostic item
		diag := &DiagnosticItem{
			Code:        "code2",
			Description: "description2",
			Level:       Info,
		}

		// Append the diagnostic item to the stack
		d.Append(diag)

		// Verify that the stack has 1 item
		assert.Equal(t, 2, len(d.stack), "Expected stack to have 2 item after appending")

		// Verify that the appended item is the same as the diagnostic item
		assert.Equal(t, diag, d.stack[1], "Expected the appended item to be the same as the diagnostic item")
	})

	t.Run("Append with existing Item", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
		}

		// Create a diagnostic item
		diag := &DiagnosticItem{
			Code:        "code1",
			Description: "description1",
			Level:       Info,
		}

		// Append the diagnostic item to the stack
		d.Append(diag)

		// Verify that the stack has 1 item
		assert.Equal(t, 1, len(d.stack), "Expected stack to have 1 item after appending")

		// Verify that the appended item is the same as the diagnostic item
		assert.Equal(t, diag, d.stack[0], "Expected the appended item to be the same as the diagnostic item")
	})
}

func TestGetDiagnostics(t *testing.T) {
	d := New()

	// Add some diagnostic items to the stack
	d.stack = []*DiagnosticItem{
		{
			Code:        "code1",
			Description: "description1",
			Level:       Info,
		},
		{
			Code:        "code2",
			Description: "description2",
			Level:       Warning,
		},
		{
			Code:        "code3",
			Description: "description3",
			Level:       Error,
		},
	}

	// Call GetDiagnostics and verify the returned stack
	stack := d.GetDiagnostics()
	assert.Equal(t, d.stack, stack, "Expected the returned stack to be the same as the internal stack")
}

func TestGetTraceID(t *testing.T) {
	d := New()
	d.traceID = "test-trace-id"

	traceID := d.GetTraceID()

	assert.Equal(t, "test-trace-id", traceID, "Expected traceID to be 'test-trace-id'")
}

func TestDiagnostics_Context(t *testing.T) {
	d := New()
	ctx := d.Context()

	assert.NotNil(t, ctx, "Expected ctx to be non-nil")
	assert.Equal(t, d.ctx, ctx, "Expected ctx to be the same as the internal context")
}

func TestHasErrors(t *testing.T) {
	t.Run("With No Errors", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
		}

		hasErrors := d.HasErrors()

		assert.False(t, hasErrors, "Expected HasErrors to return false")
	})

	t.Run("With Errors", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Error,
			},
		}

		hasErrors := d.HasErrors()

		assert.True(t, hasErrors, "Expected HasErrors to return true")
	})
}

func TestHasWarnings(t *testing.T) {
	t.Run("With No Warnings", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Error,
			},
		}

		hasWarnings := d.HasWarnings()

		assert.False(t, hasWarnings, "Expected HasWarnings to return false")
	})

	t.Run("With Warnings", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Error,
			},
		}

		hasWarnings := d.HasWarnings()

		assert.True(t, hasWarnings, "Expected HasWarnings to return true")
	})
}

func TestErrors(t *testing.T) {
	t.Run("With Errors with code", func(t *testing.T) {
		d := New()

		// Add some diagnostic items to the stack
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Error,
			},
		}

		// Call the Errors method and verify the returned errors
		errors := d.Errors()
		assert.Equal(t, 1, len(errors), "Expected the returned errors to have 1 item")
		assert.Equal(t, "error code3: description3", errors[0].Error(), "Expected the error message to be 'error code3: description3'")
	})

	t.Run("With Errors with no code", func(t *testing.T) {
		d := New()

		// Add some diagnostic items to the stack
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "",
				Description: "description3",
				Level:       Error,
			},
		}

		// Call the Errors method and verify the returned errors
		errors := d.Errors()
		assert.Equal(t, 1, len(errors), "Expected the returned errors to have 1 item")
		assert.Equal(t, "error: description3", errors[0].Error(), "Expected the error message to be 'error: description3'")
	})
}

func TestWarnings(t *testing.T) {
	t.Run("With Warning with code", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Warning,
			},
			{
				Code:        "code4",
				Description: "description4",
				Level:       Error,
			},
		}

		expected := []string{"warning code2: description2", "warning code3: description3"}
		warnings := d.Warnings()

		assert.Equal(t, expected, warnings, "Expected warnings to be %v, but got %v", expected, warnings)
	})

	t.Run("With Warning with no code", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "code1",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Warning,
			},
			{
				Code:        "code4",
				Description: "description4",
				Level:       Error,
			},
		}

		expected := []string{"warning: description2", "warning code3: description3"}
		warnings := d.Warnings()

		assert.Equal(t, expected, warnings, "Expected warnings to be %v, but got %v", expected, warnings)
	})
}

func TestInfo(t *testing.T) {
	d := New()
	d.stack = []*DiagnosticItem{
		{
			Code:        "code1",
			Description: "description1",
			Level:       Info,
		},
		{
			Code:        "code2",
			Description: "description2",
			Level:       Warning,
		},
		{
			Code:        "code3",
			Description: "description3",
			Level:       Error,
		},
	}

	expected := []string{"description1"}

	result := d.Info()

	assert.Equal(t, expected, result, "Expected the Info method to return the correct info descriptions")
}

func TestTrace(t *testing.T) {
	d := New()

	// Add some diagnostic items to the stack
	d.stack = []*DiagnosticItem{
		{
			Code:        "code1",
			Description: "description1",
			Level:       Info,
		},
		{
			Code:        "code2",
			Description: "description2",
			Level:       Warning,
		},
		{
			Code:        "code3",
			Description: "description3",
			Level:       Error,
		},
		{
			Code:        "code4",
			Description: "trace1",
			Level:       Trace,
		},
		{
			Code:        "code5",
			Description: "trace2",
			Level:       Trace,
		},
	}

	// Call the Trace method
	traces := d.Trace()

	// Verify the returned traces
	expectedTraces := []string{"trace: trace1", "trace: trace2"}
	assert.Equal(t, expectedTraces, traces, "Expected the returned traces to match the expected traces")
}

func TestStack(t *testing.T) {
	d := New()

	// Add some diagnostic items to the stack
	d.stack = []*DiagnosticItem{
		{
			Code:        "code1",
			Description: "description1",
			Level:       Info,
		},
		{
			Code:        "code2",
			Description: "description2",
			Level:       Warning,
		},
		{
			Code:        "code3",
			Description: "description3",
			Level:       Error,
		},
	}

	// Call Stack and verify the returned stack
	stack := d.Stack()
	assert.Equal(t, d.stack, stack, "Expected the returned stack to be the same as the internal stack")
}

func TestString(t *testing.T) {
	t.Run("With TraceID", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), TraceID, "test-trace-id")
		d := FromContext(ctx)
		d.traceID = "test-trace-id"
		d.stack = []*DiagnosticItem{
			{
				Code:        "",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Error,
			},
		}

		expected := "[test-trace-id][Info] description1\n[test-trace-id][Warning] code2: description2\n[test-trace-id][Error] code3: description3\n"
		result := d.String()

		assert.Equal(t, expected, result, "Expected the string representation to match the expected value")
	})

	t.Run("With auto generated TraceID", func(t *testing.T) {
		d := New()
		d.stack = []*DiagnosticItem{
			{
				Code:        "",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Error,
			},
		}

		expected := fmt.Sprintf("[%v][Info] description1\n[%v][Warning] code2: description2\n[%v][Error] code3: description3\n", d.traceID, d.traceID, d.traceID)
		result := d.String()

		assert.Equal(t, expected, result, "Expected the string representation to match the expected value")
	})

	t.Run("Without TraceID", func(t *testing.T) {
		d := New()
		d.traceID = ""
		d.stack = []*DiagnosticItem{
			{
				Code:        "",
				Description: "description1",
				Level:       Info,
			},
			{
				Code:        "code2",
				Description: "description2",
				Level:       Warning,
			},
			{
				Code:        "code3",
				Description: "description3",
				Level:       Error,
			},
		}

		expected := "[Info] description1\n[Warning] code2: description2\n[Error] code3: description3\n"
		result := d.String()

		assert.Equal(t, expected, result, "Expected the string representation to match the expected value")
	})
}
