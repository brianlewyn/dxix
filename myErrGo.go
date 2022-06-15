package dxix

import "fmt"

// Todo: [ myErrGo ] ==============================================
// * My Custom Error Message.

// My Struct Error
type errGO struct {
	Happd string
}

// Define the structure of the error message.
func myErrGO(h string) *errGO {
	return &errGO{Happd: h}
}

// Customize the error message for your debugging.
func (e *errGO) Error() string {
	return fmt.Sprintf("\nWhat happened? %s", e.Happd)
}
