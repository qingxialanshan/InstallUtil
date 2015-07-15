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

type USB_NV struct {
	Component
}

func (un *USB_NV) Install(args ...string) {
	//installing the usbdriver for nvidia.
	//1. decompress to ${installdir}/Drivers
	un.ComponentId = args[0]
	un.Version = args[1]
	un.ExecFile = args[2]
	un.InstallLocation = filepath.Join(args[3], "Drivers")

	shield_dir := filepath.Join(un.InstallLocation, "Shield")
	tablet_dir := filepath.Join(un.InstallLocation, "Shield_Tablet")

	if _, e1 := os.Open(shield_dir); e1 == nil {
		//destination directory exist
		myutils.Rmdir(shield_dir)
	}

	if _, e2 := os.Open(tablet_dir); e2 == nil {
		//destination directory exist
		myutils.Rmdir(tablet_dir)
	}

	myutils.Decompress(args[2], un.InstallLocation)

	//2. install driver for shield and shield_tablet
	//pnpu_path := filepath.Join(args[3], "_installer", "pnpuwrapper64.exe")

	shieldfname := filepath.Join(args[3], "Drivers", "Shield", "android_winusb.inf")
	stabletfname := filepath.Join(args[3], "Drivers", "Shield_Tablet", "android_winusb.inf")

	if _, se_open := os.Stat(shieldfname); se_open == nil {
		_, se := exec.Command("pnputil.exe", "-a", shieldfname).Output()
		myutils.CheckError(se)
	}

	if _, ste_open := os.Stat(stabletfname); ste_open == nil {
		_, ste := exec.Command("pnputil.exe", "-a", stabletfname).Output()

		myutils.CheckError(ste)
	}
	fmt.Println(un.InstallLocation)
}

func (un *USB_NV) Uninstall(args ...string) {
	un.ComponentId = args[0]
	un.Version = args[1]
	un.InstallLocation = filepath.Join(args[2], "Drivers")
	if _, err := os.Stat(filepath.Join(un.InstallLocation, "Devkits")); err != nil {
		myutils.Rmdir(un.InstallLocation)
	} else {
		myutils.Rmdir(filepath.Join(un.InstallLocation, "Shield"))
		myutils.Rmdir(filepath.Join(un.InstallLocation, "Shield_Tablet"))
		//myutils.Rmdir(filepath.Join(un.InstallLocation, "How_to_Install_Drivers.pdf"))
	}
	os.Remove(un.InstallLocation)
	fmt.Println(un.InstallLocation)
}
