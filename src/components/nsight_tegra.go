package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type NsightTegra struct {
	Component
}

func (ns *NsightTegra) Install(args ...string) {

	ns.ComponentId = args[0]
	ns.Version = args[1]
	ns.ExecFile = args[2]
	ns.InstallLocation = args[3]

	//check visual studio installed or not
	var node string
	reglist := []string{"Microsoft\\VisualStudio\\10.0\\Setup\\VS", "Microsoft\\VisualStudio\\11.0\\Setup\\VS", "Microsoft\\VisualStudio\\12.0\\Setup\\VS"}
	if runtime.GOARCH == "amd64" {
		node = "HKLM\\Software\\Wow6432Node\\"
	} else {
		node = "HKLM\\Software\\"
	}
	var value string
	flag := false
	for i := 0; i < len(reglist); i++ {
		value = myutils.Get_val_reg(node+reglist[i], "ProductDir")
		if value == "" {
			flag = false
		} else {
			flag = true
			break
		}
	}

	if flag == false {
		fmt.Println("Cannot detect Visual Studio(2010|2012|2013) on your system. Please install Visual Studio first to support install Nsight Tegra")
		os.Exit(2)
	}
	_, e := exec.Command(ns.ExecFile, "/passive").Output()
	myutils.CheckError(e)
	fmt.Println("Nsight Tegra has integrated into Visual Studio")
}

func (ns *NsightTegra) Uninstall(args ...string) {

	ns.ComponentId = args[0]
	ns.Version = args[1]
	path := `Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`
	//uninstall the bundle exe
	//uninstall_exe_string := myutils.QueryRegister("Nsight", "QuietUninstallString")
	uninstall_exe_string := myutils.RegQuery(path, "Nsight Tegra", "QuietUninstallString")
	sysroot := os.Getenv("SystemRoot")
	uninstall_exe_string = strings.Replace(strings.Replace(uninstall_exe_string, "\"", "", -1), "/uninstall /quiet", "", 1)
	//fmt.Println("uninstall_exe_string is :", strings.Replace(uninstall_exe_string, "/uninstall /quiet", "", 1))

	if uninstall_exe_string != "" {

		_, e := exec.Command("cmd", "/c", uninstall_exe_string, "/uninstall", "/quiet").Output()
		myutils.CheckError(e)
	}

	//uninstall the msi if still exist

	//uninstall_msi_string := myutils.QueryRegister("NVIDIA*nsight*Tegra", "UninstallString")
	uninstall_msi_string := myutils.RegQuery(path, "Nsight Tegra", "UninstallString")

	if uninstall_msi_string != "" {

		uninstall_msi_string = strings.Replace(uninstall_msi_string, "MsiExec.exe ", "", 1)
		//fmt.Println(sysroot+"/System32/msiexec.exe", "/passive", uninstall_msi_string)
		_, e1 := exec.Command(sysroot+"/System32/msiexec.exe", "/passive", uninstall_msi_string).Output()
		if e1 != nil {
			fmt.Println(e1)
		}

		//fmt.Println("end")
	}
}

func (ns *NsightTegra) Query(args ...string) string {
	ns.ComponentId = args[0]
	ns.InstallLocation = args[1]
	//revision := myutils.QRVersion("NsightTegra")
	//revision := myutils.QueryRegister("Nsight Tegra", "DisplayVersion")
	rev := myutils.RegQuery(`Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "Nsight Tegra", "DisplayVersion")
	return rev
}
