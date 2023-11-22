package main

import "fmt"

var (
	ErrInvalidInstanceType = fmt.Errorf("instance is not a Service nor a Singleton")
	ErrNotAStruct          = fmt.Errorf("instance is not a struct")
	ErrNotAPointer         = fmt.Errorf("instance is not a pointer")
	ErrTagInvalid          = fmt.Errorf("for tag \"di\" are only \"transient\" and \"direct\" (default) allowed")
	ErrInstanceNotFound    = fmt.Errorf("instance not found")
)
