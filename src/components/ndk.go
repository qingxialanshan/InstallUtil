package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type NDK struct {
	Component
}

func (n *NDK) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	n.ComponentId = args[0]
	n.Version = args[1]
	n.ExecFile = args[2]
	n.InstallLocation = args[3]

	os.Chmod(n.ExecFile, 0777)
	_, e := exec.Command(n.ExecFile, "-y", "-o"+n.InstallLocation).Output()
	myutils.CheckError(e)

	//Add enviroment variables
	//1. add platform-tools,build-tools,tools to PATH

	ndk_path := filepath.Join(args[3], "android-ndk-r"+n.Version)
	//os.Mkdir(ndk_path, os.ModePerm)
	myutils.Set_Environment("PATH", ndk_path)

	//2. add ANDROID_NDK_HOME, NDKROOT,NDK_ROOT
	myutils.Set_Environment("ANDROID_NDK_ROOT", ndk_path)
	myutils.Set_Environment("NDK_ROOT", ndk_path)
	myutils.Set_Environment("NDKROOT", ndk_path)
	myutils.Set_Environment("NVPACK_NDK_VERSION", "android-ndk-r"+n.Version)
	//3. generate standalone toolchain

	if myutils.Global_OS == "linux" && runtime.GOARCH == "amd64" {
		genscript := filepath.Join(ndk_path, "build", "tools", "make-standalone-toolchain.sh")

		//generate 32 bit standalone toolchain
		s_path := filepath.Join(ndk_path, "toolchains", "arm-linux-androideabi-4.6", "gen_standalone", "linux-x86_64")
		gencmd := "bash " + genscript + " --platform=android-14 --system=linux-x86_64 --arch=arm --install-dir=" + s_path

		_, e1 := exec.Command("bash", "-c", gencmd).Output()
		myutils.CheckError(e1)
		myutils.Set_Environment("NDK_STANDALONE_46_ANDROID9_32", s_path)

		//generate 64 bit standalone toolchain
		s_path_64 := filepath.Join(ndk_path, "toolchains", "aarch64-linux-android-4.9", "gen_standalone", "linux-x86_64")
		gencmd_64 := "bash " + genscript + " --platform=android-21 --system=linux-x86_64 --arch=arm64 --install-dir=" + s_path_64 + " --toolchain=aarch64-linux-android-4.9"

		_, e2 := exec.Command("bash", "-c", gencmd_64).Output()
		myutils.CheckError(e2)
		myutils.Set_Environment("NDK_STANDALONE_46_ANDROID9_64", s_path_64)
	}
	fmt.Println(ndk_path)
}

func (n *NDK) Uninstall(args ...string) {

	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	n.ComponentId = args[0]
	n.Version = args[1]
	n.InstallLocation = args[2]
	location := filepath.Join(n.InstallLocation, "android-ndk-r"+n.Version)

	myutils.Delete_Environment("ANDROID_NDK_ROOT", location)
	myutils.Delete_Environment("NDK_ROOT", location)
	myutils.Delete_Environment("NDKROOT", location)
	myutils.Delete_Environment("PATH", location)
	myutils.Delete_Environment("NVPACK_NDK_VERSION", "")
	myutils.Delete_Environment("NDK_STANDALONE_46_ANDROID9_32", location)
	myutils.Delete_Environment("NDK_STANDALONE_46_ANDROID9_64", location)
	os.RemoveAll(location)
	fmt.Println(location)

}

func (n *NDK) Query(args ...string) string {
	n.ComponentId = args[0]
	n.InstallLocation = args[1]

	//get the revison from the RELEASE.TXT
	revision := myutils.Get_NDK_Revision(n.InstallLocation)

	return revision
}
