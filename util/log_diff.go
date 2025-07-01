package util

import (
	"fmt"
	"strings"

	"github.com/jbcpollak/greenstalk/core"
	"github.com/jbcpollak/greenstalk/internal"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type TreeLogger[Blackboard any] interface {
	LogTree(node core.Walkable[Blackboard])
}

func MakeDiffingTreeLogger[Blackboard any]() TreeLogger[Blackboard] {
	differ := diffmatchpatch.New()
	return &diffingTreeLogger[Blackboard]{
		differ: differ,
	}
}

type diffingTreeLogger[Blackboard any] struct {
	differ   *diffmatchpatch.DiffMatchPatch
	lastTree string
}

func (d *diffingTreeLogger[Blackboard]) LogTree(node core.Walkable[Blackboard]) {
	treeStringBuilder := strings.Builder{}

	printToBuilder := func(node core.Walkable[Blackboard], level int) {
		indent := strings.Repeat("    ", level)

		status := node.Result().Status()
		symbol := symbolForStatus[status]

		treeStringBuilder.WriteString(fmt.Sprintf("%s%s %s\n", indent, node.String(), symbol))
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
