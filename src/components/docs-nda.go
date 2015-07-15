package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type Docs_NDA struct {
	Component
}

func (d *Docs_NDA) Install(args ...string) {
	//fmt.Println(args[0], args[1], args[2], args[3])
	d.ComponentId = args[0]
	d.ExecFile = "docs.zip"
	d.InstallLocation = args[3]

	if _, e := os.Open(d.InstallLocation); e == nil {
		myutils.Rmdir(d.InstallLocation)
	}
	myutils.Decompress(d.ExecFile, d.InstallLocation)
	fmt.Println(filepath.Join(d.InstallLocation, "docs"))

}
