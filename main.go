package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/manisbindra/get-github-vulnerable-container-images/fileparser"
	"github.com/manisbindra/get-github-vulnerable-container-images/ghclient"
)

type FileContainerImages map[string][]string

type ContainerImageFiles map[string][]string

func main() {

	var fileContainerImages FileContainerImages
	var containerImageFiles ContainerImageFiles
	var filter *ghclient.MatchingFileFilter
	var files []ghclient.FileDetails
	var err error

	fileContainerImages = make(map[string][]string, 0)
	containerImageFiles = make(map[string][]string, 0)

	// ***** Set flags to filter files *****
	ghOrg := flag.String("org", "", "Github organization name to scan for files")
	ghFileName := flag.String("fileName", "Dockerfile", "File name to filter")
	ghRepoName := flag.String("repoName", "", "Repository name to filter")
	ghToken := flag.String("ghToken", "", "Github token to access private repositories")
	generateOutputFile := flag.Bool("generateOutputFile", false, "Generate output file")
	outputFile := flag.String("outputFile", "./containerImageFiles.json", "Output file pull path and  name")
	noOfWorkers := flag.Int("noOfWorkers", 1, "Number of workers to be used for processing files")

	flag.Parse()

	if *ghOrg == "" {
		log.Fatal("Github organization name is required to scan the files")
	}

	filter = &ghclient.MatchingFileFilter{
		Org:      *ghOrg,
		FileName: *ghFileName,
		GHToken:  *ghToken,
	}

	if *ghRepoName != "" {
		filter.Repo = *ghRepoName
	}

	files, err = filter.GetDownloadableFileNames(context.Background())
	if err != nil {
		log.Println("Please verify that the organization / repository names are valid")
		log.Print(err)
		return
	}

	// ***** To load from file instead of making github api calls *****
	// jsonData, _ := ioutil.ReadFile("files.json")
	// err = json.Unmarshal(jsonData, &files)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }

	noOfFiles := len(files)
	log.Printf("fileCount: %d\n", noOfFiles)
	if noOfFiles == 0 {
		log.Printf("No files with name '%s' found \n", *ghFileName)
		return
	}
	// jsonFiles, err := json.MarshalIndent(files, "", " ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// _ = ioutil.WriteFile("files.json", jsonFiles, 0644)

	workQueue := make(chan *fileparser.FileDownloadDetails, noOfFiles)
	termChan := make(chan bool, *noOfWorkers)
	respChan := make(chan *fileparser.FileParseResult, noOfFiles)
	ctx := context.Background()

	if *noOfWorkers > noOfFiles {
		log.Println("Reducing number of workers to number of files......")
		*noOfWorkers = noOfFiles
	}

	for i := 0; i < *noOfWorkers; i++ {
		go fileparser.GetFileContainerImageNames(ctx, workQueue, termChan, respChan)
	}
	log.Println("All workers started...")

	for _, file := range files {
		workQueue <- &fileparser.FileDownloadDetails{
			FileName:        file.FileName,
			FileDownloadUrl: file.FileDownloadURL,
		}
	}
	log.Println("All files added for processing to worker queue...")

	for i := 0; i < noOfFiles; i++ {
		res := <-respChan
		fileContainerImages[res.FileName] = res.ContainerImageNames
		for _, containerImage := range res.ContainerImageNames {
			containerImageFiles[containerImage] = append(containerImageFiles[containerImage], res.FileDownloadUrl)
		}
	}

	log.Println("All files processed...")

	for i := 0; i < *noOfWorkers; i++ {
		termChan <- true
	}

	log.Println("All workers stopped...")

	distinctContainerImages := GetDistinctImages(containerImageFiles)

	if *generateOutputFile {
		log.Println("Writing container images, file mapping to output file...")
		WriteToFile(containerImageFiles, *outputFile)
	}

	log.Println("Printing name of distinct container image names...")
	PrintDistinctImages(distinctContainerImages)

}

func GetDistinctImages(containerImageFiles ContainerImageFiles) []string {
	distinctContainerImages := make([]string, 0, len(containerImageFiles))
	for k := range containerImageFiles {
		distinctContainerImages = append(distinctContainerImages, k)
	}
	return distinctContainerImages
}

func PrintDistinctImages(distinctContainerImages []string) {
	fmt.Println("\n\nDistinct container images:")
	fmt.Println("--------------------------")
	for _, containerImage := range distinctContainerImages {
		fmt.Printf("%s\n", containerImage)
	}
}

func WriteToFile(containerImageFiles ContainerImageFiles, file string) {
	jsonFiles, err := json.MarshalIndent(containerImageFiles, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(file, jsonFiles, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
