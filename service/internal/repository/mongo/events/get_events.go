package events

import (
	"context"
	"log"
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

	log.Printf("datefrom,dateto, %s, %s", filters.DateFrom, filters.DateTo)
	dateFrom, _ := time.Parse(time.RFC3339, filters.DateFrom)
	dateTo, _ := time.Parse(time.RFC3339, filters.DateFrom)
	if !dateFrom.IsZero() || !dateTo.IsZero() {
		if !dateFrom.IsZero() {
			filter["started_at"] = bson.M{"$gte": filters.DateFrom}
		}
		if !dateTo.IsZero() {
			filters.DateTo = filters.DateFrom
		}
		if !dateTo.IsZero() {
			log.Printf("dateTo: %s", filters.DateTo)
			t, _ := time.Parse(time.RFC3339, filters.DateTo)
			t = t.Add(time.Hour * 21)
			filter["finished_at"] = bson.M{"$lt": t.Format(time.RFC3339)}
			log.Printf("dateTo t: %s", t.Format(time.RFC3339))
		}
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

	log.Printf("getEventFilters: %v, events: %v", filters, events)
	return events, nil
}
