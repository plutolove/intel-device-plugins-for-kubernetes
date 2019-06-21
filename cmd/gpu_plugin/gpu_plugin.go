// Copyright 2017 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"regexp"
	"time"

	//"regexp"
	//"time"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	dpapi "github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
	devfsDriDirectory = "/dev"
	gpuDeviceRE       = `^cambricon_c10Dev[0-9]+$`

	// Device plugin settings.
	namespace  = "nvidia.com"
	deviceType = "gpu"
)

type devicePlugin struct {
	devfsDir string

	sharedDevNum int

	gpuDeviceReg     *regexp.Regexp
}

func newDevicePlugin(devfsDir string, sharedDevNum int) *devicePlugin {
	return &devicePlugin{
		devfsDir:         devfsDir,
		sharedDevNum:     sharedDevNum,
		gpuDeviceReg:     regexp.MustCompile(gpuDeviceRE),
	}
}

func (dp *devicePlugin) Scan(notifier dpapi.Notifier) error {
	for {
		devTree, err := dp.scan()
		if err != nil {
			return err
		}

		notifier.Notify(devTree)

		time.Sleep(5 * time.Second)
	}
}

func (dp *devicePlugin) scan() (dpapi.DeviceTree, error) {
	files, err := ioutil.ReadDir(dp.devfsDir)
	if  err != nil {
		fmt.Println("something error")
		return nil, err
	}
	i := 0
	devTree := dpapi.NewDeviceTree()
	for _, f := range files {
		var nodes []pluginapi.DeviceSpec
		if dp.gpuDeviceReg.MatchString(f.Name()) {
			devPath := path.Join(dp.devfsDir, f.Name())
			fmt.Printf("%s\n", devPath)
			nodes = append(nodes, pluginapi.DeviceSpec{
				HostPath:      devPath,
				ContainerPath: devPath,
				Permissions:   "rw",
			})
			devID := fmt.Sprintf("%s-%d", f.Name(), i)
			devTree.AddDevice(deviceType, devID, dpapi.DeviceInfo{
				State: pluginapi.Healthy,
				Nodes: nodes,
			})
			i += 1
		}
	}
	return devTree, nil
}

func main() {

	fmt.Println("GPU device plugin started")

	plugin := newDevicePlugin(devfsDriDirectory, 100)
	manager := dpapi.NewManager(namespace, plugin)
	manager.Run()
}
