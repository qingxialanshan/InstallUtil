package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type Gradle struct {
	Component
}

func (g *Gradle) Install(args ...string) {
	g.ComponentId = args[0]

	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//1. add gradle/bin to PATH
	g.Version = args[1]
	gradle_path := filepath.Join(args[3], "gradle-"+args[1])
	bin_path := filepath.Join(gradle_path, "bin")
	myutils.Set_Environment("PATH", bin_path)

	//2. add ANT_HOME
	myutils.Set_Environment("GRADLE_HOME", gradle_path)
	fmt.Println(gradle_path)
}

func (g *Gradle) Uninstall(args ...string) {

	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	g.ComponentId = args[0]
	g.Version = args[1]
	g.InstallLocation = filepath.Join(args[2], "gradle-"+args[1])

	os.RemoveAll(g.InstallLocation)
	myutils.Delete_Environment("GRADLE_HOME", g.InstallLocation)
	myutils.Delete_Environment("PATH", filepath.Join(g.InstallLocation, "bin"))
	fmt.Println(g.InstallLocation)
}
