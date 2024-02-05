package filesystem

import "fmt"

const (
	unhandledErrorCode = iota + 3000
	directoryNotFoundErrorCode
)

type UnhandledError string

func (p UnhandledError) Error() string {
	return fmt.Sprintf("error code: %d - unhandled exception: %q", unhandledErrorCode, string(p))
}

type DirectoryNotFoundError string

func (p DirectoryNotFoundError) Error() string {
	return fmt.Sprintf("error code: %d - Project Directory not Found %q", directoryNotFoundErrorCode, string(p))
}
