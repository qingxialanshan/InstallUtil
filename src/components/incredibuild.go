package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"

	"os/exec"
)

type IncrediBuild struct {
	Component
}

func (ib *IncrediBuild) Install(args ...string) {

	ib.ComponentId = args[0]
	ib.Version = args[1]
	ib.ExecFile = args[2]
	ib.InstallLocation = filepath.Join(os.Getenv("ProgramFiles(x86)"), "Xoreax", "IncrediBuild")

	_, e := exec.Command("cmd", "/c", ib.ExecFile+` /install /Components=Coordinator,Agent /tadp_bundle`).Output()

	myutils.CheckError(e)

	fmt.Println(ib.InstallLocation)
}

func (ib *IncrediBuild) Uninstall(args ...string) {
	//InstallUtil uninstall incredibuild 1606 [installdir] -p [exefile]
	ib.ComponentId = args[0]
	ib.Version = args[1]
	//ib.ExecFile = args[2] //IBSetupConsole1606.exe
	ib.InstallLocation = filepath.Join(os.Getenv("ProgramFiles(x86)"), "Xoreax", "IncrediBuild")

	if len(args) > 4 && args[3] == "-p" {
		ib.ExecFile = args[4]
	}

	_, e := exec.Command(ib.ExecFile, "/uninstall").Output()
	myutils.CheckError(e)
	fmt.Println(ib.InstallLocation)

}

func (ib *IncrediBuild) Query(args ...string) string {
	ib.ComponentId = args[0]
	ib.InstallLocation = args[1]

	revision := myutils.RegQuery(`Software\Wow6432Node\Xoreax\IncrediBuild`, "Builder", "VersionText")

	return revision
}
