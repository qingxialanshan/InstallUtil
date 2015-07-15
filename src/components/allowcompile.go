package components

import (
	"myutils"
)

type AllowCompile struct {
	Component
}

func (ac *AllowCompile) Install(args ...string) {
	ac.ComponentId = args[0]
	ac.Version = args[1]
	ac.ExecFile = args[2]
	ac.InstallLocation = args[3]

	c := &myutils.CompileCommand{}
	//fmt.Println(ac.Version, ac.InstallLocation)
	c.Run("compile", ac.InstallLocation, ac.Version)
	return
}

func (ac *AllowCompile) Uninstall(args ...string) {

	return
}
