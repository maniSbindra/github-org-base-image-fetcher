package fileparser

import (
	"bufio"
	"context"
	"net/http"
	"strings"
)

type FileDownloadDetails struct {
	FileName        string
	FileDownloadUrl string
}

type FileParseResult struct {
	FileName                string
	FileDownloadUrl         string
	ContainerImageNameFound bool
	ContainerImageNames     []string
}

func getContainerName(line string) (bool, string) {
	spltLine := strings.Split(strings.TrimSpace(line), " ")
	if spltLine[0] == "FROM" {
		return true, spltLine[1]
	}
	return false, ""
}

func ParseFile(ctx context.Context, fileDownloadUrl string) (bool, []string, error) {
	containerNames := []string{}
	containerImageFound := false
	resp, err := http.Get(fileDownloadUrl)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if exists, containerName := getContainerName(line); exists {
			containerImageFound = true
			containerNames = append(containerNames, containerName)
		}
	}
	return containerImageFound, containerNames, nil
}

func GetFileContainerImageNames(ctx context.Context, workQueue chan *FileDownloadDetails, termChan chan bool, respChan chan<- *FileParseResult) {

	var fileDetails *FileDownloadDetails
	for {
		select {
		case <-termChan: // terminate this worker if message on termination channel
			return
		case fileDetails = <-workQueue: // Scan file for container image names
			contImgFound, contImgNames, err := ParseFile(ctx, fileDetails.FileDownloadUrl)
			if err == nil && contImgFound {
				respChan <- &FileParseResult{FileName: fileDetails.FileName, FileDownloadUrl: fileDetails.FileDownloadUrl, ContainerImageNameFound: true, ContainerImageNames: contImgNames}
			} else {
				respChan <- &FileParseResult{FileName: fileDetails.FileName, ContainerImageNameFound: false}
			}
		}
	}

}
