package condition

import (
	"github.com/jbcpollak/greenstalk/v2/core"
)

// Switch activates the child at the index returned by the switch function.
// Note that an If node is a rudimentary form of a Switch node with two children
// and the function returning 0 / 1 for true / false.
func SwitchNamed(name string, switchFunc ParamSwitchFunc[int], children ...core.Node) core.Node {
	base := core.NewComposite(core.BaseParams(name), children)
	childrenMap := map[int]core.Node{}
	for i, child := range children {
		childrenMap[i] = child
	}
	return &paramSwitchNode[int]{
		Composite:  base,
		switchFunc: switchFunc,
		children:   childrenMap,
	}
}

func Switch(switchFunc ParamSwitchFunc[int], children ...core.Node) core.Node {
	return SwitchNamed("Switch", switchFunc, children...)
}
