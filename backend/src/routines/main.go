package routines

import (
	"context"
	"log"
	"time"

	"github.com/turbex-backend/src/consts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanExpiredSessions(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filters := bson.M{
		"expirationdate": bson.M{
			"$lt": time.Now().UTC().Format(consts.DATE_FORMAT),
		},
	}

	_, err := db.Collection(consts.COLLECTION_SESSIONS).DeleteMany(ctx, filters)

	if err != nil {
		log.Printf("[Session cleaning] %s", err)
	}
}
