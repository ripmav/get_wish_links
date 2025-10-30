package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/ripmav/get_wish_links/internal"
	"github.com/ripmav/get_wish_links/internal/extract"
	"github.com/ripmav/get_wish_links/internal/finder"
	"github.com/ripmav/get_wish_links/internal/gacha"

	"github.com/alecthomas/kong"
)

type Application struct {
	Verbose bool `help:"show verbose output" short:"v" name:"verbose" env:"VERBOSE" default:"false"`

	Root      string `help:"base path where 'Genshin Impact' directory resides" short:"r" name:"root"`
	UrlFilter string `help:"filter urls by this string" short:"f" default:"gacha_info/api/getGachaLog" name:"url-filter"`
}

func main() {
	ctx := internal.NewContextWithSignal(context.Background())

	options := []kong.Option{
		kong.Name("Get Genshin Wishes"),
		kong.Description("retrieve genshin wishes from the game"),
		kong.UsageOnError(),
	}

	app := &Application{}
	ktx := kong.Parse(app, options...)
	ktx.BindTo(ctx, (*context.Context)(nil))

	slog.SetDefault(slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: app.Verbose,
		})),
	)

	// TODO: Cleanup

	if strings.TrimSpace(app.Root) == "" {
		slog.Error("root path cannot be empty", "root", app.Root)
		os.Exit(1)
	}
	path, err := finder.FindData2(app.Root)
	if err != nil {
		slog.Error("failed to locate data_2", "error", err, "root", app.Root)
		os.Exit(1)
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		slog.Error("failed to open file", "error", err, "path", path)
		os.Exit(1)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	urls, err := extract.UrLs(r)

	if err != nil {
		slog.Error("failed to extract urls", "error", err)
		os.Exit(1)
	}

	urls = slices.DeleteFunc(urls, func(url string) bool {
		return !strings.Contains(url, app.UrlFilter)
	})

	selected := gacha.OnePerTypeWithEndIdZero(urls)

	// determine deterministic order of types: numeric ascending, non-numeric next lexicographically, and "unknown" last
	keys := make([]string, 0, len(selected))
	for k := range selected {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == "unknown" {
			return false
		}
		if keys[j] == "unknown" {
			return true
		}
		ai, aerr := strconv.Atoi(keys[i])
		aj, berr := strconv.Atoi(keys[j])
		if aerr == nil && berr == nil {
			return ai < aj
		}
		if aerr == nil {
			return true
		}
		if berr == nil {
			return false
		}
		return keys[i] < keys[j]
	})

	fmt.Println(selected[keys[0]])
}
