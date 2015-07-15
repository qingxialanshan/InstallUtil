package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	//	"strings"
)

type Eclipse struct {
	Component
}

func get_java_version() string {
	var value string
	hkey := "HKCU\\Environment"
	output, e := exec.Command("reg", "query", hkey, "/v", "JAVA_HOME").Output()
	if e != nil {
		fmt.Println("Please install JAVA first to continue")
		os.Exit(2)
	}
	rev := string(output)
	vlist := strings.Split(rev, "\r\n")
	for i := 0; i < len(vlist); i++ {

		if strings.Contains(vlist[i], "jdk") {

			drev := strings.Split(vlist[i], "  ")

			value = drev[len(drev)-1]

			break

		}

	}
	return value
}

func (e *Eclipse) Install(args ...string) {

	e.ComponentId = args[0]
	e.Version = args[1]
	e.ExecFile = args[2]
	e.InstallLocation = args[3]

	var cfile, tmp_file string
	if myutils.Global_OS == "macosx" {
		//exec.Command("bash", "-c", "java -version 2>a.txt").Output()
		//		java_version, _ := exec.Command("cat", ".java.txt").Output()
		//		os.Remove("a.txt")
		//		fmt.Println(string(java_version))
		//		if strings.Contains(string(java_version), "1") {

		//		} else {
		//			fmt.Println("installing java 1.6.0_65")

		//			//_, oe1 := exec.Command("hdiutil", "attach", "_installer/JavaForOSX2014-001.dmg").Output()
		//			//_, oe2 := exec.Command("osascript", "-e 'do shell script with administrator privileges'", "installer -pkg /Volumes/Java for OS X 2014-001/JavaForOSX.pkg -target /").Output()
		//			_, oe2 := exec.Command("bash", "installjava.sh").Output()
		//			//fmt.Println("osascript", "-e 'do shell script with administrator privileges'", "installer -pkg /Volumes/Java for OS X 2014-001/JavaForOSX.pkg -target /")
		//			if oe2 != nil {

		//				myutils.CheckError(oe2)
		//			}
		//			//exec.Command("hdiutil", "detach", "/Volumes/Java for OS X 2014-001").Output()
		//			//os.Remove("installjava.sh")
		//		}

		//exec.Command("bash", "-c", "unzip "+e.ExecFile+" -d "+e.InstallLocation).Output()
		//		myutils.Decompress(e.ExecFile, e.InstallLocation)
		exec.Command("bash", "-c", "cd "+e.InstallLocation+"&& tar xvf "+e.ExecFile).Output()
		cfile = filepath.Join(e.InstallLocation, "eclipse", "Eclipse.app", "Contents", "MacOS", ".eclipse.ini")
		tmp_file = filepath.Join(e.InstallLocation, "eclipse", "Eclipse.app", "Contents", "MacOS", ".tmp_file")
	} else {
		myutils.Decompress(e.ExecFile, e.InstallLocation)
		cfile = filepath.Join(e.InstallLocation, "eclipse", "eclipse.ini")
		tmp_file = filepath.Join(e.InstallLocation, "eclipse", ".tmp_file")
	}

	var oldstr, newstr string
	if myutils.Global_OS == "windows" {

		//Get JAVA version
		java_ver := get_java_version()
		oldstr = "-vmargs"
		newstr = "-pluginCustomization\r\n" + e.InstallLocation + "eclipse\\nvpref.ini\r\n-vm\r\n" + java_ver + "\\bin\\javaw.exe\r\n-vmargs"

	} else if myutils.Global_OS == "linux" {
		oldstr = "-vmargs"
		newstr = "-pluginCustomization\n" + e.InstallLocation + "eclipse/nvpref.ini\n-vmargs"
	} else if myutils.Global_OS == "macosx" {
		oldstr = "-vmargs"
		newstr = "-pluginCustomization\n" + e.InstallLocation + "eclipse/Eclipse.app/Contents/Mac"
	}
	myutils.Substitude(cfile, oldstr, newstr)

	err := os.Rename(tmp_file, cfile)
	//_, err := myutils.CopyFile(tmp_file, cfile)
	if err != nil {
		fmt.Println("Copy file error")
		fmt.Println(err)
	}
	os.Remove(tmp_file)
	fmt.Println(filepath.Join(e.InstallLocation, "eclipse"))
	//	myutils.CopyFile(filepath.Join(e.InstallLocation, "eclipse", ".tmp_file"), cfile)
	//strings.Replace()

}

func (e *Eclipse) Uninstall(args ...string) {
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}

	e.ComponentId = args[0]
	e.Version = args[1]
	e.InstallLocation = filepath.Join(args[2], "eclipse")
	err := os.RemoveAll(e.InstallLocation)
	myutils.CheckError(err)
	fmt.Println(e.InstallLocation)
}
