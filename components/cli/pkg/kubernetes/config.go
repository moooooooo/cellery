/*
 * Copyright (c) 2019 WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package kubernetes

import (
	"os"
	"os/exec"

	"cellery.io/cellery/components/cli/pkg/osexec"
)

func (kubeCli *CelleryKubeCli) UseContext(context string) error {
	cmd := exec.Command(
		kubectl,
		"config",
		"use-context",
		context,
	)
	displayVerboseOutput(cmd)
	cmd.Stderr = os.Stderr
	return cmd.Run()
	return nil
}

func (kubeCli *CelleryKubeCli) GetContexts() ([]byte, error) {
	cmd := exec.Command(
		kubectl,
		"config",
		"view",
		"-o",
		"json",
	)
	displayVerboseOutput(cmd)
	return osexec.GetCommandOutputFromTextFile(cmd)
}

func (kubeCli *CelleryKubeCli) GetContext() (string, error) {
	cmd := exec.Command(
		kubectl,
		"config",
		"current-context",
	)
	displayVerboseOutput(cmd)
	out, err := osexec.GetCommandOutput(cmd)
	return out, err
}

func (kubeCli *CelleryKubeCli) SetNamespace(namespace string) error {
	cmd := exec.Command(
		kubectl,
		"config",
		"set-context",
		"--current",
		"--namespace",
		namespace,
	)
	displayVerboseOutput(cmd)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
