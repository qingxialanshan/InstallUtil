package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type Supportlib struct {
	Component
}

func (s *Supportlib) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	s.ComponentId = args[0]

	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//add platform-tools to PATH

	sdk_path := args[3]
	sp_path := filepath.Join(sdk_path, "support-tools")
	myutils.Set_Environment("PATH", sp_path)
	fmt.Println(sp_path)
}

func (s *Supportlib) Uninstall(args ...string) {

	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	s.ComponentId = args[0]
	s.Version = args[1]
	s.InstallLocation = args[2]

	os.RemoveAll(s.InstallLocation)

	myutils.Delete_Environment("PATH", filepath.Join(s.InstallLocation, "support-tools"))
	fmt.Println(filepath.Join(s.InstallLocation, "support-tools"))
}

func (s *Supportlib) Query(args ...string) string {
	s.ComponentId = args[0]
	s.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(s.ComponentId, filepath.Join(s.InstallLocation, "support"), "Pkg.Revision")
	return revision
}
