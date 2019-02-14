/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (
	"fmt"
	"os/exec"
	"sync"
)

type mounter interface {
	AddMount(string, string, string) error
	RemoveMount(string, string, string) error
}

type fstabMounter struct {
	fileMutex *sync.Mutex
}

func newFstabMounter() *fstabMounter {
	return &fstabMounter{
		fileMutex: &sync.Mutex{},
	}
}

func (e *fstabMounter) AddMount(device, path, mountType string) error {
	block := e.createMountBlock(device, path, mountType)

	cmd := exec.Command("mount", "-t", mountType, device, path)
	log(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mount failed with error: %v, output: %s", err, out)
	}

	// Add the export block to the config file
	if err := addToFile(e.fileMutex, "/etc/fstab", block); err != nil {
		return fmt.Errorf("error adding mount block %s to /etc/fstab: %v", block, err)
	}

	return nil
}

func (e *fstabMounter) RemoveMount(device, path, mountType string) error {
	cmd := exec.Command("umount", device)
	log(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("umount failed with error: %v, output: %s", err, out)
	}

	block := e.createMountBlock(device, path, mountType)

	return removeFromFile(e.fileMutex, "/etc/fstab", block)
}

// CreateBlock creates the text block to add to the /etc/exports file.
func (e *fstabMounter) createMountBlock(device, path, mountType string) string {
	return "\n" + device + "\t" + path + "\t" + mountType + "\tdefaults\t0 2\n"
}
