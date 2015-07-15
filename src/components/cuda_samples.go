package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type CUDASamples struct {
	Component
}

func (cs *CUDASamples) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}

	cs.Version = args[1]
	cs.ExecFile = args[2]
	cs.InstallLocation = args[3]

	if _, e := os.Open(cs.InstallLocation); e == nil {
		myutils.Rmdir(cs.InstallLocation)
	}

	myutils.Decompress(args[2], args[3])
	fmt.Println(filepath.Join(cs.InstallLocation, "CUDA_Samples"))
}
