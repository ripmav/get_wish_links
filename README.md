# Get Genshin Wishes

A command-line tool to extract Genshin Impact wish/gacha history URLs from your local game installation.

## Overview

This tool scans your Genshin Impact game files to locate and extract wish history URLs. 
It finds the most recent valid URLs for each gacha type, making it easy to access your wish history without manually searching through game logs.
So you can use the URL for third-party wish tracking tools (like paimon.moe's auto import), build custom applications, export your complete wish history or whatever you want to do with it.

## Features

- Automatically locates the Genshin Impact `data_2` log file
- Extracts valid wish history URLs from binary cache data
- Filters URLs by gacha type
- Returns the most recent URL for each wish banner type
- Supports verbose logging for debugging

## Installation

### Prerequisites

- Go 1.25 or higher

### Download pre-built binaries

Check the [Releases](https://github.com/ripmav/get_wish_links/releases) page for pre-built binaries for your platform.

### Build from source
```bash
go build -o get_wish_links
```
## Usage
```bash
./get_wish_links [flags]
```
### Parameters

| Parameter | Short | Type | Default | Environment Variable | Required | Description |
|-----------|-------|------|---------|---------------------|----------|-------------|
| `--root` | `-r` | string | - | - | Yes | Base path where 'Genshin Impact' directory resides |
| `--url-filter` | `-f` | string | `gacha_info/api/getGachaLog` | - | No | Filter URLs by this string |
| `--verbose` | `-v` | boolean | `false` | `VERBOSE` | No | Show verbose output with source information |
| `--help` | `-h` | flag | - | - | No | Display help information |

### Parameter Details

#### `--root` / `-r` (Required)

Specifies the base path where the 'Genshin Impact' directory resides. This is the root installation directory of your game.

**Examples:**
- Linux (Wine): `/home/user/.wine/drive_c/Program Files/Genshin\ Impact`
- Windows: `C:\Program Files\Genshin Impact`

**What it does:** The tool searches for the `data_2` file within this directory tree. Specifically, it looks for:
```
{root}/Genshin Impact/GenshinImpact_Data/webCaches/{latest_version}/Cache/Cache_Data/data_2
```
The tool automatically finds the latest version directory (formatted as `x.y.z.w`) and locates the `data_2` file containing the game's web cache logs where wish history URLs are stored.

**Troubleshooting:**
- Make sure the path includes the main "Genshin Impact" folder
- Use quotes if the path contains spaces
- On Windows, you can use either forward slashes (`/`) or backslashes (`\`)
- The tool requires read permissions for the game directory

---

#### `--url-filter` / `-f`

Filters extracted URLs to only include those containing this substring. This helps narrow down the results to only wish history URLs.

**Examples:**
- Default: `gacha_info/api/getGachaLog` (standard wish history API endpoint)
- More permissive: `getGachaLog`
- Stricter: `https://public-operation-hk4e-sg.hoyoverse.com/gacha_info/api/getGachaLog`

**What it does:** After extracting all URLs from the binary `data_2` file, only URLs containing this filter string will be kept. The extraction process:
1. Scans the binary data for `https://` occurrences
2. Captures complete URLs using a permissive character set (RFC3986 compliant)
3. Filters out URLs that don't contain the specified substring
4. Returns only the matching wish history URLs

**Use cases:**
- Keep the default value for normal operation
- Use a shorter filter if you're getting no results (e.g., `gacha`)
- Use a longer, more specific filter if you're getting too many unrelated URLs
- Adjust if miHoYo changes their API endpoint structure in future updates

---

#### `--verbose` / `-v`

Enables verbose logging output with detailed debug information including source file locations.

---

#### `--help` / `-h`

Displays help information about all available parameters and exits.

**Example:**
```shell script
./get_wish_links --help
```

## Complete Examples

### Basic usage (Windows)

```shell script
get_wish_links.exe --root "C:\Program Files\Genshin Impact"
```


### Basic usage (Linux/Wine)

```shell script
./get_wish_links --root "/home/user/.wine/drive_c/Program Files/Genshin\ Impact"
```


### With all parameters

```shell script
get_wish_links.exe \
  --root "C:\Program Files\Genshin Impact" \
  --url-filter "gacha_info/api/getGachaLog" \
  --verbose
```


### Short form

```shell script
get_wish_links.exe -r "C:\Program Files\Genshin Impact" -f "getGachaLog" -v
```


### Using environment variable

```shell script
export VERBOSE=true
get_wish_links.exe -r "C:\Program Files\Genshin Impact"
```

## How It Works

The tool performs the following steps:

1. **Locate data file**:
    - Searches for the `webCaches` directory within the Genshin Impact installation
    - Identifies all version directories (format: `x.y.z.w`)
    - Selects the latest version by comparing version numbers
    - Constructs a path to `data_2` file: `{root}/Genshin Impact/GenshinImpact_Data/webCaches/{latest_version}/Cache/Cache_Data/data_2`

2. **Extract URLs**:
    - Opens and reads the binary `data_2` file
    - Scans for `https://` occurrences
    - Extracts complete URLs using RFC3986-compliant character set
    - Filters URLs to only include those matching `--url-filter`

3. **Group by type**:
    - Parses query parameters from URLs
    - Groups URLs by `gacha_type` parameter (banner type)
    - Handles unknown types gracefully

4. **Select best URL**:
    - For each gacha type, selects the URL with `page=1` and `end_id=0`
    - This represents the starting point for fetching complete wish history
    - Ensures the URL is the most recent and valid for API calls

5. **Output**:
    - Prints the selected URL for the first gacha type (second line)
    - Gacha types are sorted: numeric types ascending, non-numeric lexicographically, "unknown" last

## Output

The tool outputs two lines:
1. The absolute path to the `data_2` file
2. The selected wish history URL

Example output:
```
https://public-operation-hk4e.mihoyo.com/gacha_info/api/getGachaLog?authkey=t7v5Z...&authkey_ver=1&sign_type=2&game_biz=hk4e_global&lang=en&gacha_type=100&page=1&size=20&end_id=0
```


### Understanding the URL

The output URL contains several important parameters:
- `authkey`: Your authentication token (required for API access)
- `gacha_type`: The banner type (see below)
- `page`: Page number (always 1 for starting point)
- `end_id`: Last item ID (always 0 for starting point)
- `game_biz`: Game region (e.g., `hk4e_global`, `hk4e_cn`)
- `lang`: Language preference

### Using the URL

You can use this URL to:
- Import wish history into third-party wish tracking tools
- Query the API directly with tools like `curl` or Postman
- Build custom wish analysis applications
- Export your complete wish history

Example API call:
```shell script
curl "https://public-operation-hk4e.mihoyo.com/gacha_info/api/getGachaLog?authkey=...&page=1&size=20"
```


## Gacha Types

The tool automatically detects and sorts gacha types. Sorting order:
1. Numeric types in ascending order (100, 200, 301, 302)
2. Non-numeric types lexicographically
3. "unknown" type last

### Common Gacha Type Values

| Type  | Banner Name | Description |
|-------|-------------|-------------|
| `100` | Beginner's Wish | Noelle guaranteed banner (20 pulls max) |
| `200` | Standard Wish | Wanderlust Invocation (permanent banner) |
| `301` | Character Event Wish | Limited 5-star character banner |
| `302` | Weapon Event Wish | Limited 5-star weapon banner |
| `400` | Character Event Wish-2 | Second concurrent character banner (when applicable) |
| `500` | Chronicled Wish | Chronicled character banner (when applicable) |
## Troubleshooting

### No URLs found

**Symptoms:** Tool don't prints the URL

**Solutions:**
1. Open Genshin Impact and access the wish history in-game first
2. Try using a shorter `--url-filter`:
```shell script
get_wish_links.exe -r "C:\Program Files\Genshin Impact" -f "gacha"
```
3. Verify the `data_2` file is not empty:
    - Check file size (should not be empty)
    - If empty, access wish history in-game



## License

[The Unlicense](LICENSE)

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## Support

For issues or questions:
- Open an issue on GitHub
- Check existing issues for solutions

## Acknowledgments

- Built with Go 1.25
- Uses [Kong](https://github.com/alecthomas/kong) for CLI parsing
- Inspired by the Genshin Impact community's wish tracking needs

---

**Note:** This tool is not affiliated with, endorsed by, or connected to miHoYo/HoYoverse or Genshin Impact.
