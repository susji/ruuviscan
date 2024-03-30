set xdata time
set timefmt "%s"
set format x "%Y-%m-%d %H:%M"
set xtics rotate by -60
set datafile separator "|"
set key outside
set title 'RuuviTag Temperatures [°C]'
set ytics 5.0
set mytics 5
set grid xtics mxtics ytics mytics

SqliteField(f, m) = "< sqlite3 --readonly temps.db \"SELECT ".f." FROM temps WHERE mac='".m."' ORDER BY timestamp DESC\""
fields = "strftime('%s',timestamp), ROUND(temp, 2)"
macs = system('sqlite3 --readonly -column temps.db "SELECT DISTINCT(mac) FROM temps" | tr "\n" " "')
print macs
plot for [mac in macs] SqliteField(fields, mac) using 1:2 title mac[13:14].mac[16:17] with lines