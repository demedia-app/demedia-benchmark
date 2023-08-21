package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sithumonline/demedia-benchmark/models"

	"github.com/olekukonko/tablewriter"
	"github.com/sithumonline/demedia-nostr/host"
	"github.com/sithumonline/demedia-nostr/keys"
	"github.com/sithumonline/demedia-nostr/port"
	"github.com/sithumonline/demedia-nostr/relayer/ql"
)

func main() {
	var benchmarkData [][]string

	b := make([]byte, 32)
	rand.Read(b)
	_, privKey, _, _, err := keys.GetKeys(hex.EncodeToString(b))
	if err != nil {
		log.Fatalf("failed to get priv key for libp2p: %v", err)
	}

	add := host.GetAdd(fmt.Sprintf("%d", port.GetTargetAddressPort()), "1")
	h, err := host.GetHost(*privKey, add)
	if err != nil {
		log.Fatalf("failed to get host: %v", err)
	}

	start := time.Now()
	reply, sandErr := ql.QlCall(h,
		context.Background(),
		nil,
		"/ip4/192.168.1.2/tcp/10880/p2p/16Uiu2HAm11tBBtFMubGtVWty12oYHzq58k7p3ZfdPhe24qgKVgX7",
		"BridgeService",
		"Ql",
		"getAllItem",
		nil)
	if sandErr != nil {
		log.Panicf("error: failed to fetch: %s", sandErr.Error())
	}
	elapsed := time.Since(start)

	var d []models.Todo
	err = json.Unmarshal(reply.Data, &d)
	if err != nil {
		log.Fatalf("failed to unmarshal reply data: %v", err)
	}

	benchmarkData = append(benchmarkData, []string{"Ql", d[len(d)-1].Title, elapsed.String()})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Load Time", "Response Time"})
	table.AppendBulk(benchmarkData)
	table.Render()
}
