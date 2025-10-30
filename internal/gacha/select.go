package gacha

import (
	neturl "net/url"
)

func OnePerTypeWithEndIdZero(urls []string) map[string]string {
	out := make(map[string]string)
	for _, u := range urls {
		parsed, err := neturl.Parse(u)
		if err != nil {
			continue // ignore malformed; cannot test required params
		}
		q := parsed.Query()
		if q.Get("end_id") != "0" {
			continue
		}
		gt := q.Get("gacha_type")
		if gt == "" {
			gt = "unknown"
		}
		out[gt] = u
	}
	return out
}
