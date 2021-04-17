// +build tool

package tagscanner

import (
	_ "github.com/Songmu/gocredits"
)

//go:generate gocredits . -w ./CREDITS
