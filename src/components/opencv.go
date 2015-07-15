package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
)

type OpenCV struct {
	Component
}

func (ocv *OpenCV) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}

	ocv.Version = args[1]
	ocv.ExecFile = args[2]
	ocv.InstallLocation = args[3]

	if _, e := os.Open(filepath.Join(ocv.InstallLocation, "OpenCV-"+ocv.Version+"-Tegra-sdk")); e == nil {
		myutils.Rmdir(filepath.Join(ocv.InstallLocation, "OpenCV-"+ocv.Version+"-Tegra-sdk"))
	}

	myutils.Decompress(args[2], args[3])
	fmt.Println(filepath.Join(ocv.InstallLocation, "OpenCV-"+ocv.Version+"-Tegra-sdk"))
}

func (ocv *OpenCV) Uninstall(args ...string) {

	ocv.ComponentId = args[0]
	ocv.Version = args[1]
	ocv.InstallLocation = args[2]
	dir := filepath.Join(ocv.InstallLocation, "OpenCV-"+ocv.Version+"-Tegra-sdk")

	exec.Command("cmd", "/c", "del", "/F", "/S", "/Q", "/A:R", dir).Run()

	os.RemoveAll(dir)
	fmt.Println(dir)

}
