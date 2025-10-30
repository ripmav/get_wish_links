package finder

import (
    "io/fs"
    "testing"
    "testing/fstest"

    "github.com/stretchr/testify/require"
)

// dir is a helper to declare a directory entry in fstest.MapFS
func dir() *fstest.MapFile { return &fstest.MapFile{Mode: fs.ModeDir} }

func TestFindData2_SuccessSelectsHighestVersion(t *testing.T) {
    fsys := fstest.MapFS{
        "Genshin Impact":                                              dir(),
        "Genshin Impact/GenshinImpact_Data":                           dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches":                 dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/foo":             dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.43.0.0":        dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.0.0":        dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0":        dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache":  dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache/Cache_Data": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache/Cache_Data/data_2": {Data: []byte("test")},
    }

    expected := "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache/Cache_Data/data_2"

    got, err := findData2FS(fsys)
    require.NoError(t, err)
    require.Equal(t, expected, got)
}

func TestFindData2_NoVersionFolders(t *testing.T) {
    fsys := fstest.MapFS{
        "Genshin Impact":                                      dir(),
        "Genshin Impact/GenshinImpact_Data":                   dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches":         dir(),
    }

    _, err := findData2FS(fsys)
    require.Error(t, err)
}

func TestFindData2_Data2Missing(t *testing.T) {
    fsys := fstest.MapFS{
        "Genshin Impact":                                      dir(),
        "Genshin Impact/GenshinImpact_Data":                   dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches":         dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/2.44.1.0/Cache/Cache_Data": dir(),
    }

    _, err := findData2FS(fsys)
    require.Error(t, err)
}

func TestFindData2_IgnoresNonVersionDirs(t *testing.T) {
    fsys := fstest.MapFS{
        "Genshin Impact":                                      dir(),
        "Genshin Impact/GenshinImpact_Data":                   dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches":         dir(),
        // Noise directories
        "Genshin Impact/GenshinImpact_Data/webCaches/latest":  dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/backup-2024": dir(),
        // Valid versions
        "Genshin Impact/GenshinImpact_Data/webCaches/2.2.0.0": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/9.0.0.0": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/9.0.0.0/Cache": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/9.0.0.0/Cache/Cache_Data": dir(),
        "Genshin Impact/GenshinImpact_Data/webCaches/9.0.0.0/Cache/Cache_Data/data_2": {Data: []byte("test")},
    }

    expected := "Genshin Impact/GenshinImpact_Data/webCaches/9.0.0.0/Cache/Cache_Data/data_2"

    got, err := findData2FS(fsys)
    require.NoError(t, err)
    require.Equal(t, expected, got)
}
