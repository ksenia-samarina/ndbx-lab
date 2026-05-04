package events

import (
	"context"
	"samarina/ndbx/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetEvents(ctx context.Context, filters model.EventFilter) ([]model.Event, error) {
	filter := bson.M{}

	if filters.Title != "" {
		filter["title"] = bson.M{"$regex": filters.Title, "$options": "i"}
	}

	if filters.ID != "" {
		if objID, err := primitive.ObjectIDFromHex(filters.ID); err == nil {
			filter["_id"] = objID
		}
	}

	if filters.Category != "" {
		filter["category"] = filters.Category
	}

	if filters.City != "" {
		filter["location.city"] = filters.City
	}

	if filters.PriceFrom >= 0 || filters.PriceTo >= 0 {
		priceFilter := bson.M{}
		if filters.PriceFrom >= 0 {
			priceFilter["$gte"] = uint(filters.PriceFrom)
		}
		if filters.PriceTo >= 0 {
			priceFilter["$lte"] = uint(filters.PriceTo)
		}
		filter["price"] = priceFilter
	}

	if !filters.DateFrom.IsZero() || !filters.DateTo.IsZero() {
		dateFilter := bson.M{}
		if !filters.DateFrom.IsZero() {
			dateFilter["$gte"] = filters.DateFrom.Format(time.RFC3339)
		}
		if !filters.DateTo.IsZero() {
			endOfDay := time.Date(filters.DateTo.Year(), filters.DateTo.Month(), filters.DateTo.Day(), 23, 59, 59, 0, filters.DateTo.Location())
			dateFilter["$lte"] = endOfDay.Format(time.RFC3339)
		}
		filter["started_at"] = dateFilter
	}
	if filters.User != "" {
		filter["created_by"] = filters.User
	}
	findOptions := options.Find()
	if filters.Limit > 0 {
		findOptions.SetLimit(filters.Limit)
	}
	if filters.Offset > 0 {
		findOptions.SetSkip(filters.Offset)
	}
	cursor, err := s.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	events := make([]model.Event, 0)
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}
