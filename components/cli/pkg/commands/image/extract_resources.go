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

package image

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cellery.io/cellery/components/cli/cli"
	"cellery.io/cellery/components/cli/pkg/image"
	"cellery.io/cellery/components/cli/pkg/util"
)

// RunExtractResources extracts the cell image zip file and copies the resources folder to the provided path
func RunExtractResources(cli cli.Cli, cellImage string, outputPath string) error {
	parsedCellImage, err := image.ParseImageTag(cellImage)
	if err != nil {
		return fmt.Errorf("error occurred while parsing cell image, %v", err)
	}

	repoLocation := filepath.Join(cli.FileSystem().Repository(), parsedCellImage.Organization,
		parsedCellImage.ImageName, parsedCellImage.ImageVersion)
	imageLocation := filepath.Join(repoLocation, parsedCellImage.ImageName+cellImageExt)

	// Checking if the image is present in the local repo
	isImagePresent, _ := util.FileExists(imageLocation)
	if !isImagePresent {
		return fmt.Errorf(fmt.Sprintf("failed to extract resources, image %s not found", cellImage))
	}

	// Create temp directory
	currentTIme := time.Now()
	timestamp := currentTIme.Format("27065102350415")
	tempPath := filepath.Join(cli.FileSystem().TempDir(), timestamp)
	err = util.CreateDir(tempPath)
	if err != nil {
		return fmt.Errorf("error while extracting resources from cell image, %v", err)
	}
	defer func() error {
		if err = os.RemoveAll(tempPath); err != nil {
			return fmt.Errorf("error while cleaning up, %v", err)
		}
		return nil
	}()

	if err = util.Unzip(imageLocation, tempPath); err != nil {
		return fmt.Errorf("error extracting zip file, %v", err)
	}

	// Copying the image resources to the provided output directory
	resourcesList, err := util.FindRecursiveInDirectory(filepath.Join(tempPath, src), "resources")

	//resourcesDir, err := filepath.Abs(filepath.Join(tempPath, artifacts, "resources"))
	if err != nil {
		return fmt.Errorf("error occurred while extracting the image resources, %v", err)
	}
	if resourcesList != nil {
		if outputPath == "" {
			outputPath = cli.FileSystem().CurrentDir()
		}
		for _, resourceFile := range resourcesList {
			if err = util.CopyDir(resourceFile, outputPath); err != nil {
				return fmt.Errorf("error occurred while extracting the image resources, %v", err)
			}
		}

		absOutputPath, _ := filepath.Abs(outputPath)
		fmt.Fprintf(cli.Out(), fmt.Sprintf("\nExtracted Resources: %s", util.Bold(absOutputPath)))
		util.PrintSuccessMessage(fmt.Sprintf("Successfully extracted cell image resources: %s", util.Bold(cellImage)))
	} else {
		fmt.Fprintf(cli.Out(), fmt.Sprintf("\nNo resources available in %s\n", util.Bold(cellImage)))
	}
	return nil
}
