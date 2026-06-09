package validator

import (
	"net/url"
	"samarina/ndbx/internal/model"
	"strconv"
	"time"
)

const DateLayout = "20060102"

func (d *Domain) ValidateParams(query url.Values) (model.EventFilter, error) {
	limitStr := query.Get("limit")
	limit := int64(-1)
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || l < 0 {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "limit"}
		}
		limit = l
	}

	offsetStr := query.Get("offset")
	offset := int64(0)
	if offsetStr != "" {
		off, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil || off < 0 {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "offset"}
		}
		offset = off
	}

	category := query.Get("category")
	validCategories := map[string]struct{}{"": {}, "meetup": {}, "concert": {}, "exhibition": {}, "party": {}, "other": {}}
	if _, exists := validCategories[category]; !exists {
		return model.EventFilter{}, &ErrInvalidFieldName{Field: "category"}
	}

	var priceFrom int64 = -1
	if pFrom := query.Get("price_from"); pFrom != "" {
		val, err := strconv.ParseInt(pFrom, 10, 64)
		if err != nil || val < 0 {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "price_from"}
		}
		priceFrom = val
	}

	var priceTo int64 = -1
	if pTo := query.Get("price_to"); pTo != "" {
		val, err := strconv.ParseInt(pTo, 10, 64)
		if err != nil || val < 0 {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "price_to"}
		}
		priceTo = val
	}

	var dateFrom string
	if dFrom := query.Get("date_from"); dFrom != "" {
		t, err := time.Parse(DateLayout, dFrom)
		if err != nil {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "date_from"}
		}
		dateFrom = t.Format(time.RFC3339)
	}

	var dateTo string
	if dTo := query.Get("date_to"); dTo != "" {
		t, err := time.Parse(DateLayout, dTo)
		if err != nil {
			return model.EventFilter{}, &ErrInvalidFieldName{Field: "date_to"}
		}
		dateTo = t.Format(time.RFC3339)
	}

	filter := model.EventFilter{
		ID:        query.Get("id"),
		Title:     query.Get("title"),
		Category:  query.Get("category"),
		PriceFrom: priceFrom,
		PriceTo:   priceTo,
		City:      query.Get("city"),
		DateFrom:  dateFrom,
		DateTo:    dateTo,
		Offset:    offset,
		Limit:     limit,
	}
	return filter, nil
}
