package myutils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

var gloabal_path string

func find(dir, pattern string) string {
	//fmt.Println("0000", dir)

	if Exists(dir) == false {
		return ""
	}
	d, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println("Reading " + dir)

	for _, file := range files {
		new_file := filepath.Join(dir, file.Name())

		if strings.Contains(file.Name(), pattern) {
			//fmt.Println("^^^", file.Name(), dir, pattern)
			gloabal_path = gloabal_path + "," + new_file
		}
		if file.IsDir() {

			find(new_file, pattern)
		}
		//if strings.Contains(file.Name,"-debug.apk")
	}
	return gloabal_path
}

//check os architecture is 64bit or 32bit. If is 64bit then return 1,else return 0
func Get_os() int {
	var x uintptr
	os := unsafe.Sizeof(x)
	if os == 8 {
		return 1
	} else {
		return 0
	}
}

var set_env_command string

func set_env_win(installdir, version string) {

	var wow string
	wow = ""
	if Get_os() == 1 {
		wow = "\\WOW6432Node"
	} else {
		wow = ""
	}
	//set environment for pentak
	var jdk_version string
	if version == "3.0r4" {
		jdk_version = "1.7.0_45"
	} else {
		jdk_version = "1.7.0_71"
	}
	env_map := map[string]string{"sdkRoot": "android-sdk-windows",
		"ndkRoot":    "android-ndk-r10e",
		"antRoot":    "apache-ant-1.8.2",
		"jdkRoot":    "jdk" + jdk_version,
		"gradleRoot": "gradle-2.2.1"}

	Hkey := "HKEY_LOCAL_MACHINE\\SOFTWARE" + wow + "\\NVIDIA Corporation\\Nsight Tegra"
	for key, value := range env_map {
		fullpath := filepath.Join(installdir, value)
		//fmt.Println(fullpath)
		_, e := exec.Command("reg", "add", Hkey, "/v", key, "/t", "REG_SZ", "/f", "/d", fullpath).Output()
		if e != nil {

			fmt.Println(e)
			os.Exit(2)
		}
	}

}
func set_env_unix(installdir, version string) {
	var ndk, sdk string
	ndk_version := Get_NDK_Revision(installdir)
	if strings.Contains(ndk_version, "r10e") {
		ndk_version = "r10e"
	}
	ndk = "android-ndk-" + ndk_version
	//fmt.Println(ndk)

	if Global_OS == "linux" {
		sdk = "/android-sdk-linux"
	} else {
		sdk = "/android-sdk-macosx"
	}
	new_path := os.Getenv("PATH") + ":" + installdir + sdk + ":" + installdir + "/" + ndk + ":" + installdir + "/apache-ant-1.8.2/bin:" + installdir + sdk + "/platform-tools:" + installdir + sdk + "/tools"

	os.Setenv("NDK_ROOT", installdir+"/"+ndk)
	os.Setenv("NDKROOT", installdir+"/"+ndk)
	os.Setenv("NVPACK_ROOT", installdir)
	os.Setenv("ANT_HOME", installdir+"/apache-ant-1.8.2")
	os.Setenv("ANDROID_HOME", installdir+sdk)
	os.Setenv("NVPACK_NDK_VERSION", ndk)
	os.Setenv("PATH", new_path)
	os.Setenv("CUDA_TOOLKIT_ROOT", installdir+"/cuda-7.0")
	//fmt.Println(os.Getenv("PATH"))
}

func set_env(installdir, version string) {
	if runtime.GOOS == "windows" {
		set_env_win(installdir, version)
	} else {
		set_env_unix(installdir, version)
	}
}

//detect the device is T124 or T210
func detect_devices(installdir string) bool {
	//if device is arch64 return 1, else return 0
	//device_arch, _ := exec.Command("bash", "-c", `. ~/.bashrc && adb shell getprop ro.product.cpu.abi`).Output()
	adb_path := filepath.Join(installdir, "android-sdk-linux", "platform-tools", "adb")
	device_arch, _ := exec.Command("bash", "-c", adb_path+" shell getprop ro.product.cpu.abi").Output()
	//fmt.Printf("Device is: %s", device_arch)

	is_64 := false
	if strings.Contains(string(device_arch), "arm64-v8a") {
		is_64 = true
	} else if strings.Contains(string(device_arch), "armeabi-v7a") {
		is_64 = false
	} else {
		fmt.Println("cannot detect the device")
	}
	return is_64
}

// deploy and install the samples' apk to device
func deploy(installdir string) {

	out := find(filepath.Join(installdir, "Samples"), "-debug.apk")
	apks := strings.Split(out, ",")
	adb_path := filepath.Join(installdir, "android-sdk-"+Global_OS, "platform-tools", "adb")
	var deploy_command, console string
	if Global_OS == "windows" {
		console = filepath.Join(installdir, "console.exe")
		deploy_command = "cmd /c echo Deploy Samples;cmd /c echo Waiting for device...;"
	} else {
		if Global_OS == "macosx" {
			console = filepath.Join(installdir, "console.app", "Contents", "MacOS", "console")
		} else {
			console = filepath.Join(installdir, "console")
		}

		deploy_command = "echo Deploy Samples;echo Waiting for device...;"
	}

	for i := 0; i < len(apks); i++ {

		if apks[i] == "" {
			continue
		}
		deploy_command = deploy_command + adb_path + " install -r " + apks[i] + ";"
	}

	deploy_log := filepath.Join(installdir, "_installer/deploy.log")

	var cuda_samples string
	gloabal_path = ""
	if detect_devices(installdir) {
		cuda_samples = find(filepath.Join(installdir, "CUDA_Samples", "7.0"), "-debug.apk")

	} else {
		cuda_samples = find(filepath.Join(installdir, "CUDA_Samples", "6.5"), "-debug.apk")

	}

	cudaapks := strings.Split(cuda_samples, ",")

	for i := 0; i < len(cudaapks); i++ {
		//deploying cuda samples
		if cudaapks[i] == "" {
			continue
		}
		deploy_command = deploy_command + adb_path + " install -r " + cudaapks[i] + ";"
	}

	e := exec.Command(console, "-t", "Deploy Samples", "-l", deploy_log, deploy_command).Run()
	if e != nil {
		fmt.Println("Deploy Samples Failed.")
		os.Exit(2)
	}

	exec.Command(adb_path, "kill-server")
	if Global_OS == "windows" {
		Redirector(exec.Command("taskkill", "/F", "/IM", "adb.exe"))
	} else {
		cmd := exec.Command("/bin/sh", "-c", `ps -e|grep "adb"|awk '{print $1}'`)
		pid, err := cmd.Output()
		if err != nil {
			fmt.Println("failed")
		}

		pids := strings.Fields(string(pid[:]))

		if len(pids) > 0 {
			for i := 0; i < len(pids); i++ {
				exec.Command("kill", pids[i]).Start()
				//fmt.Println("Success")
			}
		}
	}
}

// deploy and install the samples' apk to device
func deploy_v3a(installdir string) {

	skip_apks := []string{"BindlessApp", "BlendedAA", "CascadedShadowMapping", "ConservativeRaster", "FeedbackParticlesApp", "MultiDrawIndirect", "NormalBlendedDecal", "PathRenderingBasic", "ShapedTextES", "SystemPerf", "TerrainTessellation", "Tiger3DES", "TigerWarpES", "WeightedBlendedOIT", "ComputeWaterSimulation"}
	out := find(filepath.Join(installdir, "Samples"), "-debug.apk")
	apks := strings.Split(out, ",")
	adb_path := filepath.Join(installdir, "android-sdk-"+Global_OS, "platform-tools", "adb")
	var deploy_command, console string
	if Global_OS == "windows" {
		console = filepath.Join(installdir, "console.exe")
		deploy_command = "cmd /c echo Deploy Samples;cmd /c echo Waiting for device...;"
	} else {
		if Global_OS == "macosx" {
			console = filepath.Join(installdir, "console.app", "Contents", "MacOS", "console")
		} else {
			console = filepath.Join(installdir, "console")
		}

		deploy_command = "echo Deploy Samples;echo Waiting for device...;"
	}

	for i := 0; i < len(apks); i++ {
		skip := 0
		for _, skip_apk := range skip_apks {
			if strings.Contains(apks[i], skip_apk) {
				skip = 1
			}
		}

		if apks[i] == "" || skip == 1 {
			//fmt.Println(apks[i])
			continue
		}

		deploy_command = deploy_command + adb_path + " install -r " + apks[i] + ";"
	}

	deploy_log := filepath.Join(installdir, "_installer/deploy.log")

	var cuda_samples string
	gloabal_path = ""

	cuda_samples = find(filepath.Join(installdir, "CUDA_Samples", "7.0"), "-debug-32.apk")

	cudaapks := strings.Split(cuda_samples, ",")

	for i := 0; i < len(cudaapks); i++ {
		//deploying cuda samples
		if cudaapks[i] == "" {
			continue
		}
		deploy_command = deploy_command + adb_path + " install -r " + cudaapks[i] + ";"
	}

	//fmt.Println(deploy_command)
	e := exec.Command(console, "-t", "Deploy Samples", "-l", deploy_log, deploy_command).Run()
	if e != nil {
		fmt.Println("Deploy Samples Failed.")
		os.Exit(2)
	}

	exec.Command(adb_path, "kill-server")
	if Global_OS == "windows" {
		Redirector(exec.Command("taskkill", "/F", "/IM", "adb.exe"))
	} else {
		cmd := exec.Command("/bin/sh", "-c", `ps -e|grep "adb"|awk '{print $1}'`)
		pid, err := cmd.Output()
		if err != nil {
			fmt.Println("failed")
		}

		pids := strings.Fields(string(pid[:]))

		if len(pids) > 0 {
			for i := 0; i < len(pids); i++ {
				exec.Command("kill", pids[i]).Start()
				//fmt.Println("Success")
			}
		}
	}
}

// check whether file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Redirector(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// compile samples
func compile_win(installdir string) {

	tdksamples := installdir + "/Samples/TDK_Samples/android_samples.sln"
	gwsamples := installdir + "/Samples/GameWorks_Samples/samples/build/vs2010android/AllSamples.sln"

	var wow, vs_str, devenv_path string
	wow = ""
	if Get_os() == 1 {
		wow = "\\WOW6432Node"
	} else {
		wow = ""
	}
	reg := "HKEY_LOCAL_MACHINE\\Software" + wow + "\\Microsoft\\VisualStudio\\"

	vs := []float64{12.0, 11.0, 10.0}

	for i := 0; i < 3; i++ {

		vspath := reg + strconv.FormatFloat(vs[i], 'f', 1, 64)

		vsdir, _ := exec.Command("reg", "query", vspath, "/v", "InstallDir").Output()
		vsstring := strings.Split(string(vsdir), "\r\n")

		for str := range vsstring {
			if strings.Contains(vsstring[str], "InstallDir") {
				vs_str = vsstring[str]
				break
			}
		}
		vs_path := strings.Split(vs_str, "    ")
		//fmt.Println("0:", vs_path[0], "1:", vs_path[1], "2:", vs_path[2], "3:", vs_path[3])
		//vsdir = string(vsdir)
		if len(vs_path) < 4 {
			continue
		}
		verifypath := vs_path[3] + "devenv.com"

		if Exists(verifypath) {
			devenv_path = verifypath
			os.Chdir(string(vsdir))
			break
		}
	}

	console := filepath.Join(installdir, "console.exe")
	compile_log := filepath.Join(installdir, "_installer/compile.log")
	compile_command := ""
	if Exists(gwsamples) {
		compile_command = "cmd /c echo ##############Compiling GameWorks_Samples##############;" + devenv_path + " /rebuild debug " + gwsamples
	}
	if Exists(tdksamples) {
		compile_command = compile_command + ";cmd /c echo '############## Compiling TDK_Samples ##############';cmd /c cd " + tdksamples + ";" + devenv_path + " /rebuild debug " + tdksamples
	}
	_, e := exec.Command(console, "-t", "Compile Samples", "-l", compile_log, compile_command).Output()
	if e != nil {
		fmt.Println("Compile Samples failed. Please check the " + installdir + "/_installer/compile.log. You can send it to TegraDeveloperPack-Users@nvidia.com if you can not resolve it.")
		os.Exit(2)
	}
}

func compile_unix(installdir, version string) {
	tdksamples := installdir + "/Samples/TDK_Samples"
	cudasamples65 := installdir + "/CUDA_Samples/6.5"
	cudasamples70 := installdir + "/CUDA_Samples/7.0"

	oclsamples := installdir + "/Samples/OpenCL_Samples/oclConvolutionSeparable"
	gwsamples := installdir + "/Samples/GameWorks_Samples/samples/build/makeandroid"
	compile_command := ""

	var console, bash_script string
	if Global_OS == "macosx" {
		console = filepath.Join(installdir, "console.app", "Contents", "MacOS", "console")
		bash_script = "~/.bash_profile"
	} else {
		console = filepath.Join(installdir, "console")
		bash_script = "~/.bashrc"
	}

	compile_log := filepath.Join(installdir, "_installer/compile.log")
	if Exists(tdksamples) {
		compile_command = compile_command + "echo '########## Compiling TDK_Samples ##########';bash -c '. " + bash_script + " && cd " + tdksamples
		//compile_command = compile_command + "for i in `find . -iname *.vcxproj|grep -v libs`;do tdk_applist=`dirname $i`;bash -c 'cd $tdk_applist';bash -c 'android update project -p . --target android-15' && ndk-build -C jni clean && ndk-build -C jni; ant debug;bash -c 'cd "+tdksamples+";done"
		os.Chdir(tdksamples)

		out, _ := exec.Command("bash", "-c", `find . -iname "*.vcxproj"|grep -v "libs"`).Output()
		tdk_applist := strings.Fields(string(out[:]))
		for i := 0; i < len(tdk_applist); i++ {
			tdk_appDir := filepath.Dir(tdk_applist[i])
			compile_command = compile_command + "&& cd " + tdk_appDir + " && android update project -p . --target android-15 && ndk-build -C jni clean && ndk-build -C jni && ant clean && ant debug && cd " + tdksamples

		}
		compile_command = compile_command + "';"

	}

	if Exists(cudasamples70) {
		cuda_path := filepath.Join(installdir, "cuda-7.0")
		compile_command = compile_command + "echo '######### Compiling CUDA Samples 7.0 ##########';bash -c '. ~/.bashrc && export CUDA_TOOLKIT_ROOT=" + cuda_path + "&& cd " + cudasamples70
		os.Chdir(cudasamples70)
		out, _ := exec.Command("bash", "-c", `ls -d ./*`).Output()
		cuda_applist := strings.Fields(string(out[:]))
		//fmt.Println(cuda_applist)
		for i := 0; i < len(cuda_applist); i++ {
			cuda_appDir := cuda_applist[i]
			if cuda_appDir == "./README.txt" {
				continue
			}
			//build for armv7l
			compile_command = compile_command + "&& cd " + cuda_appDir + "&& cd cuda && make TARGET_ARCH=armv7l TARGET_OS=android SMS=32 clean build&& cd ../ && android update project -p . --target android-15 && ndk-build clean && ndk-build -C jni && ant debug && cp bin/NativeActivity-debug.apk bin/NativeActivity-debug-32.apk && cd " + cudasamples70

			if version == "4.0" {
				//build for aarch64
				compile_command = compile_command + "&& cd " + cuda_appDir + "&&cd cuda && make TARGET_ARCH=aarch64 TARGET_OS=android SMS=53 clean build && cd ../ && android update project -p . --target android-21 &&ndk-build -C jni NV_TARGET_ARCH=aarch64 clean && ndk-build -C jni NV_TARGET_ARCH=aarch64 && ant debug && cd " + cudasamples70
			}
		}
		compile_command = compile_command + "';"
		//fmt.Println(compile_command)
	}

	if Exists(cudasamples65) {
		cuda_path := filepath.Join(installdir, "cuda-6.5")
		compile_command = compile_command + "echo '######### Compiling CUDA Samples 6.5 ##########';bash -c '. ~/.bashrc &&export CUDA_TOOLKIT_ROOT=" + cuda_path + " && cd " + cudasamples65
		os.Chdir(cudasamples65)
		out, _ := exec.Command("bash", "-c", `ls -d ./*`).Output()
		cuda_applist := strings.Fields(string(out[:]))
		//fmt.Println(cuda_applist)
		for i := 0; i < len(cuda_applist); i++ {
			cuda_appDir := cuda_applist[i]
			if cuda_appDir == "./README.txt" {
				continue
			}
			//build for armv7l
			compile_command = compile_command + "&& cd " + cuda_appDir + "&& cd cuda && make -j 4&& cd ../ && android update project -p . --target android-15 && ndk-build -C jni clean && ndk-build -C jni && ant debug && cd " + cudasamples65

		}
		compile_command = compile_command + "';"
	}

	if Exists(oclsamples) {
		compile_command = compile_command + "echo '########## Compiling OCL_Samples ##########';bash -c '. ~/.bashrc && cd " + oclsamples
		compile_command = compile_command + "&& cd opencl && make && android update project -p . --target android-15 && ndk-build -C jni && ndk-build -C jni&& ant debug';"
	}

	if Exists(gwsamples) {

		compile_command = compile_command + "echo '########## Compiling GameWorks_Samples ##########';bash -c '. ~/.bashrc &&cd " + gwsamples + " &&make clean&&make -j 4';"
	}

	//fmt.Println("compile command is:", compile_command)

	_, e := exec.Command(console, "-t", "Compile Samples", "-l", compile_log, compile_command).Output()
	if e != nil {
		log := Readfile(filepath.Join(installdir, "_installer", "compile.log"))
		if strings.Contains(log, "JAVA_HOME does not point to the JDK") {
			fmt.Println("Compile Sample failed due to your JAVA_HOME does not point to the JDK. Please install the jdk first to support this action")
		} else {
			fmt.Println("Compile Samples failed. Please check the " + installdir + "/_installer/compile.log. You can send it to TegraDeveloperPack-Users@nvidia.com if you can not resolve it.")
		}

		os.Exit(2)
	}

}

type CompileCommand struct {
}

func (*CompileCommand) Run(args ...string) {

	if len(args) < 2 {
		fmt.Println("Please supply workdirectory and deploy/compile!")
		os.Exit(2)
	}
	installdir := args[1]
	root := installdir + "/Samples"
	os.Chdir(root)

	if _, err := os.Stat(filepath.Join(args[1], "_installer")); err != nil {
		os.MkdirAll(filepath.Join(args[1], "_installer"), 0777)
	}

	//fmt.Println(os.Getenv("Path"))
	set_env(installdir, args[2])
	if args[0] == "deploy" {
		if args[2] == "1.0" {
			deploy_v3a(installdir)
		} else {
			deploy(installdir)
		}

	} else if args[0] == "compile" {

		if runtime.GOOS == "windows" {
			compile_win(args[1])
		} else {
			compile_unix(args[1], args[2])
		}
	} else {
		fmt.Println("Erro parameter!")
	}

}
