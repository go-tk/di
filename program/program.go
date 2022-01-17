// Package program represents a global Program.
package program

import "github.com/go-tk/di"

var p di.Program

// Exported methods of the global Program.
var (
	AddFunctions     = p.AddFunctions
	MustAddFunctions = p.MustAddFunctions
	Run              = p.Run
	MustRun          = p.MustRun
	Clean            = p.Clean
)
