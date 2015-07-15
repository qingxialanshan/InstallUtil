package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
)

type Samples struct {
	Component
}

func (s *Samples) Install(args ...string) {
	s.ComponentId = args[0]
	s.Version = args[1]
	s.ExecFile = args[2]
	s.InstallLocation = filepath.Join(args[3], "TDK_Samples")

	_, e := os.Stat(s.InstallLocation)
	if e == nil {
		os.RemoveAll(s.InstallLocation)
	}
	myutils.Decompress(s.ExecFile, args[3])
	oldfname := filepath.Join(args[3], "tegra_android_native_samples_"+s.Version)
	//os.Rename(oldfname, s.InstallLocation)
	myutils.CopyFile(oldfname, s.InstallLocation)
	fmt.Println(s.InstallLocation)
}

func (s *Samples) Uninstall(args ...string) {
	s.ComponentId = args[0]
	s.Version = args[1]
	s.InstallLocation = filepath.Join(args[2], "TDK_Samples")
	os.RemoveAll(s.InstallLocation)
	os.Remove(args[2])
	fmt.Println(s.InstallLocation)
}
