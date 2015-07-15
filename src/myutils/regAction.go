//+build windows

package myutils

import (
	"fmt"
	//"fmt"
	"strings"
	"syscall"
	"unsafe"
)

func QueryKey(handle syscall.Handle, mylist *[1 << 10]string) (csubkeys uint32) {

	var cvalues uint32

	var achvalue []uint16

	cbName := uint32(255)
	if syscall.RegQueryInfoKey(handle, nil, nil, nil, &csubkeys, nil, nil, &cvalues, nil, nil, nil, nil) != nil {
		return
	}
	//fmt.Println(csubkeys, cvalues, mylist)
	if csubkeys != 0 {
		var buf [1 << 10]uint16
		for i := uint32(0); i < csubkeys; i = i + 1 {
			n := uint32(len(buf))
			//fmt.Println(i, int(csubkeys), achkey)
			if syscall.RegEnumKeyEx(handle, i, &buf[0], &n, nil, nil, nil, nil) != nil {

				return
			}
			ext := syscall.UTF16ToString(buf[:])
			//fmt.Println(i, ext, (*mylist)[i])
			(*mylist)[i] = ext
		}
	}
	//fmt.Println(mylist)
	if cvalues != 0 {
		var j uint32
		for j = 0; j < cvalues; j++ {
			if syscall.RegEnumKeyEx(handle, j, (*uint16)(unsafe.Pointer(&achvalue[0])), &cbName, nil, nil, nil, nil) != nil {
				return
			}

		}
	}
	return csubkeys
}

func RegQuery(regname, installer_name, args string) (result string) {
	var handle, subtestkey syscall.Handle

	var csubkey uint32
	var list [1 << 10]string

	if syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(regname), 0, syscall.KEY_READ, &handle) != nil {
		return
	}
	defer syscall.RegCloseKey(handle)
	csubkey = QueryKey(handle, &list)
	//	fmt.Println("the args is :", args)

	for j := uint32(0); j < csubkey; j++ {

		var buffer, uninstall_buf [syscall.MAX_LONG_PATH]uint16
		n := uint32(len(buffer))
		dwSize := uint32(len(uninstall_buf))
		var display_name string
		reg1 := regname + "\\" + list[j]
		//fmt.Println("reg1 is ", reg1)
		if reg1 == "Software\\Wow6432Node\\Xoreax\\IncrediBuild\\Builder" {
			if syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(reg1), 0, syscall.KEY_READ, &subtestkey) == nil {

				e1 := syscall.RegQueryValueEx(subtestkey, syscall.StringToUTF16Ptr(args), nil, nil, (*byte)(unsafe.Pointer(&buffer[0])), &n)
				if e1 != nil {
					fmt.Println(e1)
				}
				result = syscall.UTF16ToString(buffer[:])
				//				fmt.Println(result, args)
				return
			}
		}
		if syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(reg1), 0, syscall.KEY_READ, &subtestkey) == nil {
			syscall.RegQueryValueEx(subtestkey, syscall.StringToUTF16Ptr("DisplayName"), nil, nil, (*byte)(unsafe.Pointer(&buffer[0])), &n)
			display_name = syscall.UTF16ToString(buffer[:])

			if strings.Contains(display_name, installer_name) {

				if syscall.RegQueryValueEx(subtestkey, syscall.StringToUTF16Ptr(args), nil, nil, (*byte)(unsafe.Pointer(&uninstall_buf[0])), &dwSize) == nil {

					result = syscall.UTF16ToString(uninstall_buf[:])

					if result != "" {
						return
					}

				}
			}
		}
		defer syscall.RegCloseKey(subtestkey)
	}

	return
}

//func main() {
//	path := `Software\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`

//	uninstall_string := RegQuery(path, "Nsight Tegra", "QuietUninstallString")
//	fmt.Println("The uninstall string of Nsight Tegra is :", uninstall_string)
//	rev := RegQuery(path, "Nsight Tegra", "DisplayVersion")
//	fmt.Println("The revision for Nsight Tegra is :", rev)

//	path = `Software\Microsoft\Windows\CurrentVersion\Uninstall`
//	uninstall_string = RegQuery(path, "Tegra System Profiler", "UninstallString")
//	fmt.Println("The uninstall string of Tegra System Profiler is :", uninstall_string)

//}
