package middlewares

import (
	"context"
	"time"

	"github.com/turbex-backend/src/consts"
	"github.com/turbex-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getSession(db *mongo.Database, sessionToken string) (*models.Session, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session models.Session

	err := db.Collection(consts.COLLECTION_SESSIONS).FindOne(ctx, bson.M{"cookievalue": sessionToken}).Decode(&session)

	if err != nil {
		return nil, false
	}

	return &session, true

}
