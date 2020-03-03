package bin

import (
	"encoding/gob"
	"io/ioutil"
	"os"

	"github.com/mhoertnagl/splis2/internal/vm"
)

func ReadStatic(file *os.File) vm.Ins {
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return b
}

func ReadLib(file *os.File) *Lib {
	lib := &Lib{}
	gob.Register(Lib{})
	e := gob.NewDecoder(file)
	err := e.Decode(lib)
	if err != nil {
		panic(err)
	}
	return lib
}
