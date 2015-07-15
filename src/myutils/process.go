package myutils

import (
	"os/exec"
	"strings"
)

func CheckProcess(exefile string) int {
	//check process run or not on windows
	out, _ := exec.Command("tasklist", "/fi", "IMAGENAME eq "+exefile).Output()
	if strings.Contains(string(out), exefile) {
		return 0
	} else {
		return 1
	}
}

func ManageProcess(task, exefile string) int {
	//process action
	if task == "kill" {
		_, e := exec.Command("taskkill", "/f", "/im", exefile).Output()
		if e == nil {
			return 0
		} else {
			return 1
		}
	} else {
		return 1
	}
}
