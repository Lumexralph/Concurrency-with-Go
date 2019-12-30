// Package main illustrates using context as a data storage or keep
// to store and retrive request-scoped data.
//
// some rules are suggested when using the WithValue context
// 1. Define a custom key type in your package, this prevents collisions
// with the Context.
package main

import "context"

import "fmt"

func main() {
	processRequest("lumex", "adc1234")
}

// custom type for the context key
type ctxkey int

const (
	ctxUserID ctxkey = 1 + iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	// assert and typecast to a string
	return c.Value(ctxUserID).(string)
}

func UserAuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func processRequest(userID, authToken string) {
	// store the user info in the context
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)

	// after adding value to the context, handle request
	HandleRequest(ctx)
}

func HandleRequest(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)\n",
		UserID(ctx),
		UserAuthToken(ctx),
	)
}
