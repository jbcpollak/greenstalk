package util

import (
	"fmt"
	"strings"

	"github.com/jbcpollak/greenstalk/core"

	"github.com/fatih/color"
)

// NodeToString returns a string representation
// of a tree node and all its children.
func NodeToString[Blackboard any](node core.Node[Blackboard]) string {
	var b strings.Builder
	fmt.Println()

	appendToBuffer := func(node core.Walkable[Blackboard], level int) {
		appendNode[Blackboard](node, level, &b)
	}

	node.Walk(appendToBuffer, 0)
	appendNode(node, 0, &b)
	return b.String()
}

func appendNode[Blackboard any](node core.Walkable[Blackboard], level int, b *strings.Builder) {
	indent := strings.Repeat("    ", level)
	b.WriteString(indent + node.String() + "\n")
}

// PrintTreeInColor prints the tree with colors representing node state.
//
// Red = Failure, Yellow = Running, Green = Success, Magenta = Invalid.
func PrintTreeInColor[Blackboard any](node core.Node[Blackboard]) {
	node.Walk(printInColor, 0)
	fmt.Println()
}

func printInColor[Blackboard any](node core.Walkable[Blackboard], level int) {
	indent := strings.Repeat("    ", level)

	status := node.Status()
	var symbol string
	if status.IsErroneous() {
		color.Set(color.BgRed)
		symbol = "🚨"
	} else {
		color.Set(colorForStatus[node.Status()])
		symbol = symbolForStatus[node.Status()]
	}

	fmt.Println(indent + node.String() + " " + symbol)
	color.Unset()
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
