# log_aggregator

A small Go utility that watches a directory for `*.log` files and tails them in real time. New lines are printed to standard output prefixed with the originating file name.

**Key points**
- Linux-only: uses inotify via `golang.org/x/sys/unix`.
- Watches for new `.log` files in the target directory and automatically starts/stops tailing them.
- Intended for simple aggregation or debugging workflows where multiple log files are produced in one folder.

**Prerequisites**
- Go 1.18+ installed
- Linux (inotify support)

**Build**

```bash
go build -o log_aggregator .
```

**Run**

```bash
# run the binary and pass the directory to watch
./log_aggregator /path/to/log/dir

# or run directly with `go run`
go run main.go /path/to/log/dir
```

The program expects a single argument: the directory to watch. It will tail existing `*.log` files in that directory and begin tailing any new `*.log` files that appear.

**Files of interest**
- `main.go` — watcher and coordinator (entry point)
- `tail_file.go` — tailing logic for individual files
- `watchtest.go` — test/demo helpers

**Behavior notes**
- Lines are printed to `stdout` prefixed with the file name.
- When the program receives SIGINT or SIGTERM it will stop watching and wait for active tailers to exit cleanly.

Contributions and issues welcome.

---

© Short README — add license and contribution guidelines as needed.
