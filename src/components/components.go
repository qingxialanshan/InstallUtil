package components

import (
	"fmt"
	"myutils"
	"path/filepath"
)

type Component struct {
	ComponentId     string
	Version         string
	ExecFile        string
	InstallLocation string
}

func (c *Component) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}

	c.Version = args[1]
	c.ExecFile = args[2]
	//fmt.Println(c.ComponentId, "installing ", c.Version, "using execute file : ", c.ExecFile)
	myutils.Decompress(args[2], args[3])
	c.InstallLocation = filepath.Join(args[3], "")
	fmt.Println(c.InstallLocation)
}

func (c *Component) Uninstall(args ...string) {
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	c.ComponentId = args[0]
	c.Version = args[1]
	c.InstallLocation = filepath.Join(args[2], "")
	//fmt.Println(c.ComponentId, "uninstalling from : ", c.InstallLocation)
	//if myutils.Global_OS == "windows" {
	//	exec.Command("cmd", "/c", "del", "/s", "/f", "/q", c.InstallLocation).Output()
	//}

	myutils.Rmdir(c.InstallLocation)
	fmt.Println(c.InstallLocation)

}
