package main

import (
	"github.com/mannemsolutions/pgsectest/internal"
)

func main() {
	internal.Initialize()
	internal.Handle()
}
