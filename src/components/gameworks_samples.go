package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type GameWorks struct {
	Component
}

func (g *GameWorks) Install(args ...string) {
	g.ComponentId = args[0]
	g.Version = args[1]
	g.ExecFile = args[2]
	g.InstallLocation = args[3]

	if _, e1 := os.Open(g.InstallLocation); e1 == nil {
		//destination directory exist
		myutils.Rmdir(g.InstallLocation)
	}

	//fmt.Println(c.ComponentId, "installing ", c.Version, "using execute file : ", c.ExecFile)
	myutils.Decompress(args[2], args[3])
	fmt.Println(filepath.Join(g.InstallLocation, ""))
}

func (g *GameWorks) Uninstall(args ...string) {
	g.ComponentId = args[0]
	g.Version = args[1]
	g.InstallLocation = args[2]

	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", "del", "/s", "/f", "/q", g.InstallLocation).Output()
		//		if e != nil {
		//			fmt.Println(e)
		//		}
		exec.Command("cmd", "/c", "del", "/f", "/q", filepath.Join(g.InstallLocation, "license.txt")).Run()

	}
	os.RemoveAll(g.InstallLocation)
	dir := strings.Replace(g.InstallLocation, "GameWorks_Samples", "", 1)
	os.Remove(dir)
	fmt.Println(filepath.Join(g.InstallLocation, ""))
}
