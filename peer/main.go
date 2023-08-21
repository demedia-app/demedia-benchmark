package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/sithumonline/demedia-benchmark/models"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/sithumonline/demedia-nostr/host"
	"github.com/sithumonline/demedia-nostr/keys"
	"github.com/sithumonline/demedia-nostr/relayer/ql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BridgeService struct {
	db *gorm.DB
}

func newBridgeService(db *gorm.DB) *BridgeService {
	return &BridgeService{
		db: db,
	}
}

func (t *BridgeService) Ql(ctx context.Context, argType ql.BridgeArgs, replyType *ql.BridgeReply) error {
	call := ql.BridgeCall{}
	err := json.Unmarshal(argType.Data, &call)
	if err != nil {
		return err
	}
	log.Printf("Received a Ql call, method: %s", call.Method)
	switch call.Method {
	case "getAllItem":
		return t.getAllItem(replyType)
	default:
		log.Printf("Received a call, method: %s", call.Method)
		return errors.New("method not found")
	}
}

func (t *BridgeService) getAllItem(replyType *ql.BridgeReply) error {
	list := make([]models.Todo, 0)

	start := time.Now()
	if result := t.db.Find(&list); result.Error != nil {
		log.Printf("failed to find todos: %v", result.Error)
		return result.Error
	}
	elapsed := time.Since(start)
	list = append(list, models.Todo{Id: "time", Title: elapsed.String()})

	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	replyType.Data = b
	return nil
}

func database(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	log.Print("database connected")
	return db
}

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

	db := database("postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil")
	bridgeService := newBridgeService(db)
	if err := rpcHost.Register(bridgeService); err != nil {
		log.Fatalf("failed to register rpc server: %v", err)
	}

	select {}
}
