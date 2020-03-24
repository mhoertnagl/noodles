package bin

import (
	"bufio"
	"encoding/gob"
	"os"

	"github.com/mhoertnagl/splis2/internal/vm"
)

func WriteStatic(code vm.Ins, file *os.File) {
	w := bufio.NewWriter(file)
	_, err := w.Write(code)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func WriteLib(lib *Lib, file *os.File) {
	gob.Register(Lib{})
	e := gob.NewEncoder(file)
	err := e.Encode(lib)
	if err != nil {
		panic(err)
	}
}
