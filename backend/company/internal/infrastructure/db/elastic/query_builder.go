package elastic

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

func vacancyPublicReqToQuery(
	requirements *list.Requirements,
	order list.Order,
	cursor any,
	limit int,
) (map[string]any, error) {
	sort, err := getPubVacSort(order)
	if err != nil {
		return nil, err
	}

	searchAfter, err := getPubVacSearchAfter(order, cursor)
	if err != nil {
		return nil, err
	}

	var must []any
	var filter []any

	if requirements != nil {
		must = getPubVacMust(*requirements)
		filter = getPubVacFilter(*requirements)
	}

	query := map[string]any{
		"size": limit,
		"query": map[string]any{
			"bool": map[string]any{
				"must":   must,
				"filter": filter,
			},
		},
		"sort": sort,
	}

	if len(searchAfter) > 0 {
		query["search_after"] = searchAfter
	}

	return query, nil
}

func getPubVacMust(req list.Requirements) []any {
	var must []any

	if req.Query != nil {
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query": *req.Query,
				"fields": []string{
					"title^3",
					"company_name^2",
					"description",
				},
			},
		})
	}

	return must
}

func getPubVacFilter(req list.Requirements) []any {
	var filter []any

	filter = append(filter, map[string]any{
		"term": map[string]any{
			"status": "published",
		},
	})

	filter = addWorkFormats(filter, req.WorkFormat)
	filter = addCompanies(filter, req.Companies)

	filter = addBoolTerm(filter, "is_paid", req.IsPaid)
	filter = addBoolTerm(filter, "internship_to_offer", req.InternshipToOffer)
	filter = addBoolTerm(filter, "flexible_schedule", req.FlexibleSchedule)

	filter = addSalaryIntersection(filter, req.Salary)
	filter = addHoursIntersection(filter, req.HoursPerWeek)
	filter = addDurationIntersection(filter, req.Duration)

	return filter
}

func addWorkFormats(filter []any, wfs *[]vacancy.WorkFormat) []any {
	if wfs != nil {
		values := make([]string, 0)

		for _, wf := range *wfs {
			values = append(values, string(wf))
		}

		filter = append(filter, map[string]any{
			"terms": map[string]any{
				"work_format": values,
			},
		})
	}

	return filter
}

func addCompanies(filter []any, comps *[]uuid.UUID) []any {
	if comps != nil {
		values := make([]string, 0)

		for _, id := range *comps {
			values = append(values, id.String())
		}

		filter = append(filter, map[string]any{
			"terms": map[string]any{
				"company_id": values,
			},
		})
	}

	return filter
}

func addCities(filter []any, cities *[]string) []any {
	if cities != nil {
		filter = append(filter, map[string]any{
			"terms": map[string]any{
				"city": *cities,
			},
		})
	}

	return filter
}

func addBoolTerm(filter []any, field string, value *bool) []any {
	if value != nil {
		filter = append(filter, map[string]any{
			"term": map[string]any{
				field: *value,
			},
		})
	}

	return filter
}

func addSalaryIntersection(filter []any, rng *list.RangeInt) []any {
	if rng == nil {
		return filter
	}

	if rng.Min != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"salary_to": map[string]any{
					"gte": *rng.Min,
				},
			},
		})
	}

	if rng.Max != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"salary_from": map[string]any{
					"lte": *rng.Max,
				},
			},
		})
	}

	return filter
}

func addHoursIntersection(filter []any, rng *list.RangeInt) []any {
	if rng == nil {
		return filter
	}

	if rng.Min != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"hours_per_week_to": map[string]any{
					"gte": *rng.Min,
				},
			},
		})
	}

	if rng.Max != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"hours_per_week_from": map[string]any{
					"lte": *rng.Max,
				},
			},
		})
	}

	return filter
}

func addDurationIntersection(filter []any, rng *list.RangeInt) []any {
	if rng == nil {
		return filter
	}

	if rng.Min != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"duration_to_days": map[string]any{
					"gte": *rng.Min,
				},
			},
		})
	}

	if rng.Max != nil {
		filter = append(filter, map[string]any{
			"range": map[string]any{
				"duration_from_days": map[string]any{
					"lte": *rng.Max,
				},
			},
		})
	}

	return filter
}

func getPubVacSort(order list.Order) ([]any, error) {
	var sort []any

	switch order {
	case list.OrderPublishedAtDesc:
		sort = []any{
			map[string]any{"published_at": map[string]any{"order": "desc"}},
			map[string]any{"id": map[string]any{"order": "desc"}},
		}
	case list.OrderSalaryAsc:
		sort = []any{
			map[string]any{"salary_from": map[string]any{
				"order":   "asc",
				"missing": "_last",
			},
			},
			map[string]any{
				"salary_to": map[string]any{
					"order":   "asc",
					"missing": "_last",
				},
			},
			map[string]any{"id": map[string]any{"order": "asc"}},
		}
	default:
		return nil, common.ErrUnsupportedListOrder
	}

	return sort, nil
}

func getPubVacSearchAfter(order list.Order, cursor any) ([]any, error) {
	var searchAfter []any

	if cursor == nil || reflect.ValueOf(cursor).IsNil() {
		return searchAfter, nil
	}

	switch order {
	case list.OrderPublishedAtDesc:
		curs, ok := cursor.(*list.PublishedAtCursor)
		if !ok {
			return nil, common.ErrCursorOrderMismatch
		}
		searchAfter = []any{
			curs.PublishedAt,
			curs.ID.String(),
		}

	case list.OrderSalaryAsc:
		curs, ok := cursor.(*list.SalaryCursor)
		if !ok {
			return nil, common.ErrCursorOrderMismatch
		}

		searchAfter = []any{
			curs.SalaryFrom,
			curs.SalaryTo,
			curs.ID.String(),
		}
	default:
		return nil, common.ErrUnsupportedListOrder
	}

	return searchAfter, nil
}
