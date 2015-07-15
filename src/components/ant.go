package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type Ant struct {
	Component
}

func (a *Ant) Install(args ...string) {
	a.ComponentId = args[0]

	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//1. add ant/bin to PATH
	a.Version = args[1]
	ant_path := filepath.Join(args[3], "apache-ant-"+args[1])
	bin_path := filepath.Join(ant_path, "bin")
	myutils.Set_Environment("PATH", bin_path)

	//2. add ANT_HOME
	myutils.Set_Environment("ANT_HOME", ant_path)
	fmt.Println(ant_path)
}

func (a *Ant) Uninstall(args ...string) {
	//fmt.Println("uninstalling ant")
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	a.ComponentId = args[0]
	a.Version = args[1]
	a.InstallLocation = filepath.Join(args[2], "apache-ant-"+args[1])

	os.RemoveAll(a.InstallLocation)
	myutils.Delete_Environment("ANT_HOME", a.InstallLocation)
	myutils.Delete_Environment("PATH", filepath.Join(a.InstallLocation, "bin"))
}
