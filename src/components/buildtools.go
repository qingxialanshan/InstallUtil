package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type BuildTools struct {
	Component
}

func (b *BuildTools) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	b.ComponentId = args[0]
	b.Version = args[1]
	bt_path := args[3]
	myutils.Decompress(args[2], bt_path)
	flocation := filepath.Join(bt_path, "android-5.1")

	//os.Rename(flocation, filepath.Join(bt_path, b.Version))
	myutils.CopyFile(flocation, filepath.Join(bt_path, b.Version))
	//Add enviroment variables
	//add build-tools to PATH

	myutils.Set_Environment("PATH", bt_path)
	fmt.Println(filepath.Join(bt_path, ""))
}

func (b *BuildTools) Uninstall(args ...string) {
	//fmt.Println("uninstalling sdk Build Tools")
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	b.ComponentId = args[0]
	b.Version = args[1]
	b.InstallLocation = filepath.Join(args[2], "")

	os.RemoveAll(b.InstallLocation)
	myutils.Delete_Environment("PATH", b.InstallLocation)
	fmt.Println(b.InstallLocation)
}

func (b *BuildTools) Query(args ...string) string {
	b.ComponentId = args[0]
	b.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(b.ComponentId, b.InstallLocation, "Pkg.Revision")
	return revision
}
