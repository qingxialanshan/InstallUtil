package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
)

type PlatformTools struct {
	Component
}

func (p *PlatformTools) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	p.ComponentId = args[0]

	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//add platform-tools to PATH

	sdk_path := args[3]
	pt_path := filepath.Join(sdk_path, "platform-tools")
	myutils.Set_Environment("PATH", pt_path)
	fmt.Println(pt_path)
}

func (p *PlatformTools) Uninstall(args ...string) {

	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	p.ComponentId = args[0]
	p.Version = args[1]
	p.InstallLocation = args[2]

	if myutils.Global_OS == "windows" {
		exec.Command("taskkill", "/f", "/IM", "adb.exe")

	}
	os.RemoveAll(filepath.Join(p.InstallLocation, "platform-tools"))
	//myutils.Delete_Environment("ANDROID_HOME", s.InstallLocation)
	myutils.Delete_Environment("PATH", filepath.Join(p.InstallLocation, "platform-tools"))
	fmt.Println(filepath.Join(p.InstallLocation, "platform-tools"))
}

func (p *PlatformTools) Query(args ...string) string {
	p.ComponentId = args[0]
	p.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(p.ComponentId, p.InstallLocation, "Pkg.Revision")
	return revision
}

func (p *PlatformTools) PreInstall(args ...string) int {
	if len(args) > 5 || len(args) < 4 {
		fmt.Println("Wrong args input")
		os.Exit(-1)
		return -1
	}
	p.ComponentId = args[0]
	p.Version = args[1]
	task := args[2]
	target := args[3]

	if len(args) == 5 {
		//check the condition is met or not
		if args[4] == "-c" {
			return myutils.CheckProcess(target)
		} else {
			os.Exit(-1)
			return -1
		}
	}
	//if the -c is off , task action will be performed
	return myutils.ManageProcess(task, target)
}
