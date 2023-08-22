package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/sithumonline/demedia-benchmark/models"
	"github.com/sithumonline/demedia-nostr/relayer/ql"
	"gorm.io/gorm"
)

type BridgeService struct {
	db *gorm.DB
}

func NewBridgeService(db *gorm.DB) *BridgeService {
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
		return t.GetAllItem(replyType)
	default:
		log.Printf("Received a call, method: %s", call.Method)
		return errors.New("method not found")
	}
}

func (t *BridgeService) GetAllItem(replyType *ql.BridgeReply) error {
	list := make([]models.Todo, 0)

	if result := t.db.Find(&list); result.Error != nil {
		log.Printf("failed to find todos: %v", result.Error)
		return result.Error
	}

	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	replyType.Data = b
	return nil
}
