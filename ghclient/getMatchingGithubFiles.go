package ghclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type MatchingFileFilter struct {
	Org      string
	Repo     string
	FileName string
	GHToken  string
}

type FileDetails struct {
	FileName        string `json:"fileName"`
	Repository      string `json:"repository"`
	Path            string `json:"path"`
	FileDownloadURL string `json:"fileDownloadURL"`
	BlobAccessToken string `json:"blobAccessToken"`
}

func getGHClient(ctx context.Context, ghToken string) *github.Client {

	if ghToken == "" {
		return github.NewClient(nil)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func getRawBlobUrl(HTMLURL string, blobAccessToken string) string {
	rawURL := strings.Replace(HTMLURL, "/blob/", "/", 1)
	rawURL = strings.Replace(rawURL, "github.com", "raw.githubusercontent.com", 1)
	if blobAccessToken != "" {
		rawURL = strings.Replace(rawURL, "https://", "https://"+blobAccessToken+"@", 1)
	}
	return rawURL
}

func (f *MatchingFileFilter) GetDownloadableFileNames(ctx context.Context) ([]FileDetails, error) {
	var matchingFiles []FileDetails
	client := getGHClient(ctx, f.GHToken)
	query := fmt.Sprintf("org:%s filename:%s", f.Org, f.FileName)

	if f.Repo != "" {
		query = fmt.Sprintf("repo:%s/%s filename:%s", f.Org, f.Repo, f.FileName)
	}

	searchOpt := &github.SearchOptions{}
	res, resp, err := client.Search.Code(ctx, query, searchOpt)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	for {
		for _, r := range res.CodeResults {
			fd := FileDetails{
				FileName:        *r.Name,
				Path:            *r.Path,
				Repository:      *r.Repository.Name,
				FileDownloadURL: getRawBlobUrl(*r.HTMLURL, f.GHToken),
			}
			matchingFiles = append(matchingFiles, fd)
		}
		if resp.NextPage != 0 {
			searchOpt.ListOptions.Page = resp.NextPage
			res, resp, err = client.Search.Code(ctx, query, searchOpt)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return matchingFiles, nil
}
