package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cli/browser"
)

func die(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func runCommandTrimmedOutput(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command given")
	}
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("command %v failed: %w", args, err)
	}
	return strings.TrimSpace(string(out)), nil
}
func runCommandSplitLines(args ...string) ([]string, error) {
	out, err := runCommandTrimmedOutput(args...)
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, l := range strings.Split(out, "\n") {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines, nil
}

func getGitSha() (string, error) {
	return runCommandTrimmedOutput("git", "rev-parse", "HEAD")
}

func isGitSha(sha string) bool {
	output, err := runCommandTrimmedOutput("git", "rev-parse", sha)
	if err != nil {
		return false
	}
	return output == sha
}

func getRemoteBranches(sha string) ([]string, error) {
	return runCommandSplitLines("git", "branch", "-r", "--contains", sha)
}

func getRemoteURL(remoteName string) (string, error) {
	return runCommandTrimmedOutput("git", "config", "--get", "remote."+remoteName+".url")
}

const sshGitPrefix = "git@github.com:"

func formatGithubURL(user, repo, gitSha, path string, line int) string {
	url := fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s", user, repo, gitSha, path)
	if line >= 0 {
		url += fmt.Sprintf("#L%d", line)
	}
	return url
}

func getUserAndRepo(s string) (string, string, error) {
	parts := strings.Split(strings.TrimSuffix(s, ".git"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("failed to split %s", s)
	}
	return parts[0], parts[1], nil
}

func getGithubUserAndRepo(gitURL string) (string, string, error) {
	var urlPath string
	if strings.HasPrefix(gitURL, sshGitPrefix) {
		urlPath = strings.TrimPrefix(gitURL, sshGitPrefix)
	} else {
		u, err := url.Parse(gitURL)
		if err != nil {
			return "", "", err
		}
		urlPath = u.Path
	}
	user, repo, err := getUserAndRepo(urlPath)
	if err != nil {
		return "", "", err
	}
	return user, repo, nil
}

func getGithubURL(gitURL, gitSha, path string, line int) (string, error) {
	user, repo, err := getGithubUserAndRepo(gitURL)
	if err != nil {
		return "", err
	}
	return formatGithubURL(user, repo, gitSha, path, line), nil
}

// getGithubCommitURL points to a single commit
func getGithubCommitURL(gitURL, gitSha string) (string, error) {
	user, repo, err := getGithubUserAndRepo(gitURL)
	if err != nil {
		return "", err
	}
	return formatGithubCommitURL(user, repo, gitSha), nil
}

func formatGithubCommitURL(user, repo, gitSha string) string {
	return fmt.Sprintf("https://github.com/%s/%s/commit/%s", user, repo, gitSha)
}

func getFullRepoPath(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if fileInfo.IsDir() {
		return path, nil
	}
	return runCommandTrimmedOutput("git", "ls-files", "--full-name", "--error-unmatch", path)
}

// https://github.com/earthly/earthly/commit/dab70be66fecacaa57ba63018dbfa282865eded9

func main() {
	progName := "githublink"
	if len(os.Args) > 0 {
		progName = os.Args[0]
	}

	if len(os.Args) != 2 && len(os.Args) != 3 {
		die("usage: %s <path> [<line-number>]", progName)
	}

	line := -1
	path := os.Args[1]
	if len(os.Args) == 3 {
		var err error
		line, err = strconv.Atoi(os.Args[2])
		if err != nil {
			die("failed to convert %s to int: %v", os.Args[2], err)
		}
	}

	var gitSha string
	showSingleCommit := false

	fullPath, err := getFullRepoPath(path)
	if err != nil {
		if len(path) == 40 && isGitSha(path) {
			gitSha = path
			showSingleCommit = true
		} else {
			die("%s is not tracked by git: %v", path, err)
		}
	}

	if gitSha == "" {
		gitSha, err = getGitSha()
		if err != nil {
			die("failed to get git sha: %v", err)
		}
	}
	remoteBranches, err := getRemoteBranches(gitSha)
	if err != nil {
		die("failed to get remote branches: %v", err)
	}
	if len(remoteBranches) == 0 {
		die("current git commit doesn't exist on any remote branches")
	}
	remoteName := strings.Split(remoteBranches[0], "/")[0]
	remoteURL, err := getRemoteURL(remoteName)
	if err != nil {
		die("failed to get remote url: %v", err)
	}

	var webURL string
	if showSingleCommit {
		webURL, err = getGithubCommitURL(remoteURL, gitSha)
		if err != nil {
			die("failed to get remote url: %v", err)
		}
	} else {
		webURL, err = getGithubURL(remoteURL, gitSha, fullPath, line)
		if err != nil {
			die("failed to get remote url: %v", err)
		}
	}

	err = browser.OpenURL(webURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open url automatically: %s\n", err)
		fmt.Println(webURL)
	}
}
