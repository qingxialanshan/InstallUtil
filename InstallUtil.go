package main

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) == 3 {
		if os.Args[1] == "run" {
			//run the Poller.exe
			poller_dir := os.Args[2]
			e := exec.Command(poller_dir).Start()
			fmt.Println("Done", e)
			os.Exit(0)

		} else if os.Args[1] == "rusb" {
			var pnpu string
			if runtime.GOARCH == "amd64" {
				pnpu = "C:\\windows\\system32\\pnputil.exe"
			} else {
				pnpu = "C:\\windows\\system\\pnputil.exe"
			}
			fmt.Println(pnpu)
			_, e := exec.Command(pnpu, "-a", os.Args[2]).Output()
			if e != nil {
				fmt.Println(e)
				os.Exit(2)
			}
			os.Exit(0)
		}
	}
	if len(os.Args) < 4 {
		fmt.Println("Please supply a command and component id!\nUsage: InstallUtil [Install] [compName] [version] [exefile] [install_location]")
		fmt.Println("      		   [uninstall] [compName] [version] [install_location]")
		fmt.Println("      		   [query] [compName] [location]")
		fmt.Println("      		   [deploy] [location] [version]")
		fmt.Println("      		   [preInstall] [compName] [version] [task] [target] {-c}")
		fmt.Println("      		   [reg_uninstaller] [-i|-u] [installdir]")
		os.Exit(-1)
	} else {
		var pkgname string
		if os.Args[1] == "deploy" {
			c := &myutils.CompileCommand{}
			c.Run(os.Args[1], os.Args[2], os.Args[3])
			os.Exit(0)
		} else if os.Args[1] == "env_add" {

			envadd := &myutils.EnvaddCommand{}
			envadd.Run(os.Args[2], os.Args[3])
			os.Exit(0)
		} else if os.Args[1] == "reg_uninstaller" {
			//add uninstaller to windows registery
			// [Key Name]
			// HKEY_LOCAL_MACHINE\SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall\NVIDIA Tegra Android Development Pack
			//[Name]                    [Value]
			//DisplayIcon               C:\NVPACK\_installer\TegraDeveloperKit.ico
			//NoModify                  0
			//ModifyPath                C:\NVPACK\Chooser.exe
			//UninstallString           C:\NVPACK\tadp_uninstaller.exe
			regName := `HKLM\Software\Wow6432Node\microsoft\windows\CurrentVersion\Uninstall\NVIDIA Tegra Android Development Pack`
			cmdExe := "C:\\Windows\\system32\\cmd.exe"
			option := os.Args[2]
			if option == "-i" {
				path_chooser := filepath.Join(os.Args[3], "Chooser.exe")

				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v NoModify /d 0 /f`).Output()
				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v DisplayIcon /d `+path_chooser+` /f`).Output()
				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v Publisher /d "NVIDIA Corporation" /f`).Output()
				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v DisplayName /d "NVIDIA Tegra Android Development Pack" /f`).Output()
				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v UninstallString /d "`+path_chooser+` -u" /f`).Output()
				exec.Command(cmdExe, "/c", `reg add "`+regName+`" /v ModifyPath /d `+path_chooser+` /f`).Output()
			} else if option == "-u" {
				_, e := exec.Command(cmdExe, "/c", `reg delete `+regName+` /f`).Output()
				myutils.CheckError(e)
			}
		}
		if strings.Contains(os.Args[2], "api") {
			pkgname = "api"
		} else {
			pkgname = os.Args[2]
		}
		if comp, ok := comps[pkgname]; ok {
			if os.Args[1] == "install" {
				if len(os.Args) < 6 {
					fmt.Println("Wrong args input. Usage: InstallUtil install compName version exefile install_location")
					os.Exit(-1)
					return
				}
				comp.Install(os.Args[2:]...)
			} else if os.Args[1] == "uninstall" {
				if len(os.Args) < 5 {
					fmt.Println("Wrong args input. Usage: InstallUtil uninstall compName version install_location")
					os.Exit(-1)
					return
				}
				comps[pkgname].Uninstall(os.Args[2:]...)
			} else if os.Args[1] == "query" {
				if len(os.Args) < 4 {
					fmt.Println("Wrong args input. Usage: InstallUtil uninstall compName version install_location")
					os.Exit(-1)
					return
				}
				if autocomp, q := autoComps[pkgname]; q && os.Args[2] != "api" {
					revision := autocomp.Query(os.Args[2], os.Args[3])
					fmt.Println(revision)
					//fmt.Println(myutils.Get_revision(os.Args[2], os.Args[3]))
				} else {
					//fmt.Println(os.Args[2])
					os.Exit(-1)
				}
			} else if os.Args[1] == "preInstall" {
				if len(os.Args) < 6 {
					fmt.Println("Wrong args input. Usage: InstallUtil uninstall compName version install_location")
					os.Exit(-1)
					return
				}
				if pacomp, q := preActionComps[pkgname]; q {
					os.Exit(pacomp.PreInstall(os.Args[2:]...))
				}
			} else {
				panic(os.Args[1] + " command is not defined!")
			}
		} else {
			fmt.Println(os.Args[2] + " component is not defined!")
			os.Exit(-1)
		}
	}
}
