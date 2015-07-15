package components

import (
	"fmt"
	"os"
	"os/exec"
	//	"strings"
	//	"fmt"
	"myutils"
	//	"os"
	"path/filepath"
)

type USB struct {
	Component
}

func (u *USB) Install(args ...string) {
	u.ComponentId = args[0]

	location := filepath.Join(args[3], "Drivers", "Devkits")
	//if os.IsExist(os.Chdir(location)) == false {
	//	os.Create(location)
	//}
	myutils.Decompress(args[2], location)
	install_location := filepath.Join(location, "Win7+Vista+XP")
	os.RemoveAll(install_location)
	//os.Rename(filepath.Join(location, "usb_driver"), install_location)
	myutils.CopyFile(filepath.Join(location, "usb_driver"), install_location)

	if myutils.OS_Win() == "win7" {
		//substitude the file Google.NTx86
		fname := filepath.Join(install_location, "android_winusb.inf")
		oldstr := "[Google.NTx86]"
		tegra_str := ";NVIDIA Tegra\r\n" + "%SingleAdbInterface%\t= USB_Install,USB\\VID_0955&PID_7000\r\n%CompositeAdbInterface%\t= USB_Install,USB\\VID_0955&PID_7100&MI_01\r\n\r\n"
		note_str := ";NVIDIA Tegra\r\n%SingleAdbInterface%\t= USB_Install, USB\\VID_0955&PID_CF01\r\n%CompositeAdbInterface%\t= USB_Install, USB\\VID_0955&PID_CF00&MI_01\r\n\r\n"
		note8_str := ";NVIDIA Tegra\r\n%SingleAdbInterface%\t= USB_Install, USB\\VID_0955&PID_CF05\r\n%CompositeAdbInterface%\t= USB_Install, USB\\VID_0955&PID_CF05&MI_01\r\n\r\n"
		shield_str := ";NVIDIA Tegra\r\n%SingleAdbInterface%\t= USB_Install, USB\\VID_0955&PID_B400&MI_00\r\n%CompositeAdbInterface%\t= USB_Install, USB\\VID_0955&PID_B400&MI_01\r\n\r\n"
		shield1_str := ";NVIDIA Tegra\r\n%SingleAdbInterface%\t= USB_Install, USB\\VID_0955&PID_B400\r\n%CompositeAdbInterface%\t= USB_Install, USB\\VID_0955&PID_B401&MI_01\r\n\r\n"

		newstr := "[Google.NTx86]\r\n\r\n" + tegra_str + note_str + note8_str + shield_str + shield1_str
		myutils.Substitude(fname, oldstr, newstr)
		//os.Rename(filepath.Join(install_location, ".tmp_file"), fname)
		myutils.CopyFile(filepath.Join(install_location, ".tmp_file"), fname)

		//substitude the Google.NTamd64
		oldstr = "[Google.NTamd64]"
		newstr = "[Google.NTamd64]\r\n\r\n" + tegra_str + note_str + note8_str + shield_str + shield1_str
		myutils.Substitude(fname, oldstr, newstr)
		//os.Rename(filepath.Join(install_location, ".tmp_file"), fname)
		myutils.CopyFile(filepath.Join(install_location, ".tmp_file"), fname)

		//installing
		pnpu_path := filepath.Join("C:\\Windows", "system32", "pnputil.exe")
		_, e := exec.Command(pnpu_path, "-a", fname).Output()
		//fmt.Println("meet error is", pnpu_path, "-i -a", fname, string(out), e)

		myutils.CheckError(e)
	}
	//win8 action
	myutils.Decompress(args[2], location)
	os.RemoveAll(filepath.Join(location, "Win8"))
	myutils.CopyFile(filepath.Join(location, "usb_driver"), filepath.Join(location, "Win8"))
	fmt.Println(location)
}

func (u *USB) Uninstall(args ...string) {
	u.ComponentId = args[0]
	u.Version = args[1]
	u.InstallLocation = filepath.Join(args[2], "Drivers", "Devkits")
	e := os.RemoveAll(u.InstallLocation)
	myutils.CheckError(e)
	os.Remove(filepath.Join(args[2], "Drivers"))
	fmt.Println(u.InstallLocation)
}
