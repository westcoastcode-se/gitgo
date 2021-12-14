package server

import (
	"errors"
	"regexp"
	"strings"
)

var (
	InvalidGitCommand     = errors.New("invalid git command")
	UnsupportedCommand    = errors.New("command is not a git command")
	InvalidRepositoryPath = errors.New("repository path is not valid")
)

type Command struct {
	Command         string
	Repository      string
	OriginalCommand string
}

func getRepository(s string) (string, bool) {
	idx := strings.Index(s, "'")
	if idx == -1 {
		return "", false
	}
	repository := s[idx+1:]
	idx = strings.Index(repository, "'")
	if idx == -1 {
		return "", false
	}
	repository = repository[0:idx]
	if repository[0] == '/' {
		repository = repository[1:]
	}
	return repository, true
}

var validRepositoryPathRegex = regexp.MustCompile(`^[a-zA-Z\\-_0-9\\.]+$`)

func testRepositoryPath(s string) bool {
	return validRepositoryPathRegex.MatchString(s) &&
		strings.Index(s, "..") == -1
}

func Parse(s string) (*Command, error) {
	parts := SplitOnSpace(s)
	if len(parts) < 2 {
		return &Command{
			OriginalCommand: s,
		}, InvalidGitCommand
	}

	switch parts[0] {
	case "git-upload-pack":
	case "git-upload-archive":
	case "git-receive-pack":
	default:
		return &Command{
			OriginalCommand: s,
		}, UnsupportedCommand
	}

	repository, ok := getRepository(parts[1])
	if !ok {
		return &Command{
			OriginalCommand: s,
		}, InvalidGitCommand
	}

	if !testRepositoryPath(repository) {
		return &Command{
			OriginalCommand: s,
		}, InvalidRepositoryPath
	}

	return &Command{
		Command:         parts[0],
		Repository:      repository,
		OriginalCommand: s,
	}, nil
}
