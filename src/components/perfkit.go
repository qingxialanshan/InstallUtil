package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
)

type Perfkit struct {
	Component
}

func (pk *Perfkit) Uninstall(args ...string) {
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}

	pk.ComponentId = args[0]
	pk.Version = args[1]
	pk.InstallLocation = filepath.Join(args[2], "PerfKit")

	if myutils.Global_OS == "windows" {
		exec.Command("C:\\Windows\\system32\\cmd.exe", "/c", "del", "/f", "/s", "/q", pk.InstallLocation).Output()
	}
	err := os.RemoveAll(pk.InstallLocation)
	myutils.CheckError(err)
	//fmt.Println(componentId, "uninstalling ", version, "using execute file : ", execFile)
}
