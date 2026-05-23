package sync

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cthulhu/dot-agent/internal/config"
	"github.com/cthulhu/dot-agent/internal/ignore"
	"github.com/cthulhu/dot-agent/internal/paths"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Direction int

const (
	LocalToSource Direction = iota
	SourceToLocal
)

type Options struct {
	DryRun bool
	Backup bool
	Force  bool
}

type FileAction struct {
	RelPath string
	Action  string // add, update, delete, skip, blocked, ignored
}

type Result struct {
	Actions []FileAction
}

func ResolvePaths(sourceRoot string, entry config.AssistantEntry) (sourceDir, targetDir string, err error) {
	sourceDir = filepath.Join(sourceRoot, filepath.FromSlash(entry.Source))
	targetDir, err = paths.ExpandPath(entry.Target)
	if err != nil {
		return "", "", err
	}
	return sourceDir, targetDir, nil
}

func Add(sourceRoot string, entry config.AssistantEntry, opts Options) (*Result, error) {
	return syncTrees(sourceRoot, entry, LocalToSource, opts)
}

func Apply(sourceRoot string, entry config.AssistantEntry, opts Options) (*Result, error) {
	return syncTrees(sourceRoot, entry, SourceToLocal, opts)
}

func syncTrees(sourceRoot string, entry config.AssistantEntry, dir Direction, opts Options) (*Result, error) {
	sourceDir, targetDir, err := ResolvePaths(sourceRoot, entry)
	if err != nil {
		return nil, err
	}

	matcher := ignore.New(entry.Ignore...)

	var fromRoot, toRoot string
	switch dir {
	case LocalToSource:
		fromRoot, toRoot = targetDir, sourceDir
	case SourceToLocal:
		fromRoot, toRoot = sourceDir, targetDir
	default:
		return nil, fmt.Errorf("invalid direction")
	}

	result := &Result{}
	if err := os.MkdirAll(toRoot, 0o755); err != nil && !opts.DryRun {
		return nil, err
	}

	err = filepath.WalkDir(fromRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}

		rel, err := filepath.Rel(fromRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		relSlash := filepath.ToSlash(rel)

		if d.IsDir() {
			ignored, err := matcher.Ignored(relSlash)
			if err != nil {
				return err
			}
			if ignored {
				return filepath.SkipDir
			}
			dest := filepath.Join(toRoot, rel)
			if !opts.DryRun {
				if err := os.MkdirAll(dest, 0o755); err != nil {
					return err
				}
			}
			return nil
		}

		if dir == LocalToSource {
			if ignore.Blocked(relSlash) {
				result.Actions = append(result.Actions, FileAction{RelPath: relSlash, Action: "blocked"})
				return nil
			}
			ignored, err := matcher.Ignored(relSlash)
			if err != nil {
				return err
			}
			if ignored {
				result.Actions = append(result.Actions, FileAction{RelPath: relSlash, Action: "ignored"})
				return nil
			}
		}

		dest := filepath.Join(toRoot, rel)
		action, err := copyFile(path, dest, opts)
		if err != nil {
			return err
		}
		result.Actions = append(result.Actions, FileAction{RelPath: relSlash, Action: action})
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if dir == LocalToSource {
		if err := pruneOrphans(sourceDir, targetDir, matcher, opts, result); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func pruneOrphans(sourceDir, targetDir string, matcher *ignore.Matcher, opts Options, result *Result) error {
	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		rel, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if rel == "." || d.IsDir() {
			return nil
		}
		relSlash := filepath.ToSlash(rel)
		localPath := filepath.Join(targetDir, rel)
		if _, err := os.Stat(localPath); err == nil {
			return nil
		}
		if !opts.DryRun {
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		result.Actions = append(result.Actions, FileAction{RelPath: relSlash, Action: "delete"})
		return nil
	})
}

func copyFile(src, dest string, opts Options) (string, error) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return "", err
	}

	action := "add"
	if _, err := os.Stat(dest); err == nil {
		same, err := filesEqual(src, dest)
		if err != nil {
			return "", err
		}
		if same {
			return "skip", nil
		}
		action = "update"
	}

	if opts.DryRun {
		return action, nil
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}

	if opts.Backup && action == "update" {
		if err := backupFile(dest); err != nil {
			return "", err
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", err
	}
	return action, nil
}

func backupFile(path string) error {
	backup := path + ".bak." + time.Now().Format("20060102-150405")
	return copyPath(path, backup)
}

func copyPath(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func filesEqual(a, b string) (bool, error) {
	af, err := os.ReadFile(a)
	if err != nil {
		return false, err
	}
	bf, err := os.ReadFile(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(af, bf), nil
}

type DriftEntry struct {
	RelPath string
	Status  string // only_local, only_source, modified
}

type DriftReport struct {
	Entries []DriftEntry
}

func Compare(sourceRoot string, entry config.AssistantEntry) (*DriftReport, error) {
	sourceDir, targetDir, err := ResolvePaths(sourceRoot, entry)
	if err != nil {
		return nil, err
	}
	matcher := ignore.New(entry.Ignore...)
	report := &DriftReport{}

	localFiles := map[string]struct{}{}
	sourceFiles := map[string]struct{}{}

	if err := collectFiles(targetDir, matcher, localFiles, true); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err := collectFiles(sourceDir, matcher, sourceFiles, false); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	all := map[string]struct{}{}
	for k := range localFiles {
		all[k] = struct{}{}
	}
	for k := range sourceFiles {
		all[k] = struct{}{}
	}

	for rel := range all {
		_, inLocal := localFiles[rel]
		_, inSource := sourceFiles[rel]
		switch {
		case inLocal && !inSource:
			report.Entries = append(report.Entries, DriftEntry{RelPath: rel, Status: "only_local"})
		case !inLocal && inSource:
			report.Entries = append(report.Entries, DriftEntry{RelPath: rel, Status: "only_source"})
		default:
			same, err := filesEqual(filepath.Join(targetDir, rel), filepath.Join(sourceDir, rel))
			if err != nil {
				return nil, err
			}
			if !same {
				report.Entries = append(report.Entries, DriftEntry{RelPath: rel, Status: "modified"})
			}
		}
	}
	return report, nil
}

func collectFiles(root string, matcher *ignore.Matcher, out map[string]struct{}, checkBlocked bool) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		relSlash := filepath.ToSlash(rel)
		if d.IsDir() {
			ignored, err := matcher.Ignored(relSlash)
			if err != nil {
				return err
			}
			if ignored {
				return filepath.SkipDir
			}
			return nil
		}
		if checkBlocked && ignore.Blocked(relSlash) {
			return nil
		}
		ignored, err := matcher.Ignored(relSlash)
		if err != nil {
			return err
		}
		if ignored {
			return nil
		}
		out[relSlash] = struct{}{}
		return nil
	})
}

func Diff(sourceRoot string, entry config.AssistantEntry) (string, error) {
	sourceDir, targetDir, err := ResolvePaths(sourceRoot, entry)
	if err != nil {
		return "", err
	}
	matcher := ignore.New(entry.Ignore...)
	dmp := diffmatchpatch.New()

	var buf strings.Builder
	localFiles := map[string]struct{}{}
	sourceFiles := map[string]struct{}{}
	if err := collectFiles(targetDir, matcher, localFiles, true); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if err := collectFiles(sourceDir, matcher, sourceFiles, false); err != nil && !os.IsNotExist(err) {
		return "", err
	}

	all := map[string]struct{}{}
	for k := range localFiles {
		all[k] = struct{}{}
	}
	for k := range sourceFiles {
		all[k] = struct{}{}
	}

	for rel := range all {
		localPath := filepath.Join(targetDir, rel)
		sourcePath := filepath.Join(sourceDir, rel)
		_, localOK := localFiles[rel]
		_, sourceOK := sourceFiles[rel]

		switch {
		case localOK && !sourceOK:
			fmt.Fprintf(&buf, "--- local only: %s\n", rel)
		case !localOK && sourceOK:
			fmt.Fprintf(&buf, "+++ source only: %s\n", rel)
		default:
			lb, err := os.ReadFile(localPath)
			if err != nil {
				return "", err
			}
			sb, err := os.ReadFile(sourcePath)
			if err != nil {
				return "", err
			}
			if bytes.Equal(lb, sb) {
				continue
			}
			diffs := dmp.DiffMain(string(sb), string(lb), true)
			fmt.Fprintf(&buf, "diff %s (source -> local)\n%s\n", rel, dmp.DiffPrettyText(diffs))
		}
	}
	return buf.String(), nil
}

func PrintResult(result *Result) {
	counts := map[string]int{}
	for _, a := range result.Actions {
		counts[a.Action]++
		if a.Action != "skip" && a.Action != "ignored" {
			fmt.Printf("  %s: %s\n", a.Action, a.RelPath)
		}
	}
	fmt.Printf("\nSummary: ")
	parts := make([]string, 0, len(counts))
	for _, k := range []string{"add", "update", "delete", "blocked", "ignored", "skip"} {
		if n := counts[k]; n > 0 {
			parts = append(parts, fmt.Sprintf("%s=%d", k, n))
		}
	}
	fmt.Println(strings.Join(parts, ", "))
}

func PrintDrift(report *DriftReport) {
	if len(report.Entries) == 0 {
		fmt.Println("No drift detected.")
		return
	}
	for _, e := range report.Entries {
		fmt.Printf("  %s: %s\n", e.Status, e.RelPath)
	}
	fmt.Printf("\n%d file(s) differ.\n", len(report.Entries))
}
