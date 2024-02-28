package diagnostics_test

import (
	"errors"
	"fmt"

	"github.com/cjlapao/common-go-diagnostics/diagnostics"
)

type Item struct {
	ID   string
	Name string
}

func main() {
	diag := diagnostics.New()
	diag.AddInfo("This is an info")
	diag.AddWarning("This is a warning")
	diag.AddError(errors.New("this is an error"))

	fmt.Println(diag)
}
