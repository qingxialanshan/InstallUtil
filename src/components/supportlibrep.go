package components

import (
	"fmt"
	"myutils"
	"path/filepath"
	"strings"
)

type SupportlibRep struct {
	Component
}

func (sr *SupportlibRep) Install(args ...string) {

	sr.ComponentId = args[0]
	sr.Version = args[1]
	sr.ExecFile = args[2]
	sr.InstallLocation = args[3]

	myutils.Decompress(args[2], args[3])
	fmt.Println(filepath.Join(sr.InstallLocation, "m2repository"))
}

func (s *SupportlibRep) Query(args ...string) string {
	s.ComponentId = args[0]
	s.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(s.ComponentId, filepath.Join(s.InstallLocation, "m2repository"), "Pkg.Revision")
	return strings.Split(revision, ".")[0]
}
