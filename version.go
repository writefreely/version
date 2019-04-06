package version

import (
	"context"
	"github.com/google/go-github/github"
)

func GetLatest(org, repo string) (string, error) {
	c := github.NewClient(nil)
	rel, _, err := c.Repositories.GetLatestRelease(context.Background(), org, repo)
	if err != nil {
		return "", err
	}
	return *rel.TagName, nil
}
