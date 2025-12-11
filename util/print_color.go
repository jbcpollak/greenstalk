package util

import (
	"fmt"
	"strings"

	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"

	"github.com/fatih/color"
)

// NodeToString returns a string representation
// of a tree node and all its children.
func NodeToString(node core.Node) string {
	var b strings.Builder
	fmt.Println()

	appendToBuffer := func(node core.Walkable, level int) {
		appendNode(node, level, &b)
	}

	node.Walk(appendToBuffer, 0)
	appendNode(node, 0, &b)
	return b.String()
}

func appendNode(node core.Walkable, level int, b *strings.Builder) {
	indent := strings.Repeat("    ", level)
	b.WriteString(indent + node.String() + "\n")
}

// PrintTreeInColor prints the tree with colors representing node state.
//
// Red = Failure, Yellow = Running, Green = Success, Magenta = Invalid.
func PrintTreeInColor(node core.Walkable) {
	node.Walk(printInColor, 0)
	fmt.Println()
}

func printInColor(node core.Walkable, level int) {
	indent := strings.Repeat("    ", level)

	status := node.Result().Status()
	color.Set(colorForStatus[status])
	symbol := symbolForStatus[status]

	fmt.Println(indent + node.String() + " " + symbol)
	color.Unset()
}

// Logs the tree to a logger
func PrintTreeToLog(node core.Walkable) {
	node.Walk(printToLog, 0)
}

func printToLog(node core.Walkable, level int) {
	indent := strings.Repeat("    ", level)

	status := node.Result().Status()
	symbol := symbolForStatus[status]

	internal.Logger.Info(indent + node.String() + " " + symbol)
}

var colorForStatus = map[core.Status]color.Attribute{
	core.StatusFailure: color.FgRed,
	core.StatusRunning: color.FgYellow,
	core.StatusSuccess: color.FgGreen,
	core.StatusInvalid: color.FgMagenta,
}

var symbolForStatus = map[core.Status]string{
	core.StatusFailure: "❌",
	core.StatusRunning: "🏃‍➡️",
	core.StatusSuccess: "✅",
	core.StatusInvalid: "❓",
}
