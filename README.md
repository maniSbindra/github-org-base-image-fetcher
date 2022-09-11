# github-org-base-image-fetcher

Base Image Fetcher is a golang cli utility which can be used to fetch names of base container images used across repositories in a Github organization or a specific repository in the organization. 

# Utility in action

```
go run main.go --org=manisbindra  --fileName=Dockerfile  --generateOutputFile --noOfWorkers=10
```
The above command scans all files in the github organization "manisbindra", which have "Dockerfile" in their file names, and parses those files to retrieve distinct base container images being used. Other command line arguments enable output of json file containing the distinct container images along with all files where the base container image has been used.  Command output is shown below

```sh  
2022/09/11 23:35:47 fileCount: 6
2022/09/11 23:35:47 Reducing number of workers to number of files......
2022/09/11 23:35:47 All workers started...
2022/09/11 23:35:47 All files added for processing to worker queue...
2022/09/11 23:35:48 All files processed...
2022/09/11 23:35:48 All workers stopped...
2022/09/11 23:35:48 Writing container images, file mapping to output file...
2022/09/11 23:35:48 Printing name of distinct container image names...


Distinct container images:
--------------------------
openjdk:18-slim-buster
openjdk:18-jdk-alpine3.14
golang:1.10-stretch
python:2.7-slim
mcr.microsoft.com/dotnet/sdk:5.0
mcr.microsoft.com/dotnet/aspnet:5.0
alpine:latest
```


# Command Line Arguments

* **org**: The Github Organization to scan for files. This argument **is mandatory**
* **fileName**: By default files with **"Dockerfile"** in their name are parsed for base container images. This file name can be modified using this argument
* **repoName**: This can **optionally** be supplied, if aim is to scan only specific repository in organization
* **ghToken**: This needs to be supplied if you wish to scan **private github repositories**. **Please note!** if this parameter is supplied along with **--generateOutputFile** then the output file generated will contain the token as a part of the file download url
* **generateOutputFile**: This argument needs to be specified to enable generation of output json file
* **outputFile**: This is the output file along with the full path. Format of output file is as follows (distinc container image name, along with list of files this image is used in)
  
  ```json
  {
    "alpine:latest": [
        "https://raw.githubusercontent.com/maniSbindra/k8s-delete-validation-webhook/98e33089208081360dbccd384ffdbb815b98ca2a/build/Dockerfile-local-go-build",
        "https://raw.githubusercontent.com/maniSbindra/k8s-delete-validation-webhook/98e33089208081360dbccd384ffdbb815b98ca2a/build/Dockerfile"
    ],
    "golang:1.10-stretch": [
        "https://raw.githubusercontent.com/maniSbindra/k8s-delete-validation-webhook/98e33089208081360dbccd384ffdbb815b98ca2a/build/Dockerfile"
    ]
  }
  ```
* **noOfWorkers**: This is the number of go routines used to parse the files. If an organization with a large number of repositories and Dockerfiles needs to be parsed, setting a high value like 10 is recommended. The default value is 1

## More Sample commands in action

* Full organization scan across public repositories in organization (with no output file)
  
  ```sh
  $ go run main.go --org=manisbindra  --fileName=Dockerfile   --noOfWorkers=10 
    2022/09/12 00:01:08 fileCount: 6
    2022/09/12 00:01:08 Reducing number of workers to number of files......
    2022/09/12 00:01:08 All workers started...
    2022/09/12 00:01:08 All files added for processing to worker queue...
    2022/09/12 00:01:08 All files processed...
    2022/09/12 00:01:08 All workers stopped...
    2022/09/12 00:01:08 Printing name of distinct container image names...


    Distinct container images:
    --------------------------
    openjdk:18-slim-buster
    openjdk:18-jdk-alpine3.14
    mcr.microsoft.com/dotnet/sdk:5.0
    mcr.microsoft.com/dotnet/aspnet:5.0
    python:2.7-slim
    golang:1.10-stretch
    alpine:latest
    ```
    * Scan for single public repository
    * Full scan arcross organization including private repositories
   ```

* Scan of single public repository
  
  ```sh
  $go run main.go --org=manisbindra --repoName=k8s-delete-validation-webhook  --fileName=Dockerfile
    2022/09/12 00:03:37 fileCount: 2
    2022/09/12 00:03:37 All workers started...
    2022/09/12 00:03:37 All files added for processing to worker queue...
    2022/09/12 00:03:37 All files processed...
    2022/09/12 00:03:37 All workers stopped...
    2022/09/12 00:03:37 Printing name of distinct container image names...


    Distinct container images:
    --------------------------
    alpine:latest
    golang:1.10-stretch
  ```

* Specific private repository scan
  
  ```sh
   $ go run main.go --org=manisbindra --repoName=privateRepoName --fileName=Dockerfile  --generateOutputFile --outputFile="./tempfile.json"  --ghToken=$GITHUB_TOKEN
   2022/09/11 23:57:20 fileCount: 2
    2022/09/11 23:57:20 All workers started...
    2022/09/11 23:57:20 All files added for processing to worker queue...
    2022/09/11 23:57:21 All files processed...
    2022/09/11 23:57:21 All workers stopped...
    2022/09/11 23:57:21 Writing container images, file mapping to output file...
    2022/09/11 23:57:21 Printing name of distinct container image names...


    Distinct container images:
    --------------------------
    node:lts
    mcr.microsoft.com/azure-cli
    node:lts-alpine
  ```

* Full Scan across all private and public repositories in organization

  ```sh
  $go run main.go --org=manisbindra  --fileName=Dockerfile  --generateOutputFile --outputFile="./tempfile.json" --noOfWorkers=10 --ghToken=$GITHUB_TOKEN
    2022/09/12 00:05:13 fileCount: 18
    2022/09/12 00:05:13 All workers started...
    2022/09/12 00:05:13 All files added for processing to worker queue...
    2022/09/12 00:05:14 All files processed...
    2022/09/12 00:05:14 All workers stopped...
    2022/09/12 00:05:14 Writing container images, file mapping to output file...
    2022/09/12 00:05:14 Printing name of distinct container image names...


    Distinct container images:
    --------------------------
    mcr.microsoft.com/dotnet/sdk:5.0
    golang:1.15
    scratch
    openjdk:17-slim
    mcr.microsoft.com/dotnet/sdk:6.0
    mcr.microsoft.com/dotnet/aspnet:6.0-alpine
    ghcr.io/cse-labs/webvalidate:latest
    nginx
    node:lts
    ghcr.io/cse-labs/k3d:latest
    openjdk:8-jdk-alpine
    alpine:latest
    golang:1.10-stretch
    openjdk:18-jdk-alpine3.14
    mcr.microsoft.com/vscode/devcontainers/universal:1-focal
    mcr.microsoft.com/azure-cli
    python:2.7-slim

    openjdk:18-slim-buster
    node:lts-alpine
    maven:3-openjdk-17-slim
    mcr.microsoft.com/dotnet/aspnet:5.0
    gcr.io/distroless/static:nonroot
  ```
