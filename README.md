# COROS CLI

A command-line application written in **Go** that talks to the COROS Training Hub API. It can authenticate with your COROS account, list your activities, download activity files (FIT, TCX or GPX), and run a basic FIT file inspection.

> This is an independent project, under development. It is not an official COROS product.

## Features

- **Authentication:** logs in to the COROS API with your email and password (MD5-hashed, as the API expects) and `coros auth` saves your credentials (password in the OS keyring), with `.env` / environment variables as a fallback.
- **Activity listing:** prints date, name, label ID and sport type code for your activities.
- **Date filtering:** `--from` / `--to` (`YYYY-MM-DD`) on `list-activities`, `export-activities` and `process-fit`.
- **Activity export:** downloads activity files as `fit`, `tcx` or `gpx` into a directory of your choice.
- **FIT summaries:** downloads activities as FIT files, decodes them in memory with [`muktihari/fit`](https://github.com/muktihari/fit) and prints a JSON summary per activity (see [process-fit](#process-fit)).
- **Automatic pagination:** all result pages of an activity query are fetched (20 per page).

## Tech Stack

| Technology                                          | Purpose                            |
| --------------------------------------------------- | ---------------------------------- |
| Go                                                  | Main language                      |
| [Cobra](https://github.com/spf13/cobra)             | CLI framework                      |
| [muktihari/fit](https://github.com/muktihari/fit)   | FIT file decoding                  |
| [go-keyring](https://github.com/zalando/go-keyring) | Password storage in the OS keyring |
| [godotenv](https://github.com/joho/godotenv)        | Loading `.env`                     |

## Project Structure

```text
cmd/coros/        Entry point
internal/cli/     Cobra commands (auth, list-activities, export-activities, process-fit)
internal/app/     Use cases (list, export, process FIT)
internal/coros/   COROS API client (login, query, download, sport/export types)
internal/fit/     FIT decoding
internal/config/  Configuration loading (saved credentials, env, JSON file, keyring)
```

## Installation

Requires Go 1.26 or newer.

```bash
git clone <repository-url>
cd coros-api-go
go build -o coros ./cmd/coros
```

## Configuration

### Credentials

The easiest way is `coros auth` (see [auth](#auth)), which saves your region, email and password.

Alternatively, copy `.env.example` to `.env` (in the directory you run the binary from) and fill it in:

```bash
COROS_API_URL=https://teamapi.coros.com   # America
# https://teameuapi.coros.com             # Europe
# https://teamcnapi.coros.com             # China
COROS_EMAIL=john.doe@example.com
COROS_PASSWORD=mysecurepassword
```

Each value is resolved independently, in this order:

1. Saved by `coros auth` (API URL and email in `<user config dir>/coros/credentials.json`, password in the OS keyring).
2. Environment variables (`COROS_API_URL`, `COROS_EMAIL`, `COROS_PASSWORD`), including those loaded from `.env`. `.env` is optional.

### `.config/coros/config.json` (optional)

Read from the current working directory. If the file is missing, the defaults below are used.

```json
{
  "output_dir": "./activities",
  "request_timeout_seconds": 30
}
```

- `output_dir`: default directory for `export-activities` when `--output-dir` is not given.
- `request_timeout_seconds` (> 0): timeout for the initial login request.

## Usage

Every command except `auth` logs in to COROS first, so credentials must be configured (via `coros auth` or `.env`).

```bash
./coros --help
```

### auth

Saves your credentials so you don't need a `.env` file.

```bash
./coros auth
```

1. Type your region: `America`, `Europe` or `China` (case-insensitive). The matching API URL is stored.
2. Type your email.
3. Type your password (hidden while typing).

The credentials are verified with a real login before anything is saved. On success you see a confirmation. On failure you see the reason and a hint to copy `.env.example` to `.env` and fill in your credentials instead.

The region and email are saved in your user config directory (e.g. `~/.config/coros/credentials.json`, mode `0600`) and the password in the OS keyring (service `coros-api-go`). Values saved by `auth` take priority over environment variables, which act as the fallback.

### list-activities

Lists activities, one per line: `date - name - labelId - sport type: <code>`.

```bash
./coros list-activities --from 2026-10-01 --to 2026-10-31
```

| Flag     | Description              |
| -------- | ------------------------ |
| `--from` | Start date, `YYYY-MM-DD` |
| `--to`   | End date, `YYYY-MM-DD`   |

### export-activities

Downloads activity files into a directory. Files are named `YYYY-MM-DD <activity name> <labelId>.<ext>`. Failed downloads are logged and counted, and the remaining activities continue to be processed.

```bash
./coros export-activities \
  --from 2026-10-01 \
  --to 2026-10-31 \
  --export-type fit \
  --output-dir ./activities
```

| Flag                 | Description                                                                                       |
| -------------------- | ------------------------------------------------------------------------------------------------- |
| `--output-dir`, `-o` | Directory where files are saved (created if missing). Defaults to `output_dir` from `config.json` |
| `--export-type`      | `fit` (default), `tcx` or `gpx`                                                                   |
| `--from`, `--to`     | Date range, `YYYY-MM-DD`                                                                          |
| `--sport-types`      | Comma-separated sport type names (e.g. `run,trailRun,bike`), see below                            |

Available sport type names: `all`, `run`, `indoorRun`, `trailRun`, `trackRun`, `hike`, `mtnClimb`, `climb`, `bike`, `indoorBike`, `roadEbike`, `gravelRoadBike`, `mountainRiding`, `mountainEbike`, `helmetBike`, `poolSwim`, `openWater`, `triathlon`, `multiSport`, `strength`, `gymCardio`, `gpsCardio`, `ski`, `snowboard`, `xcSki`, `skiTouring`, `skiTouringOld`, `speedsurfing`, `windsurfing`, `row`, `indoorRow`, `whitewater`, `flatwater`, `multiPitch`, `indoorClimb`, `bouldering`, `walk`, `jumpRope`, `climbStairs`, `customSport`.

### process-fit

Downloads activities in the given range as FIT files, decodes each one in memory (FIT files are not saved) and prints one JSON summary per activity, as NDJSON (one object per line). Activities are processed one at a time. A file that fails is logged to stderr and skipped. At the end the command prints `N processed, M failed` to stderr and exits with an error if any failed.

```bash
./coros process-fit --from 2026-10-01 --to 2026-10-31 --sport-types run,trailRun
./coros process-fit --from 2026-10-01 -o summaries.ndjson
./coros process-fit --from 2026-10-01 | jq '.distance_m'
```

| Flag             | Description                                              |
| ---------------- | -------------------------------------------------------- |
| `--from`, `--to` | Date range, `YYYY-MM-DD`                                 |
| `--sport-types`  | Comma-separated sport type names (all sports if omitted) |
| `--output`, `-o` | Write to this file instead of stdout                     |

Each summary comes from the first session of the FIT file:

```json
{
  "label_id": "...",
  "sport": "running",
  "start_time": "2026-10-05T07:30:00Z",
  "distance_m": 10000,
  "timer_time_s": 3000,
  "elapsed_time_s": 3100,
  "calories_kcal": 700,
  "ascent_m": null,
  "descent_m": null,
  "avg_heart_rate_bpm": 150,
  "max_heart_rate_bpm": 175,
  "avg_speed_ms": 3.333,
  "max_speed_ms": null,
  "avg_power_w": null,
  "avg_cadence_rpm": null,
  "laps": 1,
  "records": 2
}
```

Metrics missing from the file (for example power on a swim) are `null`, not `0`. Per-record data (GPS points, per-second heart rate) is not exported yet.

## Disclaimer

This project uses the COROS API to access sports activity data. It is an independent project and is not an official COROS product or affiliated with COROS.
