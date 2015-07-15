package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Quadd struct {
	Component
}

func (q *Quadd) Install(args ...string) {

	q.ComponentId = args[0]
	q.Version = args[1]
	q.ExecFile = args[2]
	q.InstallLocation = filepath.Join(args[3], "Tegra_System_Profiler")

	if _, err := os.Stat(filepath.Join(args[3], "_installer")); err != nil {
		os.MkdirAll(filepath.Join(args[3], "_installer"), 0777)
	}
	if myutils.Global_OS == "windows" {
		sysroot := os.Getenv("SystemRoot")
		//_, e := exec.Command(sysroot+"/System32/msiexec.exe", "/i", q.ExecFile, "/q", "/l*vx", "_installer/profile.log", "INSTALLDIR="+q.InstallLocation).Output()
		_, e := exec.Command(sysroot+"/System32/msiexec.exe", "/i", q.ExecFile, "/q", "/l*vx", "_installer/profile.log").Output()
		myutils.CheckError(e)
		q.InstallLocation = filepath.Join(os.Getenv("ProgramFiles"), "NVIDIA Corporation")

	} else if myutils.Global_OS == "linux" {
		myutils.Decompress(q.ExecFile, args[3])
		os.RemoveAll(q.InstallLocation)
		oldname := filepath.Base(strings.Replace(q.ExecFile, ".tar.gz", "", 1))
		oldname = filepath.Join(args[3], oldname)
		//os.Rename(oldname, q.InstallLocation)
		myutils.CopyFile(oldname, q.InstallLocation)

		//add to desktop entry
		quadd_entry := `[Desktop Entry]
Type=Application
Name=Tegra System Profiler
GenericName=Tegra System Profiler
Icon=` + args[3] + `
Exec=env PATH=usr/lib/lightdm/lightdm:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:` + filepath.Join(args[3], "android-sdk-linux", "platform-tool") + " " + filepath.Join(q.InstallLocation, "Host-x86_64", "TegraSystemProfiler") + `
TryExec=` + filepath.Join(q.InstallLocation, "Host-x86_64", "TegraSystemProfiler") + `
Keywords=tegra;nvidia;
Terminal=No
Categories=Development;`
		entry_folder := os.Getenv("HOME") + "/.local/share/applications"
		if _, se := os.Stat(entry_folder); se != nil {
			//fmt.Println("no folder exist")
			_, ce := exec.Command("mkdir", "-p", entry_folder).Output()
			if ce != nil {
				fmt.Println(ce)
			}
		}
		fp, fe := os.OpenFile(entry_folder+"/tegrasystemprofiler.desktop", os.O_CREATE|os.O_RDWR, 0600)
		if fe != nil {
			fmt.Println(fe)
		}
		fp.WriteString(quadd_entry)
	} else {
		myutils.Decompress(q.ExecFile, q.InstallLocation)
	}
	fmt.Println(q.InstallLocation)
}

func (q *Quadd) Uninstall(args ...string) {

	q.ComponentId = args[0]
	q.Version = args[1]
	q.InstallLocation = filepath.Join(args[2], "Tegra_System_Profiler")

	if myutils.Global_OS == "windows" {
		q.InstallLocation = filepath.Join(os.Getenv("ProgramFiles"), "NVIDIA Corporation")
		//uninstall_string := myutils.QueryRegister("tegra*system", "UninstallString")
		uninstall_string := myutils.RegQuery(`Software\Microsoft\Windows\CurrentVersion\Uninstall`, "Tegra System Profiler", "UninstallString")

		uninstall_string = strings.Replace(uninstall_string, "MsiExec.exe ", "", 1)
		sysroot := os.Getenv("SystemRoot")
		_, e := exec.Command(sysroot+"/System32/msiexec.exe", "/passive", uninstall_string).Output()
		myutils.CheckError(e)
	} else {
		if myutils.Global_OS == "linux" {
			os.Remove(os.Getenv("HOME") + "/.local/share/applications/tegrasystemprofiler.desktop")
		}
		os.RemoveAll(q.InstallLocation)
	}
	fmt.Println(q.InstallLocation)
}

func (q *Quadd) Query(args ...string) string {
	q.ComponentId = args[0]
	q.InstallLocation = args[1]
	//revision := myutils.QueryRegister("tegra*system", "DisplayVersion")
	revision := myutils.RegQuery(`Software\Microsoft\Windows\CurrentVersion\Uninstall`, "Tegra System Profiler", "DisplayVersion")

	return revision
}
func (q *Quadd) PreInstall(args ...string) int {
	if len(args) > 5 || len(args) < 4 {
		fmt.Println("Wrong args input")
		os.Exit(-1)
		return -1
	}
	q.ComponentId = args[0]
	q.Version = args[1]
	task := args[2]
	target := args[3]

	if len(args) == 5 {
		//check the condition is met or not
		if args[4] == "-c" {
			return myutils.CheckProcess(target)
		} else {
			os.Exit(-1)
			return -1
		}
	}
	//if the -c is off , task action will be performed
	return myutils.ManageProcess(task, target)
}
