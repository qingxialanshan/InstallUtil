package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	//	"os"
)

type Rebel struct {
	Component
}

func (r *Rebel) Install(args ...string) {
	fmt.Println(args[0], args[1], args[2], args[3])
	r.ComponentId = args[0]
	r.ExecFile = "TOD.zip"
	r.InstallLocation = args[3]

	//if _, e := os.Stat(r.InstallLocation); e != nil {
	//	os.Mkdir(r.InstallLocation, 0777)
	//}

	myutils.Decompress(r.ExecFile, r.InstallLocation)

}

func (r *Rebel) Uninstall(args ...string) {
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	r.ComponentId = args[0]
	r.Version = args[1]
	r.InstallLocation = args[2]
	fmt.Println(r.ComponentId, "uninstalling from : ", r.InstallLocation)
	if myutils.Global_OS == "windows" {
		exec.Command("cmd", "/c", "del", "/s", "/f", "/q", r.InstallLocation).Output()
	} else {
		exec.Command("rm", "-rf", r.InstallLocation).Output()
	}
	os.RemoveAll(r.InstallLocation)

}
