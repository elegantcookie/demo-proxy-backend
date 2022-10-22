package db

import (
	"context"
	"log"
	"proxy_crud/internal/config"
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"testing"
)

func TestDb_Update(t *testing.T) {
	logger := logging.GetLogger("trace")
	psqlClient, err := postgresql.NewClient(context.Background(), 3, config.GetConfig().Storage)
	if err != nil {
		log.Fatalf("%v", err)
	}
	proxyStorage := NewStorage(psqlClient, &logger)
	if proxyStorage == nil {
		t.Fatal("storage is nil")
	}

	ctx := context.Background()
	t.Run("update correct fields", func(t *testing.T) {
		id := "5e7df0e5-db4a-4d2e-9605-d0c4d4df83d7"
		proxy, err := proxyStorage.FindById(ctx, id)
		if err != nil {
			t.Fatalf("failed to get record with id: %s, err: %v, create or use another one", id, err)
		}
		proxy.BLCheck = 1
		err = proxyStorage.Update(ctx, proxy)
		if err != nil {
			t.Fatalf("failed to update: %s", err)
		}
	})
}
