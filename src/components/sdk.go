package components

import (
	"fmt"
	//"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type SDK struct {
	Component
}

func (s *SDK) Install(args ...string) {
	s.ComponentId = args[0]

	myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//1. add platform-tools,build-tools,tools to PATH

	sdk_path := filepath.Join(args[3], "android-sdk-"+myutils.Global_OS)
	t_path := filepath.Join(sdk_path, "tools")
	myutils.Set_Environment("PATH", t_path)

	//2. add ANDROID_HOME
	myutils.Set_Environment("ANDROID_HOME", sdk_path)
	fmt.Println(sdk_path)
}

func (s *SDK) Uninstall(args ...string) {
	//fmt.Println("uninstalling sdk")
	s.ComponentId = args[0]
	s.Version = args[1]
	s.InstallLocation = filepath.Join(args[2], "android-sdk-"+myutils.Global_OS, "tools")

	installLocation := filepath.Join(args[2], "android-sdk-"+myutils.Global_OS)
	os.RemoveAll(s.InstallLocation)
	os.Remove(filepath.Join(installLocation, "AVD Manager.exe"))
	os.Remove(filepath.Join(installLocation, "SDK Manager.exe"))
	os.Remove(filepath.Join(installLocation, "SDK Readme.txt"))
	os.Remove(filepath.Join(installLocation, "add-ons"))
	os.Remove(filepath.Join(installLocation, "platforms"))
	os.Remove(filepath.Join(installLocation))
	myutils.Delete_Environment("ANDROID_HOME", installLocation)
	myutils.Delete_Environment("PATH", filepath.Join(installLocation, "tools"))
	fmt.Println(installLocation)
}

func (s *SDK) Query(args ...string) string {
	s.ComponentId = args[0]
	s.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(s.ComponentId, s.InstallLocation, "Pkg.Revision")
	return revision
}
