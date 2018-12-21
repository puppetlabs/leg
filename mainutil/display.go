package mainutil

import (
	"fmt"
	"io"
	"os"

	errawr "github.com/puppetlabs/errawr-go"
)

func ExitWithCLIError(w io.Writer, code int, err errawr.Error) {
	fmt.Fprintln(w, err.FormattedDescription())

	os.Exit(code)
}
