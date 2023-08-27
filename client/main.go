package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sithumonline/demedia-benchmark/database"
	"github.com/sithumonline/demedia-benchmark/models"
	"github.com/sithumonline/demedia-benchmark/service"
	"github.com/sithumonline/demedia-benchmark/util"

	"github.com/olekukonko/tablewriter"
	"github.com/sithumonline/demedia-nostr/host"
	"github.com/sithumonline/demedia-nostr/keys"
	"github.com/sithumonline/demedia-nostr/port"
	"github.com/sithumonline/demedia-nostr/relayer/ql"
	"github.com/spf13/viper"
)

func call(url string) time.Duration {
	rest_start := time.Now()
	_, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to fetch rest: %v", err)
	}
	rest_elapsed := time.Since(rest_start)
	return rest_elapsed
}

func directDatabaseCall() time.Duration {
	db := database.Database(util.EnvOrDefault("DATABASE_URL", "postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil"))
	bridgeService := service.NewBridgeService(db)
	data := ql.BridgeReply{}

	time_start := time.Now()
	err := bridgeService.GetAllItem(&data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	elapsed := time.Since(time_start)

	return elapsed
}

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
	_, sandErr := ql.QlCall(h,
		context.Background(),
		nil,
		util.EnvOrDefault("PEER_ADDRESS", "/ip4/192.168.1.2/tcp/10880/p2p/16Uiu2HAm11tBBtFMubGtVWty12oYHzq58k7p3ZfdPhe24qgKVgX7"),
		"BridgeService",
		"Ql",
		"getAllItem",
		nil)
	if sandErr != nil {
		log.Panicf("error: failed to fetch: %s", sandErr.Error())
	}
	elapsed := time.Since(start)

	data_elapsed := call(util.EnvOrDefault("PEER_REST_ENDPOINT", "http://localhost:8080/getAllItem"))

	benchmarkData = append(benchmarkData, []string{"REST API Call vs DeMedia API Call", data_elapsed.String(), elapsed.String()})

	database_elapsed := directDatabaseCall()
	benchmarkData = append(benchmarkData, []string{"Direct DB Fetch vs DB Fetch Through DeMedia", database_elapsed.String(), elapsed.String()})
	benchmarkData = append(benchmarkData, []string{"Direct DB Fetch vs REST API Call", database_elapsed.String(), data_elapsed.String()})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config models.RunList
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	for _, run := range config.Runs {
		rest_elapsed := call(run.Rest)
		ipfs_elapsed := call(run.Ipfs)

		benchmarkData = append(benchmarkData, []string{fmt.Sprintf("Time Elapsed To Fetch %s", run.Name), rest_elapsed.String(), ipfs_elapsed.String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Rest", "IPFS/DeMedia"})
	table.AppendBulk(benchmarkData)
	table.Render()
}
