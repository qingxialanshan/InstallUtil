package components

import (
	"fmt"
	"myutils"
	"os"
	"os/exec"
	"path/filepath"
)

type CUDA_70 struct {
	Component
}

func (c *CUDA_70) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	c.ComponentId = args[0]
	c.Version = "7.0"

	exec.Command("bash", "-c", "cd "+args[3]+"&&tar xvf "+args[2]).Run()
	//myutils.Decompress(args[2], args[3])

	//Add enviroment variables
	//1. add cuda/bin to PATH

	cuda_path := filepath.Join(args[3], "cuda-"+c.Version)
	bin_path := filepath.Join(cuda_path, "bin")
	myutils.Set_Environment("PATH", bin_path)

	//2. add CUDA_TOOLKIT_ROOT
	myutils.Set_Environment("CUDA_TOOLKIT_ROOT", cuda_path)
	myutils.Set_Environment("CUDA_TOOLKIT_ROOT_7_0", cuda_path)
	fmt.Println(cuda_path)
}

func (c *CUDA_70) Uninstall(args ...string) {
	//fmt.Println("uninstalling cuda")
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	c.ComponentId = args[0]
	c.Version = "7.0"
	c.InstallLocation = filepath.Join(args[2], "cuda-7.0")

	os.RemoveAll(c.InstallLocation)
	myutils.Delete_Environment("CUDA_TOOLKIT_ROOT", c.InstallLocation)
	myutils.Delete_Environment("CUDA_TOOLKIT_ROOT_7_0", c.InstallLocation)
	myutils.Delete_Environment("PATH", filepath.Join(c.InstallLocation, "bin"))
}

type CUDA_65 struct {
	Component
}

func (c65 *CUDA_65) Install(args ...string) {
	if len(args) < 4 {
		fmt.Println("wrong args input")
		return
	}
	c65.ComponentId = args[0]
	c65.Version = "6.5"
	//myutils.Decompress(args[2], args[3])

	exec.Command("bash", "-c", "cd "+args[3]+"&&tar xvf "+args[2]).Run()
	//Add enviroment variables
	//1. add cuda/bin to PATH

	cuda_path := filepath.Join(args[3], "cuda-6.5")
	//bin_path := filepath.Join(cuda_path, "bin")
	//myutils.Set_Environment("PATH", bin_path)

	//2. add CUDA_TOOLKIT_ROOT
	myutils.Set_Environment("CUDA_TOOLKIT_ROOT_6_5", cuda_path)
	fmt.Println(cuda_path)
}

func (c65 *CUDA_65) Uninstall(args ...string) {
	fmt.Println("uninstalling cuda")
	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	c65.ComponentId = args[0]
	c65.Version = "6.5"
	c65.InstallLocation = filepath.Join(args[2], "cuda-"+c65.Version)

	os.RemoveAll(c65.InstallLocation)
	myutils.Delete_Environment("CUDA_TOOLKIT_ROOT_6_5", c65.InstallLocation)
	//myutils.Delete_Environment("PATH", filepath.Join(c.InstallLocation, "bin"))
}
