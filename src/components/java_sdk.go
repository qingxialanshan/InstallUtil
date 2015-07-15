package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
)

type JavaSDK struct {
	Component
}

func (j *JavaSDK) Install(args ...string) {
	j.ComponentId = args[0]
	j.Version = args[1]
	j.InstallLocation = args[3]

	if _, e1 := os.Open(filepath.Join(j.InstallLocation, "jdk"+j.Version)); e1 == nil {
		//destination directory exist
		myutils.Rmdir(filepath.Join(j.InstallLocation, "jdk"+j.Version))
	}
	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//1. add java/bin to PATH
	j.Version = args[1]
	java_path := filepath.Join(args[3], "jdk"+args[1])
	bin_path := filepath.Join(java_path, "bin")
	myutils.Set_Environment("PATH", bin_path)

	//2. add ANT_HOME
	myutils.Set_Environment("JAVA_HOME", java_path)
	fmt.Println(java_path)
}

func (j *JavaSDK) Uninstall(args ...string) {
	//fmt.Println("uninstalling java")
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	j.ComponentId = args[0]
	j.Version = args[1]
	j.InstallLocation = filepath.Join(args[2], "jdk"+args[1])

	exec.Command("cmd", "/c", "del", "/F", "/S", "/Q", "/A:R", j.InstallLocation).Run()
	os.RemoveAll(j.InstallLocation)

	myutils.Delete_Environment("JAVA_HOME", j.InstallLocation)
	myutils.Delete_Environment("PATH", filepath.Join(j.InstallLocation, "bin"))
	fmt.Println(j.InstallLocation)
}
