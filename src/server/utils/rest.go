package utils

import (
	"fmt"

	"github.com/torchiaf/code-editor/server/models"
)

func GitInfo(git models.GitSource) (string, string) {
	return fmt.Sprintf("https://github.com/%s/%s/%s", git.Org, git.Repo, git.Branch), fmt.Sprintf("folder=/git/%s", git.Repo)
}
