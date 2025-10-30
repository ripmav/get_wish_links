package finder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// FindData2 searches for the file named "data_2" located under:
//
//	{root}/Genshin Impact/GenshinImpact_Data/webCaches/{latest_version}/Cache/Cache_Data/data_2
//
// The {latest_version} is determined by picking the greatest directory name that
// matches the version schema "x.y.z.w" where each segment is a non-negative integer.
// Returns the absolute path to data_2.
func FindData2(root string) (string, error) {
	if root == "" {
		return "", errors.New("root path must not be empty")
	}

	rel, err := findData2FS(os.DirFS(root))
	if err != nil {
		return "", err
	}
	candidate := filepath.Join(root, filepath.FromSlash(rel))
	abs, err := filepath.Abs(candidate)
	if err != nil {
		return candidate, nil // return as-is if Abs fails
	}
	return abs, nil
}

func findData2FS(fsys fs.FS) (string, error) {
	webCaches := path.Join("GenshinImpact_Data", "webCaches")
	entries, err := fs.ReadDir(fsys, webCaches)
	if err != nil {
		return "", fmt.Errorf("read webCaches dir: %w", err)
	}

	type ver struct {
		name  string
		parts [4]int
	}
	versions := make([]ver, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		p, err := parseVersion(name)
		// ignore invalid version directories
		if err != nil {
			continue
		}
		versions = append(versions, ver{name: name, parts: p})
	}

	if len(versions) == 0 {
		return "", errors.New("no version directories found in webCaches")
	}

	sort.Slice(versions, func(i, j int) bool {
		a, b := versions[i].parts, versions[j].parts
		if a[0] != b[0] {
			return a[0] > b[0]
		}
		if a[1] != b[1] {
			return a[1] > b[1]
		}
		if a[2] != b[2] {
			return a[2] > b[2]
		}
		return a[3] > b[3]
	})

	latest := versions[0].name
	relCandidate := path.Join(webCaches, latest, "Cache", "Cache_Data", "data_2")
	if _, err := fs.Stat(fsys, relCandidate); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("data_2 not found in latest version %s", latest)
		}
		return "", fmt.Errorf("stat data_2: %w", err)
	}
	return relCandidate, nil
}

func parseVersion(s string) ([4]int, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return [4]int{}, fmt.Errorf("invalid version format: %s", s)
	}
	var res [4]int
	for i, seg := range parts {
		if seg == "" {
			return [4]int{}, fmt.Errorf("invalid version segment: %s", seg)
		}
		v, err := strconv.Atoi(seg)
		if err != nil || v < 0 {
			return [4]int{}, fmt.Errorf("invalid version segment: %s", seg)
		}
		res[i] = v
	}
	return res, nil
}
