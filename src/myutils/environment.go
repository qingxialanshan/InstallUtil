package myutils

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var Global_OS = Get_OS()
var Global_Home = os.Getenv("HOME")

func Get_OS() string {
	var os_name string
	if runtime.GOOS == "darwin" {
		os_name = "macosx"
	} else {
		os_name = runtime.GOOS
	}
	return os_name
}

func CheckError(e error) {
	if e != nil {

		fmt.Println(e)
		os.Exit(2)
	}
}

func Set_Environment(env_name, env_value string) {

	var set_command string
	if Global_OS == "windows" {
		if env_name == "PATH" {
			user_path := Get_User_Path()
			_, e := exec.Command(os.Getenv("SystemRoot")+"\\System32\\setx.exe", "PATH", env_value+";"+user_path).Output()
			CheckError(e)
		} else {
			set_command = "%SystemRoot%\\System32\\setx.exe " + env_name + " " + env_value
			_, e := exec.Command(os.Getenv("systemRoot")+"\\System32\\setx.exe", env_name, env_value).Output()

			CheckError(e)
		}

	} else if Global_OS == "linux" {

		config_file := "/.bashrc"

		if env_name == "PATH" {
			set_command = "\nexport PATH=$PATH:\"" + env_value + "\"\n"
			Write_To_File(set_command, Global_Home+config_file)

		} else {
			Delline(env_name, env_value, Global_Home+config_file)
			set_command = "\nexport " + env_name + "=\"" + env_value + "\"\n"
			Write_To_File(set_command, Global_Home+config_file)
		}

	} else if Global_OS == "macosx" {
		config_file := "/.bash_profile"

		if env_name == "PATH" {
			set_command = "\nexport PATH=$PATH:\"" + env_value + "\"\n"
			Write_To_File(set_command, Global_Home+config_file)

		} else {
			Delline(env_name, env_value, Global_Home+config_file)
			set_command = "\nexport " + env_name + "=\"" + env_value + "\"\n"
			Write_To_File(set_command, Global_Home+config_file)
		}
	}
}

func Delete_Environment(env_name, env_value string) {
	var del_command string
	if Global_OS == "windows" {
		if env_name == "PATH" {
			old_path := Get_User_Path()
			new_path := strings.Replace(old_path, env_value+";", "", -1)

			_, e := exec.Command(os.Getenv("SystemRoot")+"\\System32\\setx.exe", "PATH", new_path).Output()
			CheckError(e)
		} else {
			del_command = os.Getenv("SystemRoot") + "\\System32\\reg.exe delete HKCU\\Environment /F /V " + env_name
			if os.Getenv(env_name) != env_value {
				//fmt.Println("skip the ", env_name, env_value)
				return
			}
			_, e := exec.Command("cmd", "/c", del_command).Output()
			CheckError(e)
		}
	} else {
		var config_file string
		if Global_OS == "linux" {
			config_file = "/.bashrc"
		} else {
			config_file = "/.bash_profile"
		}
		if env_name == "PATH" {
			del_command = env_value
			//fmt.Println(del_command, env_value)
			Delline(del_command, env_value, Global_Home+config_file)
		} else {
			Delline(env_name, env_value, Global_Home+config_file)
		}
	}
}

func OS_Win() string {
	output, e := exec.Command("wmic", "os", "get", "Caption").Output()
	CheckError(e)
	var os string
	if strings.Contains(string(output), "Windows 7") || strings.Contains(string(output), "Windows XP") {
		os = "win7"
	} else {
		os = "win8"
	}
	return os
}

type EnvaddCommand struct {
}

func pre_action() {
	cmd := exec.Command("bash", "-c", `adb wait-for-device && adb shell svc power stayon true && adb root && adb wait-for-devices`)
	Redirector(cmd)
}

func set_adb_env(curr_path string) {
	path := os.Getenv("PATH")
	sdk_path := "/android-sdk-" + Global_OS
	var split string
	if Global_OS == "windows" {
		split = ";"
	} else {
		split = ":"
	}
	new_path := path + split + curr_path + sdk_path + "/platform-tools/"

	os.Setenv("PATH", new_path)
	//fmt.Println(new_path)
}
func (*EnvaddCommand) Run(args ...string) {
	fset := flag.NewFlagSet(args[0], flag.ExitOnError)
	err := fset.Parse(args[1:])
	if nil == err {
		if len(args) < 2 {
			fmt.Println("Please supply excute script file!")
			os.Exit(2)
		}
		set_adb_env(args[0])
		dir := filepath.Dir(args[0])
		os.Chdir(dir)
		//fmt.Println(os.Getenv("PATH"))

		pre_action()

		var cmd *exec.Cmd
		if Global_OS == "windows" {
			cmd = exec.Command("cmd", "/c", args[1])
		} else {
			cmd = exec.Command("sh", args[1])
		}

		Redirector(cmd)
		fmt.Println("end")

	} else {
		fmt.Println(err)
	}
}
