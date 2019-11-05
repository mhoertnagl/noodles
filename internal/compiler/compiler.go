package compiler

import "github.com/mhoertnagl/splis2/internal/vm"

type Compiler interface {
	Compile(node Node) vm.Ins
}

type compiler struct {
}

func NewCompiler() Compiler {
	return &compiler{}
}

func (c *compiler) Compile(node Node) vm.Ins {
	switch node.(type) {
	case bool:

	}
	return nil
}
