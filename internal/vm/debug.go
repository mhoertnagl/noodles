package vm

import "fmt"

func (m *VM) InspectStack(offset int64) Val {
	a := m.sp - offset - 1
	if a >= 0 {
		return m.stack[a]
	}
	return nil
}

func (m *VM) StackSize() int64 {
	return m.sp
}

func (m *VM) InspectFrames(offset int64) Val {
	a := m.fp + offset
	// Adresses above the FSP are invalid.
	if a >= 0 && a < m.fsp {
		return m.frames[a]
	}
	return nil
}

func (m *VM) FramesSize() int64 {
	return m.fsp
}

func (m *VM) printStack() {
	fmt.Print("STACK ⫣")
	for i := int64(0); i < m.sp; i++ {
		fmt.Printf(" %v │", m.stack[i])
	}
	fmt.Print("\n")
}

func (m *VM) printFrames() {
	w := m.fsp - m.fp + 2
	fmt.Print("FRAME ⫣")
	for i := int64(0); i < m.fsp; i++ {
		if i > 0 && i%w == w-1 {
			fmt.Printf(" %v │", m.frames[i])
		} else {
			fmt.Printf(" %v ⏐", m.frames[i])
		}
	}
	fmt.Print("\n")
}
