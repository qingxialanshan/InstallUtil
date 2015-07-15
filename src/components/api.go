package components

import (
	"fmt"
	"myutils"
	"os"
	"path/filepath"
	"strings"
)

type API struct {
	Component
}

var api_list = map[string]string{

	"api8":  "android-2.2_r03-linux",
	"api9":  "android-2.3.1_r02-linux",
	"api10": "android-2.3.3_r02-linux",
	"api11": "android-3.0_r02-linux",
	"api12": "android-3.1_r03-linux",
	"api13": "android-3.2_r01-linux",
	"api14": "android-4.0.2",
	"api15": "android-4.0.4",
	"api16": "android-4.1.2",
	"api17": "android-4.2.2",
	"api18": "android-4.3.1",
	"api19": "android-4.4.2",
	"api20": "android-4.4W",
	"api21": "android-5.0.1",
	"api22": "android-5.1.1",
}

func (api *API) Install(args ...string) {
	api.ComponentId = args[0]
	api.Version = args[1]
	if api.ComponentId == "api21" || api.ComponentId == "api22" {
		api_list[api.ComponentId] = "android-" + api.Version
	}
	location := filepath.Join(args[3], api_list[api.ComponentId])
	new_location := filepath.Join(args[3], strings.Replace(args[0], "api", "android-", 1))
	myutils.Decompress(args[2], args[3])

	os.RemoveAll(new_location)
	//err := os.Rename(location, new_location)
	_, err := myutils.CopyFile(location, new_location)
	myutils.CheckError(err)
	fmt.Println(new_location)

}

func (api *API) Uninstall(args ...string) {

	if len(args) < 3 {
		fmt.Println("wrong args input")
		return
	}
	api.ComponentId = args[0]
	api.Version = args[1]
	api.InstallLocation = filepath.Join(args[2], strings.Replace(args[0], "api", "android-", 1))
	fmt.Println(api.InstallLocation)
	os.RemoveAll(api.InstallLocation)

}

func (api *API) Query(args ...string) string {
	api.ComponentId = args[0]
	api.InstallLocation = args[1]
	revision := myutils.Get_By_Tags(api.ComponentId, api.InstallLocation, "Platform.Version")
	return revision
}
