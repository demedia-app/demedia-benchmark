package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sithumonline/demedia-benchmark/database"
	"github.com/sithumonline/demedia-benchmark/service"
	"github.com/sithumonline/demedia-benchmark/util"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/sithumonline/demedia-nostr/host"
	"github.com/sithumonline/demedia-nostr/keys"
	"github.com/sithumonline/demedia-nostr/relayer/ql"
)

func main() {
	_, privKey, _, _, err := keys.GetKeys("fad5c8833b841a0b1ed4c323dbad0f11a83a49cad6b3fe8d5234ac83d38b6a19")
	if err != nil {
		log.Fatalf("failed to get priv key for libp2p: %v", err)
	}

	add := host.GetAdd("10880", "1")
	h, err := host.GetHost(*privKey, add)
	if err != nil {
		log.Fatalf("failed to get host: %v", err)
	}

	log.Printf("Peer: listening on %s\n", host.GetMultiAddr(h))

	rpcHost := gorpc.NewServer(h, "/p2p/1.0.0")

	db := database.Database(util.EnvOrDefault("DATABASE_URL", "postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil"))
	bridgeService := service.NewBridgeService(db)
	if err := rpcHost.Register(bridgeService); err != nil {
		log.Fatalf("failed to register rpc server: %v", err)
	}

	http.HandleFunc("/getAllItem", func(w http.ResponseWriter, r *http.Request) {
		data := ql.BridgeReply{}
		err := bridgeService.GetAllItem(&data)
		if err != nil {
			fmt.Fprintf(w, "error: %v", err)
		}
		w.Write(data.Data)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", util.EnvOrDefault("PORT", "8080")), nil))
}
