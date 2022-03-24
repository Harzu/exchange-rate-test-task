package rates

import (
	"context"

	"github.com/Harzu/exchange-rate-test-task/internal/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const ratesCollection = "rates"

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, database string) *Repository {
	return &Repository{collection: client.Database(database).Collection(ratesCollection)}
}

func (r *Repository) GetPairRates(ctx context.Context, sourceSymbols, targetSymbols []string) ([]entities.Rate, error) {
	filter := bson.D{
		{
			"$and", bson.A{
				bson.D{{"SourceSymbol", bson.A{"$in", sourceSymbols}}},
				bson.D{{"TargetSymbol", bson.A{"$in", targetSymbols}}},
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var pairs []entities.Rate
	for cursor.Next(ctx) {
		var pair entities.Rate
		if err := cursor.Decode(&pair); err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	if err := cursor.Close(ctx); err != nil {
		return nil, err
	}

	return pairs, nil
}
