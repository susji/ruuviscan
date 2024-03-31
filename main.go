// Based on TinyGo's Bluetooth example:
// https://github.com/tinygo-org/bluetooth/blob/12b6f0bc25b665a683fd3143cfaef76dadfb619b/examples/scanner/main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/susji/ruuvi/data"
	"github.com/susji/ruuvi/data/rawv2"
	"tinygo.org/x/bluetooth"
)

var VERBOSE = os.Getenv("VERBOSE")

func main() {
	v := false
	if VERBOSE == "1" {
		v = true
	}
	adapter := bluetooth.DefaultAdapter
	must("enable BLE stack", adapter.Enable())
	fmt.Fprintln(os.Stderr, "scanning...")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		mds := device.AdvertisementPayload.ManufacturerData()
		for _, md := range mds {
			if md.CompanyID != 0x0499 {
				continue
			}
			if len(md.Data) < 1 {
				continue
			}
			if v {
				fmt.Fprintln(os.Stderr,
					"[*]",
					time.Now().Format(time.DateTime),
					device.Address.String(),
					device.RSSI,
					device.LocalName())
			}
			switch md.Data[0] {
			case data.VERSION_RAWV2, data.VERSION_CUTRAWV2:
				r2, err := rawv2.Parse(md.Data)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Decoding Raw V2 packet failed:", err)
					continue
				}
				out, err := json.Marshal(&r2)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Packet marshalling failed:", err)
					continue
				}
				fmt.Println(string(out))
			default:
				fmt.Fprintf(os.Stderr, "v=%X\n", md.Data[0])
				continue
			}
		}
	})
	must("start scan", err)
}
func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
