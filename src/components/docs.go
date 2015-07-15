package components

import (
	"fmt"
	"path/filepath"
	//	"io/ioutil"
	"myutils"
	"os"
	//"path/filepath"
)

type Docs struct {
	Component
}

func (d *Docs) Uninstall(args ...string) {
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	d.ComponentId = args[0]
	d.Version = args[1]
	d.InstallLocation = args[2]
	dir := d.InstallLocation
	myutils.Chmod(dir)

	os.RemoveAll(d.InstallLocation)

}

func (d *Docs) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}

	d.Version = args[1]
	d.ExecFile = args[2]
	d.InstallLocation = args[3]

	if _, e := os.Open(d.InstallLocation); e == nil {
		myutils.Rmdir(d.InstallLocation)
	}

	myutils.Decompress(args[2], args[3])
	fmt.Println(filepath.Join(d.InstallLocation, "docs"))
}
