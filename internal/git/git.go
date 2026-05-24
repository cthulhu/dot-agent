package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RequireGit() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git not found in PATH; install git to use dot-agent")
	}
	return nil
}

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

func Init(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
		return nil
	}
	_, err := run(dir, "init")
	return err
}

func Clone(url, dir string) error {
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", url, dir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

func AddAll(dir string) error {
	_, err := run(dir, "add", "-A")
	return err
}

func IsDirty(dir string) (bool, error) {
	out, err := run(dir, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return out != "", nil
}

func Commit(dir, message string) error {
	dirty, err := IsDirty(dir)
	if err != nil {
		return err
	}
	if !dirty {
		return nil
	}
	_, err = run(dir, "commit", "-m", message)
	return err
}

func Push(dir string) error {
	branch, err := CurrentBranch(dir)
	if err != nil {
		return err
	}
	_, err = run(dir, "push", "-u", "origin", branch)
	return err
}

func Pull(dir string) error {
	branch, err := CurrentBranch(dir)
	if err != nil {
		return err
	}
	_, err = run(dir, "pull", "--ff-only", "origin", branch)
	return err
}

func CurrentBranch(dir string) (string, error) {
	return run(dir, "rev-parse", "--abbrev-ref", "HEAD")
}

func StatusPorcelain(dir string) (string, error) {
	return run(dir, "status", "--porcelain")
}

func RemoteURL(dir string) (string, error) {
	return run(dir, "remote", "get-url", "origin")
}

func SetRemote(dir, url string) error {
	if _, err := run(dir, "remote", "get-url", "origin"); err != nil {
		_, err = run(dir, "remote", "add", "origin", url)
		return err
	}
	_, err := run(dir, "remote", "set-url", "origin", url)
	return err
}
