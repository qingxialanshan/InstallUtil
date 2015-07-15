package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Battle struct {
	Component
}

func (b *Battle) Install(args ...string) {

	b.ComponentId = args[0]
	b.Version = args[1]
	b.ExecFile = args[2]
	b.InstallLocation = filepath.Join(args[3], "Tegra_Graphics_Debugger")

	var e error
	if _, err := os.Stat(filepath.Join(args[3], "_installer")); err != nil {
		os.MkdirAll(filepath.Join(args[3], "_installer"), 0777)
	}
	if myutils.Global_OS == "windows" {
		b.InstallLocation = filepath.Join(os.Getenv("ProgramFiles(x86)"), "NVIDIA Corporation")
		sysroot := os.Getenv("SystemRoot")
		//_, e = exec.Command(sysroot+"/System32/msiexec.exe", "/i", b.ExecFile, "/q", "/l*vx", "_installer/battle.log", "INSTALLDIR="+b.InstallLocation).Output()
		_, e = exec.Command(sysroot+"/System32/msiexec.exe", "/i", b.ExecFile, "/q", "/l*vx", "_installer/battle.log").Output()
	} else if myutils.Global_OS == "linux" {
		os.Chmod(b.ExecFile, 0777)
		//_, e = exec.Command("xterm", "-e", b.ExecFile+" -targetpath="+b.InstallLocation+" -noprompt").Output()
		_, e = exec.Command(b.ExecFile, " -targetpath="+b.InstallLocation, " -noprompt", "--nox11").Output()

		//add to desktop entry
		battle_entry := `[Desktop Entry]
Type=Application
Name=Tegra System Profiler
GenericName=Tegra System Profiler
Icon=` + args[3] + `
Exec=env PATH=usr/lib/lightdm/lightdm:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:` + filepath.Join(args[3], "android-sdk-linux", "platform-tool") +
			" " + filepath.Join(b.InstallLocation, "nvidia-gfx-debugger") + `
TryExec=` + filepath.Join(b.InstallLocation, "nvidia-gfx-debugger") + `
Keywords=tegra;graphics;nvidia;
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
		fp, fe := os.OpenFile(entry_folder+"/tegragraphicsdebugger.desktop", os.O_CREATE|os.O_RDWR, 0600)
		if fe != nil {
			fmt.Println(fe)
		}
		fp.WriteString(battle_entry)
	} else if myutils.Global_OS == "macosx" {
		if _, exist := os.Stat(b.InstallLocation); exist != nil {
			os.Mkdir(b.InstallLocation, 0777)
		}
		_, e1 := exec.Command("hdiutil", "attach", b.ExecFile).Output()

		_, e2 := exec.Command("cp", "-r", "/Volumes/NVIDIA Tegra Graphics Debugger/NVIDIA Tegra Graphics Debugger.app", b.InstallLocation).Output()
		exec.Command("cp", "/Volumes/NVIDIA Tegra Graphics Debugger/EULA.txt", b.InstallLocation).Output()

		_, e = exec.Command("hdiutil", "detach", "/Volumes/NVIDIA Tegra Graphics Debugger").Output()

		myutils.CheckError(e1)
		myutils.CheckError(e2)
	} else {
		fmt.Println("unknow platforms, only support Windows ,Linux and Mac")
		os.Exit(2)
	}

	myutils.CheckError(e)
	fmt.Println(b.InstallLocation)
}

func (b *Battle) Uninstall(args ...string) {

	b.ComponentId = args[0]
	b.Version = args[1]
	b.InstallLocation = filepath.Join(args[2], "Tegra_Graphics_Debugger")
	if myutils.Global_OS == "windows" {
		uninstall_string := myutils.RegQuery(`Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "NVIDIA Tegra Graphics Debugger", "UninstallString")
		//uninstall_string := myutils.QueryRegister("graphics*debugger", "UninstallString")
		//fmt.Println(uninstall_string)
		uninstall_string = strings.Replace(uninstall_string, "MsiExec.exe ", "", 1)
		sysroot := os.Getenv("SystemRoot")
		_, e := exec.Command(sysroot+"/System32/msiexec.exe", "/passive", uninstall_string).Output()
		if e != nil {
			fmt.Println(e)
			os.Exit(2)
		}
		b.InstallLocation = filepath.Join(os.Getenv("ProgramFiles(x86)"), "NVIDIA Corporation")

	} else {
		if myutils.Global_OS == "linux" {
			os.Remove(os.Getenv("HOME") + "/.local/share/applications/tegragraphicsdebugger.desktop")
		}
		os.RemoveAll(b.InstallLocation)
	}
	fmt.Println(b.InstallLocation)
}

func (b *Battle) Query(args ...string) string {
	b.ComponentId = args[0]
	b.InstallLocation = args[1]
	//revision := myutils.QueryRegister("graphics*debugger", "DisplayVersion")
	revision := myutils.RegQuery(`Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "NVIDIA Tegra Graphics Debugger", "DisplayVersion")
	return revision
}

func (b *Battle) PreInstall(args ...string) int {
	if len(args) > 5 || len(args) < 4 {
		fmt.Println("Wrong args input")
		os.Exit(-1)
		return -1
	}
	b.ComponentId = args[0]
	b.Version = args[1]
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
