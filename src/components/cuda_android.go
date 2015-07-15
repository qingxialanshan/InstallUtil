package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
	"strings"
)

type CUDAAndroid struct {
	Component
}

func (ca *CUDAAndroid) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}

	ca.Version = args[1]
	ca.ExecFile = args[2]
	ca.InstallLocation = args[3]

	if _, e := os.Open(ca.InstallLocation); e == nil {
		myutils.Rmdir(ca.InstallLocation)
	}

	//fmt.Println(c.ComponentId, "installing ", c.Version, "using execute file : ", c.ExecFile)
	myutils.Decompress(args[2], args[3])
	if strings.Contains(ca.Version, "6.5") {
		fmt.Println(filepath.Join(ca.InstallLocation, "cuda-android-6.5"))
	} else if strings.Contains(ca.Version, "7.0") {
		fmt.Println(filepath.Join(ca.InstallLocation, "cuda-android-7.0"))
	}

}
