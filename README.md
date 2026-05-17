# 📦 Android Package Manager

Fast Go CLI for querying and managing Android packages via ADB.

## Features
- List all / user-only packages
- Search by name, size, install date
- Categorize (system, user, updated, bloat)
- Export to JSON/CSV
- Estimate storage usage per package
- Parallel queries (10x faster than ADB alone)

## Install
```bash
go install github.com/OutrageousStorm/android-package-manager@latest
apm --help
```

## Usage
```bash
apm list              # all packages
apm list --user      # user-installed only
apm search firefox    # find packages matching "firefox"
apm size              # total storage per package
apm export packages.json
```
