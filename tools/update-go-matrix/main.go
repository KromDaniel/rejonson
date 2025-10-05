package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	downloadsAPI  = "https://go.dev/dl/?mode=json&include=all"
	matrixPattern = `go-version:\s*\[[^\]]*\]`
)

type goRelease struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

type majorMinor struct {
	Major int
	Minor int
}

func main() {
	versions, err := latestMinorVersions(4)
	if err != nil {
		panic(fmt.Errorf("fetching go versions: %w", err))
	}

	ciPath, err := workflowPath()
	if err != nil {
		panic(err)
	}

	changed, err := updateMatrix(ciPath, versions)
	if err != nil {
		panic(fmt.Errorf("updating matrix: %w", err))
	}

	if changed {
		fmt.Printf("Updated Go version matrix to: %s\n", strings.Join(versions, ", "))
	} else {
		fmt.Println("Go version matrix already up to date")
	}
}

func latestMinorVersions(limit int) ([]string, error) {
	releases, err := fetchReleases()
	if err != nil {
		return nil, err
	}

	seen := make(map[majorMinor]bool)
	var ordered []majorMinor

	for _, rel := range releases {
		if !rel.Stable {
			continue
		}

		mm, err := parseMajorMinor(rel.Version)
		if err != nil {
			continue
		}

		if seen[mm] {
			continue
		}
		seen[mm] = true
		ordered = append(ordered, mm)
	}

	if len(ordered) == 0 {
		return nil, errors.New("no stable Go releases found")
	}

	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].Major == ordered[j].Major {
			return ordered[i].Minor > ordered[j].Minor
		}
		return ordered[i].Major > ordered[j].Major
	})

	if len(ordered) > limit {
		ordered = ordered[:limit]
	}

	versions := make([]string, len(ordered))
	for i, mm := range ordered {
		versions[i] = fmt.Sprintf("%d.%d.x", mm.Major, mm.Minor)
	}

	return versions, nil
}

func fetchReleases() ([]goRelease, error) {
	resp, err := http.Get(downloadsAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var releases []goRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func parseMajorMinor(raw string) (majorMinor, error) {
	s := strings.TrimPrefix(raw, "go")
	s = strings.SplitN(s, "rc", 2)[0]
	s = strings.SplitN(s, "beta", 2)[0]

	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return majorMinor{}, fmt.Errorf("invalid version: %q", raw)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return majorMinor{}, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return majorMinor{}, err
	}

	return majorMinor{Major: major, Minor: minor}, nil
}

func workflowPath() (string, error) {
	base := ".github/workflows/ci.yml"
	if _, err := os.Stat(base); err == nil {
		return base, nil
	}

	// Attempt to locate relative to repo root for safety when invoked elsewhere.
	if dir, err := os.Getwd(); err == nil {
		candidate := filepath.Join(dir, base)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("workflow file %q not found", base)
}

func updateMatrix(path string, versions []string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	quoted := make([]string, len(versions))
	for i, v := range versions {
		quoted[i] = fmt.Sprintf("\"%s\"", v)
	}

	replacement := fmt.Sprintf("go-version: [%s]", strings.Join(quoted, ", "))

	re := regexp.MustCompile(matrixPattern)
	updated := re.ReplaceAll(data, []byte(replacement))

	if bytes.Equal(data, updated) {
		return false, nil
	}

	return true, os.WriteFile(path, updated, 0o644)
}
