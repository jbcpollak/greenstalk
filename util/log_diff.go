package util

import (
	"fmt"
	"strings"

	"github.com/jbcpollak/greenstalk/v2/core"
	"github.com/jbcpollak/greenstalk/v2/internal"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type TreeLogger interface {
	LogTree(node core.Walkable)
}

func MakeDiffingTreeLogger() TreeLogger {
	differ := diffmatchpatch.New()
	return &diffingTreeLogger{
		differ: differ,
	}
}

type diffingTreeLogger struct {
	differ   *diffmatchpatch.DiffMatchPatch
	lastTree string
}

func (d *diffingTreeLogger) LogTree(node core.Walkable) {
	treeStringBuilder := strings.Builder{}

	printToBuilder := func(node core.Walkable, level int) {
		indent := strings.Repeat("    ", level)

		status := node.Result().Status()
		symbol := symbolForStatus[status]

		fmt.Fprintf(&treeStringBuilder, "%s%s %s\n", indent, node.String(), symbol)
	}

	node.Walk(printToBuilder, 0)

	currentTree := treeStringBuilder.String()

	if d.lastTree == "" {
		internal.Logger.Info(currentTree)
	} else {
		diff := d.differ.DiffMain(d.lastTree, currentTree, true)
		delta := d.differ.DiffToDelta(diff)
		internal.Logger.Info(delta)
	}

	d.lastTree = currentTree
}
