// Package main illustrates an opinionated way on how to handle error
// propagation in a complex system
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

type MyError struct {
	Inner      error // error we are wrapping for investigation
	Message    string
	StackTrace string
	Misc       map[string]interface{} // storing miscellaneous info
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}),
	}
}

// implementing the Error interface
func (e MyError) Error() string {
	return e.Message
}

// low-level module
type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{wrapError(err, err.Error())}
	}
	return info.Mode().Perm()&0100 == 0100, nil
}

// intermediate module
type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const jobBinaryPath = "/bad/job/path"
	isExecutable, err := isGloballyExec(jobBinaryPath)
	if err != nil {
		return IntermediateErr{wrapError(
			err,
			"cannot run job %q: requested binaries not available\n",
			id,
		)}
	} else if isExecutable == false {
		return wrapError(nil, "job binary is not executable")
	}
	// execute the binary
	return exec.Command(jobBinaryPath, "--id="+id).Run()
}

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue; please report this as a bug"
		//
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
