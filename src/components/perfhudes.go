package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	//	"strings"
)

type PerfHUDES struct {
	Component
}

func (p *PerfHUDES) Install(args ...string) {

	p.ComponentId = args[0]
	p.Version = args[1]
	p.ExecFile = args[2]
	p.InstallLocation = args[3]

	var pkg_name, pkg_installer string
	_, exist := os.Stat(filepath.Join(p.InstallLocation, "perfhudes_tmp", "uninstall.iss"))
	if exist == nil {
		os.Chmod(filepath.Join(p.InstallLocation, "perfhudes_tmp", "uninstall.iss"), 0777)
		os.RemoveAll(filepath.Join(p.InstallLocation, "perfhudes_tmp"))
	}
	myutils.Decompress(args[2], args[3])
	if myutils.Global_OS == "linux" { //install perfhude on linux
		if runtime.GOARCH == "386" {
			pkg_name = "perfhudes_release.x86.deb.tgz"
			pkg_installer = "perfhudes_" + p.Version + "_i386.deb"
		} else if runtime.GOARCH == "amd64" {
			pkg_name = "perfhudes_release.x64.deb.tgz"
			pkg_installer = "perfhudes_" + p.Version + "_amd64.deb"
		} else {
			fmt.Println("unknown os arch")
			os.Exit(-1)
		}

		myutils.Decompress(filepath.Join(args[3], pkg_name), p.InstallLocation)

		os.Remove(filepath.Join(args[3], pkg_name))
		_, e := exec.Command("xterm", "-e", "sudo dpkg -i "+filepath.Join(args[3], pkg_installer)+" && sudo apt-get install -f -y && echo success>>.p.txt").Output()
		myutils.CheckError(e)
		os.Remove(filepath.Join(p.InstallLocation, "perfhudes_"+p.Version+"_amd64.deb"))
		if _, pe := os.Stat(".p.txt"); pe != nil {
			fmt.Println("Failed to install perfhudes")
			os.Exit(2)
		} else {
			os.Remove(".p.txt")
		}
	} else if myutils.Global_OS == "windows" {
		setup_file := filepath.Join(p.InstallLocation, "perfhudes_tmp", "setup.iss")
		setup_log := filepath.Join(p.InstallLocation, "perfhudes_tmp", "setup.log")
		myutils.Substitude(setup_file, "szDir=c:\\perfhud_es", "szDir="+args[3]+"\\perfhud_es")
		os.Remove(setup_file)
		//re := os.Rename(filepath.Join(args[3], "perfhudes_tmp", ".tmp_file"), setup_file)
		_, re := myutils.CopyFile(filepath.Join(args[3], "perfhudes_tmp", ".tmp_file"), setup_file)
		myutils.CheckError(re)
		pkg_name = "NVIDIA_PerfHUD_ES_Release.exe"
		pkg_installer = filepath.Join(p.InstallLocation, "perfhudes_tmp", pkg_name)
		_, e := exec.Command("cmd", "/c", pkg_installer, "/s", "/f1"+setup_file+"\n/f2"+setup_log).Output()
		myutils.CheckError(e)

		dst := filepath.Join(p.InstallLocation, "perfhud_es", "uninstall.iss")
		org := filepath.Join(p.InstallLocation, "perfhudes_tmp", "uninstall.iss")

		_, ce := myutils.CopyFile(org, dst)
		myutils.CheckError(ce)
		os.Chmod(org, 0777)
		os.RemoveAll(filepath.Join(p.InstallLocation, "perfhudes_tmp"))

	} else if myutils.Global_OS == "macosx" { //not finished
		perfhudes_mac := "perfhudes-" + p.Version
		pkg_installer = "perfhudes-" + p.Version + ".dmg"
		_, e1 := exec.Command("hdiutil", "attach", filepath.Join(args[3], "perfhudes_tmp", pkg_installer)).Output()
		myutils.CheckError(e1)
		_, e2 := exec.Command("cp", "-rf", "/Volumes/"+perfhudes_mac+"_x86/perfhudes.app", args[3]).Output()
		myutils.CheckError(e2)
		//detach
		_, e3 := exec.Command("hdiutil", "detach", "/Volumes/"+perfhudes_mac+"_x86").Output()
		myutils.CheckError(e3)
	} else {
		fmt.Println("error")
	}
	fmt.Println("System")
}

func (p *PerfHUDES) Uninstall(args ...string) {

	p.ComponentId = args[0]
	p.Version = args[1]
	p.InstallLocation = filepath.Join(args[2], "perfhud_es")

	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/usr/bin/perfhudes"); err == nil {
			_, e := exec.Command("xterm", "-e", "sudo dpkg -r perfhudes").Output()
			if e != nil {
				fmt.Println(e)
				os.Exit(2)
			}
		}
	} else if runtime.GOOS == "windows" { //has issues
		installer := filepath.Join(os.Getenv("ProgramFiles(x86)"), "InstallShield Installation Information", "{207C2D06-4314-4E2D-B72B-843FB7C6B745}", "setup.exe")

		if _, se := os.Stat(installer); se != nil {
			os.Exit(0)
		}
		uninstallstr := filepath.Join(p.InstallLocation, "uninstall.iss")
		_, e := exec.Command(installer, "/s", "/f1"+uninstallstr).Output()

		if e != nil {
			fmt.Println(e)
			os.Exit(2)
		}
	} else if runtime.GOOS == "darwin" {
		os.RemoveAll(filepath.Join(args[2], "perfhudes.app"))
	} else {
		fmt.Println("error")
	}
	fmt.Println("System")
}
