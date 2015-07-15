package myutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Query_reg(regname, description string) string {

	var reg, revision string

	reg = ""

	//common_reg := "HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows\\CurrentVersion\\Installer\\UserData\\S-1-5-18\\Products"

	cmd_reg := os.Getenv("SystemRoot") + "\\system32\\reg.exe"

	output, e := exec.Command(cmd_reg, "query", regname, "/s", "/f", description).Output()
	//fmt.Println(description, cmd_reg, "query", regname, "/s", "/f", description, string(output))

	if e != nil {

		//fmt.Print("query failed")

		return ""

	} else {

		revision = string(output)

	}

	rv := strings.Split(revision, "\r\n")

	for i := 0; i < len(rv); i++ {

		if rv[i] != "" {

			reg = rv[i]

			break

		}

	}

	return reg

}

func Get_val_reg(reg, name string) string {

	var rev string

	value := ""

	cmd_reg := os.Getenv("SystemRoot") + "\\system32\\reg.exe"

	output, e := exec.Command(cmd_reg, "query", reg, "/s", "/f", name).Output()

	if e != nil {

		//fmt.Println(e)

	} else {

		rev = string(output)

	}

	vlist := strings.Split(rev, "\r\n")

	for i := 0; i < len(vlist); i++ {

		if strings.Contains(vlist[i], name) {

			drev := strings.Split(vlist[i], "  ")

			value = drev[len(drev)-1]

			break

		}

	}

	return value

}

func QueryRegister(name, value string) string {

	var regvalue, common_reg string

	if name == "Nsight" {
		common_reg = "HKEY_LOCAL_MACHINE\\Software\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall"
		reg := Query_reg(common_reg, "Nsight*Tegra*exe*uninstall*quiet")
		if reg == "" {
			return reg
		}
		regvalue = Get_val_reg(reg, "Nsight_Tegra_Installer_Bundle_x86_Release.exe")
		//reg := Query_reg(common_reg, "{283EB5AF-C351-42B4-9781-E0AB8B330BBA}")
		//regvalue = Get_val_reg(reg, value)
		//fmt.Println("the regvalue is ", regvalue)

	} else {

		common_reg = "HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows\\CurrentVersion\\Installer\\UserData\\S-1-5-18\\Products"
		reg := Query_reg(common_reg, name)
		if reg == "" {
			return reg
		}
		regvalue = Get_val_reg(reg, value)
	}

	return regvalue

}

func QRVersion(name string) string {
	var value string

	//common_reg := "\"HKEY_LOCAL_MACHINE\\Software\\Wow6432Node\\NVIDIA Corporation\\Registration\\" + name + "\""
	common_reg := "HKEY_LOCAL_MACHINE\\Software\\Wow6432Node\\NVIDIA Corporation\\Registration\\" + name
	cmd_reg := os.Getenv("SystemRoot") + "\\system32\\reg.exe"
	reg, e := exec.Command(cmd_reg, "query", common_reg, "/v", "ProductVer").Output()
	//fmt.Println(string(reg))
	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}

	vlist := strings.Split(string(reg), "\r\n")
	//fmt.Println("vlist is ", vlist, len(vlist))
	//fmt.Println(vlist[0], "\n", vlist[1], "\n", vlist[2], "\n", vlist[3], "\n", vlist[4])
	for i := 0; i < len(vlist); i++ {

		if strings.Contains(vlist[i], "ProductVer") {
			//fmt.Println("@@@", vlist[i])

			drev := strings.Split(vlist[i], " ")

			value = drev[len(drev)-1]

			break

		}

	}

	return value
}

func Get_NDK_Revision(flocation string) string {
	release_txt := get_source_properties(flocation, "RELEASE.TXT")
	revision := Readfile(release_txt)
	return strings.Split(revision, " ")[0]

}

func Get_Reg(description string) string {
	common_reg := "HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows\\CurrentVersion\\Installer\\UserData\\S-1-5-18\\Products"
	cmd_reg := os.Getenv("SystemRoot") + "system32\reg.exe"
	output, e := exec.Command(cmd_reg, "query", common_reg, "/s", "/f", description, ">x.txt").Output()
	if e != nil {
		//fmt.Println(output)
		return string(output)
	} else {
		return "cannot find"
	}
}

func Get_By_Tags(pkgname, location, tags string) string {

	var conf, res string

	if strings.Contains(pkgname, "api") {
		location = filepath.Join(location, strings.Replace(pkgname, "api", "android-", 1))
		//fmt.Println(location)
	} else if pkgname == "sdkbase" {
		location = filepath.Join(location, "android-sdk-"+Global_OS, "tools")
		//fmt.Println(location)
	} else if pkgname == "platformtools" {
		location = filepath.Join(location, "platform-tools")
	} else {
		location = location
	}

	conf = get_source_properties(location, "source.properties")
	if conf == "none" {
		os.Exit(-1)
		return "none"
	}
	//fmt.Println("the config file is :", conf)
	result := Readfile(conf)
	//fmt.Println(result)

	result = strings.Replace(result, "\n\n", "", 3)
	tmp_r := strings.Split(result, "\n")
	//fmt.Println(len(tmp_r))

	var i int
	for i = 0; i < len(tmp_r); i++ {
		//fmt.Println(i, tmp_r[i])
		if strings.Contains(tmp_r[i], tags) {
			res = tmp_r[i]
			break
		}
	}
	if i == len(tmp_r) {
		return "cannot find"
	}
	//fmt.Println(res)
	return strings.Split(res, "=")[1]
}

var fpath string

func get_source_properties(path, fname string) string {
	files, e := ioutil.ReadDir(path)
	if e != nil {
		return "none"
	}

	var flag bool
	flag = false
	for _, file := range files {
		if file.Name() == fname {
			flag = true
			//fmt.Println("find the path is: ", filepath.Join(path+"/"+file.Name()))
			fpath = filepath.Join(path + "/" + file.Name())
		}
	}
	if flag == false {
		for _, folder := range files {
			if folder.IsDir() {
				//fmt.Println(folder)
				new_path := filepath.Join(path, folder.Name())
				//fmt.Println(new_path)
				get_source_properties(new_path, fname)

			}
		}
	}

	return fpath
}

func Get_User_Path() string {

	path, e := exec.Command("reg", "query", "HKCU\\Environment", "-v", "path").Output()
	if e != nil {
		return ""
	}
	value := string(path)
	var tmp string
	vlist := strings.Split(value, "\n")
	for i := 0; i < len(vlist); i++ {
		if strings.Contains(vlist[i], "path") {
			tmp = vlist[i]
			break
		}
	}

	x := strings.Split(tmp, "    ")[3]
	return x
}
