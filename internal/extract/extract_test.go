package extract

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestExtractURLs_SingleWithPrefix(t *testing.T) {
	data := []byte("\x00random\x01data1/0/https://public-operation-hk4e-sg.hoyoverse.com/gacha_info/api/getGachaLog?win_mode=fullscreen\x00tail")

	urls, err := UrLs(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://public-operation-hk4e-sg.hoyoverse.com/gacha_info/api/getGachaLog?win_mode=fullscreen",
	}, urls)
}

func TestExtractURLs_MultipleMixedAndKeepOrder(t *testing.T) {
	data := []byte("junk1/0/https://a.example/path?a=1&b=2 more-bytes https://b.example/ok#frag\x00end")

	urls, err := UrLs(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://a.example/path?a=1&b=2",
		"https://b.example/ok#frag",
	}, urls)
}

func TestExtractURLs_WithoutPrefixStillCaptured(t *testing.T) {
	data := []byte("prefix https://no-prefix.example/works too")

	urls, err := UrLs(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://no-prefix.example/works",
	}, urls)
}

func TestExtractURLs_StopsAtInvalidChars(t *testing.T) {
	// URL followed by a space and a non-URL char; extraction should stop before the space
	data := []byte("1/0/https://host.tld/abc_def?x=1 y trailing")

	urls, err := UrLs(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://host.tld/abc_def?x=1",
	}, urls)
}

func TestExtractURLs_NoUrlsReturnsEmpty(t *testing.T) {
	data := []byte("binary\x00\x01\x02no links here")

	urls, err := UrLs(bytes.NewReader(data))
	require.NoError(t, err)
	require.Empty(t, urls)
}

func TestExtractURLs_FromFSData2(t *testing.T) {
	// Build a fake filesystem with a data_2 file at the expected path
	const data2Path = "Genshin Impact/GenshinImpact_Data/webCaches/2.0.0.0/Cache/Cache_Data/data_2"

	content := []byte(
		"\x00\x00garbage1/0/https://u.example/a?x=1&y=2\x00more1/0/https://another.example/b#c\x01",
	)

	fsys := fstest.MapFS{
		"Genshin Impact":                                                       {Mode: fs.ModeDir},
		"Genshin Impact/GenshinImpact_Data":                                    {Mode: fs.ModeDir},
		"Genshin Impact/GenshinImpact_Data/webCaches":                          {Mode: fs.ModeDir},
		"Genshin Impact/GenshinImpact_Data/webCaches/2.0.0.0":                  {Mode: fs.ModeDir},
		"Genshin Impact/GenshinImpact_Data/webCaches/2.0.0.0/Cache":            {Mode: fs.ModeDir},
		"Genshin Impact/GenshinImpact_Data/webCaches/2.0.0.0/Cache/Cache_Data": {Mode: fs.ModeDir},
		data2Path: {Data: content},
	}

	f, err := fsys.Open(data2Path)
	require.NoError(t, err)
	defer f.Close()

	urls, err := UrLs(f)
	require.NoError(t, err)
	require.Equal(t, []string{
		"https://u.example/a?x=1&y=2",
		"https://another.example/b#c",
	}, urls)
}
