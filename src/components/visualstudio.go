package components

import (
	"myutils"
	"os/exec"
	"strconv"
)

type VisualStudio struct {
	Component
}

func (vs *VisualStudio) Install(args ...string) {

}

func (vs *VisualStudio) Uninstall(args ...string) {

}
func (vs *VisualStudio) Query(args ...string) string {
	vs.ComponentId = args[0]
	vs.InstallLocation = args[1]

	var wow, vs_reg string
	wow = ""
	if myutils.Get_os() == 1 {
		wow = "\\WOW6432Node"
	} else {
		wow = ""
	}
	reg := "HKEY_LOCAL_MACHINE\\Software" + wow + "\\Microsoft\\VisualStudio\\"

	version := []float64{12.0, 11.0, 10.0}

	vs_detect := ""
	for i := 0; i < 3; i++ {

		vspath := reg + strconv.FormatFloat(version[i], 'f', 1, 64)

		vs_reg = vspath + "\\Setup\\VS"
		vsdir, _ := exec.Command("reg", "query", vs_reg, "/v", "ProductDir").Output()

		if string(vsdir) != "" {
			vs_detect = strconv.FormatFloat(version[i], 'f', 1, 64)
			break
		}
		//		vsstring := strings.Split(string(vsdir), "\r\n")

		//		for str := range vsstring {
		//			if strings.Contains(vsstring[str], "InstallDir") {
		//				vs_str = vsstring[str]
		//				break
		//			}
		//		}
		//		vs_path := strings.Split(vs_str, "    ")
		//		//fmt.Println("0:", vs_path[0], "1:", vs_path[1], "2:", vs_path[2], "3:", vs_path[3])
		//		//vsdir = string(vsdir)
		//		if len(vs_path) < 4 {
		//			continue
		//		}
		//		verifypath := vs_path[3] + "devenv.com"

		//		if Exists(verifypath) {
		//			devenv_path = verifypath
		//			os.Chdir(string(vsdir))
		//			break
		//		}
	}
	//	rev := myutils.RegQuery(`Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "Nsight Tegra", "DisplayVersion")
	return vs_detect
}
