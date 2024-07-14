package exit

import (
	"fmt"
	"os"
)

func ExitWithError(err error) {
	fmt.Fprintln(os.Stderr, "Грешка:", err)
	os.Exit(1)
}
