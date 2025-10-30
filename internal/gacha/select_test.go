package gacha

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// keys returns the set of keys of the provided map in arbitrary order.
func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestOnePerTypeWithPage1End0_SelectsLastPerTypeAndFiltersByParams(t *testing.T) {
	urls := []string{
		// params for type=1 with different page value -> still acceptable because page is optional
		"https://example.com/gacha?gacha_type=1&page=2&end_id=0",
		// also acceptable for type=1 (first candidate with page=1)
		"https://example.com/gacha?gacha_type=1&page=1&end_id=0&x=a",
		// another acceptable for type=1 (should win because it's last)
		"https://example.com/gacha?x=y&gacha_type=1&end_id=0&page=1",
		// unknown type, correct params
		"https://example.com/gacha?foo=bar&page=1&end_id=0",
		// unknown type, wrong params (ignored due to end_id!=0)
		"https://example.com/gacha?foo=bar&page=1&end_id=9",
		// type=11 with page!=1 is still acceptable since page is optional
		"https://example.com/gacha?gacha_type=11&page=3&end_id=0",
		// malformed URL ignored
		"not a url at all",
	}

	got := OnePerTypeWithEndIdZero(urls)

	// Expect keys 1, 11 and unknown
	require.ElementsMatch(t, []string{"1", "11", "unknown"}, keys(got))

	require.Equal(t, "https://example.com/gacha?x=y&gacha_type=1&end_id=0&page=1", got["1"]) // last correct wins
	require.Equal(t, "https://example.com/gacha?gacha_type=11&page=3&end_id=0", got["11"])   // accepted though page!=1
	require.Equal(t, "https://example.com/gacha?foo=bar&page=1&end_id=0", got["unknown"])    // last correct unknown
}

func TestOnePerTypeWithPageOptional_EndIDRequired(t *testing.T) {
	urls := []string{
		"https://host/path?gacha_type=301&page=1",          // missing end_id -> excluded
		"https://host/path?gacha_type=301&end_id=0",        // missing page -> accepted (page optional)
		"https://host/path?gacha_type=301&page=1&end_id=1", // wrong end_id -> excluded
	}
	got := OnePerTypeWithEndIdZero(urls)
	require.Equal(t, map[string]string{
		"301": "https://host/path?gacha_type=301&end_id=0",
	}, got)
}

func TestOnePerTypeWithPage1End0_DuplicatesAndOrder(t *testing.T) {
	urls := []string{
		"https://x/y?gacha_type=12&page=1&end_id=0",
		"https://x/y?gacha_type=12&page=1&end_id=0&a=b",
	}
	got := OnePerTypeWithEndIdZero(urls)
	require.Equal(t, "https://x/y?gacha_type=12&page=1&end_id=0&a=b", got["12"]) // last one wins
}
