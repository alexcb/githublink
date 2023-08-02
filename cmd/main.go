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
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func runCommandTrimmedOutput(args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no command given")
	}
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
func runCommandSplitLines(args ...string) ([]string, error) {
	out, err := runCommandTrimmedOutput(args...)
	if err != nil {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}

func getGitSha() (string, error) {
	return runCommandTrimmedOutput("git", "rev-parse", "HEAD")
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

func getGithubURL(gitURL, gitSha, path string, line int) (string, error) {
	var urlPath string
	if strings.HasPrefix(gitURL, sshGitPrefix) {
		urlPath = strings.TrimPrefix(gitURL, sshGitPrefix)
	} else {
		u, err := url.Parse(gitURL)
		if err != nil {
			return "", err
		}
		urlPath = u.Path
	}
	user, repo, err := getUserAndRepo(urlPath)
	if err != nil {
		return "", err
	}
	return formatGithubURL(user, repo, gitSha, path, line), nil
}

func getFullRepoPath(path string) (string, error) {
	return runCommandTrimmedOutput("git", "ls-files", "--full-name", "--error-unmatch", path)
}

func main() {
	progName := "githublink"
	if len(os.Args) > 0 {
		progName = os.Args[0]
	}

	if len(os.Args) != 2 && len(os.Args) != 3 {
		die("usage: %s <path> [<line-number>]\n", progName)
	}

	line := -1
	path := os.Args[1]
	if len(os.Args) == 3 {
		var err error
		line, err = strconv.Atoi(os.Args[2])
		if err != nil {
			die("failed to convert %s to int: %v\n", os.Args[2], err)
		}
	}

	fullPath, err := getFullRepoPath(path)
	if err != nil {
		die("%s is not tracked by git: %v\n", path, err)
	}

	gitSha, err := getGitSha()
	if err != nil {
		die("failed to get git sha: %v\n", err)
	}
	remoteBranches, err := getRemoteBranches(gitSha)
	if err != nil {
		die("failed to get remote branches: %v\n", err)
	}
	if len(remoteBranches) == 0 {
		die("current git commit doesn't exist on any remote branches")
	}
	remoteName := strings.Split(remoteBranches[0], "/")[0]
	remoteURL, err := getRemoteURL(remoteName)
	if err != nil {
		die("failed to get remote url: %v\n", err)
	}
	webURL, err := getGithubURL(remoteURL, gitSha, fullPath, line)
	if err != nil {
		die("failed to get remote url: %v\n", err)
	}
	err = browser.OpenURL(webURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open url automatically: %s\n", err)
		fmt.Println(webURL)
	}
}
