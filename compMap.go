package main

import (
	"components"
)

type Action interface {
	Install(arg ...string)
	Uninstall(arg ...string)
}

var comps = map[string]Action{
	"eclipse":       &components.Eclipse{},
	"adt":           &components.ADT{},
	"api":           &components.API{},
	"ant":           &components.Ant{},
	"ardbeg":        &components.Ardbeg{},
	"battle":        &components.Battle{},
	"cuda65":        &components.CUDA_65{},
	"cudaandroid65": &components.CUDAAndroid{},
	"cudasamples65": &components.CUDASamples{},
	"cuda70":        &components.CUDA_70{},
	"cudaandroid70": &components.CUDAAndroid{},
	"cudasamples70": &components.CUDASamples{},
	//"cuda":              &components.CUDA{},
	//"cudaandroid":       &components.CUDAAndroid{},
	//"cudasamples":       &components.CUDASamples{},
	"docs":              &components.Docs{},
	"docs-nda":          &components.Docs_NDA{},
	"gameworks_samples": &components.GameWorks{},
	"gradle":            &components.Gradle{},
	"java":              &components.JavaSDK{},
	"ndk":               &components.NDK{},
	"nsight":            &components.NsightTegra{},
	"ocl_samples":       &components.OclSamples{},
	"opencv":            &components.OpenCV{},
	"perfhudes":         &components.PerfHUDES{},
	"perfkit":           &components.Perfkit{},
	"physx":             &components.PhySX{},
	"rebel":             &components.Rebel{},
	"tdksample":         &components.Samples{},
	"sdkbase":           &components.SDK{},
	"shield":            &components.Shield{},
	"supportlib":        &components.Supportlib{},
	"supportlibrep":     &components.SupportlibRep{},
	"platformtools":     &components.PlatformTools{},
	"buildtools":        &components.BuildTools{},
	"tegraprofiler":     &components.Quadd{},
	"usbdriver":         &components.USB{},
	"usbnv":             &components.USB_NV{},
	"allowcompile":      &components.AllowCompile{},
	"drivertools":       &components.DriverTools{},
	"os":                &components.FlashOS{},
	"incredibuild":      &components.IncrediBuild{},
	"visualstudio":      &components.VisualStudio{},
}

type QAction interface {
	Query(args ...string) string
}

var autoComps = map[string]QAction{
	"sdkbase":       &components.SDK{},
	"supportlib":    &components.Supportlib{},
	"supportlibrep": &components.SupportlibRep{},
	"platformtools": &components.PlatformTools{},
	"buildtools":    &components.BuildTools{},
	"api":           &components.API{},
	"tegraprofiler": &components.Quadd{},
	"battle":        &components.Battle{},
	"nsight":        &components.NsightTegra{},
	"ndk":           &components.NDK{},
	"incredibuild":  &components.IncrediBuild{},
	"visualstudio":  &components.VisualStudio{},
}

type PreAction interface {
	PreInstall(args ...string) int
}

var preActionComps = map[string]PreAction{
	"platformtools": &components.PlatformTools{},
	"battle":        &components.Battle{},
	"tegraprofiler": &components.Quadd{},
}
