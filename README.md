# ruuviscan

This is a small utility for scanning Bluetooth Low Energy (BLE) Advertisements
from [Ruuvi sensors](https://ruuvi.com) and dumping the values as JSON to
standard output.

# Building

Grab a precompiled [release](https://github.com/susji/ruuviscan/releases) or use
the Go toolchain.

# Usage

Set the environment value `VERBOSE=1` if you wish to see more details when a
RuuviTag Advertisement is received.

## Dumping human-readable temperature values in real time

```
$ ./ruuviscan \
      | jq -rc 'select(.Temperature.Valid) | [.Timestamp, .MAC, .Temperature.Value] | @tsv'
```

## Storing temperature values in an SQLite database

First create a suitable database file with a table for temperature data:

```
$ sqlite3 temps.db <<'EOF'
CREATE TABLE temps (
    timestamp TIMESTAMP NOT NULL,
    mac STRING NOT NULL,
    temp REAL NOT NULL
)
EOF
```

Then you can run the program to scan and store for the values:

```
$ ./ruuviscan \
    | jq --unbuffered -rc \
        'select(.Temperature.Valid) |
             [.Timestamp, .MAC, .Temperature.Value] | @sh' \
    | while read -r _ts _mac _temp; do
          sqlite3 temps.db "INSERT INTO TEMPS VALUES ($_ts, $_mac, $_temp)"
      done
```

You may also observe the latest values as they are recorded:

```
$ watch sqlite3 --readonly \
      "temps.db \
           'SELECT COUNT(*) FROM temps;
            SELECT * FROM temps ORDER BY timestamp DESC LIMIT 10;'"
```

See [temps.plt](./temps.plt) for an example how to plot the temperature data
with gnuplot.
