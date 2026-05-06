package helpers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listbycomp"
)

func ListVacRequestFromQuery(r *http.Request) (*list.Request, error) {
	q := r.URL.Query()

	order, err := parseVacListOrderQuery(r)
	if err != nil {
		return nil, err
	}

	req := &list.Request{
		Limit:         ParseLimit(r, "limit", 20),
		Order:         order,
		EncodedCursor: q.Get("cursor"),
		Requirements:  &list.Requirements{},
	}

	applyRanges(q, req)
	applyBoolFilters(q, req)
	applyWorkFormat(q, req)
	applyCompanies(q, req)
	applyCity(q, req)

	return req, nil
}

func applyRanges(q url.Values, req *list.Request) {
	req.Requirements.Salary = parseRangeInt(q, "salary_min", "salary_max")
	req.Requirements.HoursPerWeek = parseRangeInt(q, "hours_min", "hours_max")
	req.Requirements.Duration = parseRangeInt(q, "duration_min", "duration_max")
}

func applyBoolFilters(q url.Values, req *list.Request) {
	parseBoolToPtr(q, "is_paid", &req.Requirements.IsPaid)
	parseBoolToPtr(q, "internship_to_offer", &req.Requirements.InternshipToOffer)
	parseBoolToPtr(q, "flexible_schedule", &req.Requirements.FlexibleSchedule)
}

func parseBoolToPtr(q url.Values, key string, target **bool) {
	if str := q.Get(key); str != "" {
		if val, err := strconv.ParseBool(str); err == nil {
			*target = &val
		}
	}
}

func applyWorkFormat(q url.Values, req *list.Request) {
	values, ok := q["work_format"]
	if !ok || len(values) == 0 {
		return
	}

	var wfs []vacancy.WorkFormat
	for _, str := range values {
		wf := vacancy.WorkFormat(str)
		if wf.IsValid() {
			wfs = append(wfs, wf)
		}
	}

	if len(wfs) > 0 {
		req.Requirements.WorkFormat = &wfs
	}
}

func applyCompanies(q url.Values, req *list.Request) {
	values, ok := q["company_id"]
	if !ok || len(values) == 0 {
		return
	}

	ids := make([]uuid.UUID, 0, len(values))
	for _, str := range values {
		if id, err := uuid.Parse(str); err == nil {
			ids = append(ids, id)
		}
	}

	if len(ids) > 0 {
		req.Requirements.Companies = &ids
	}
}

func applyCity(q url.Values, req *list.Request) {
	if cities, ok := q["city"]; ok && len(cities) > 0 {
		req.Requirements.City = &cities
	}
}

func parseVacListOrderQuery(r *http.Request) (list.Order, error) {
	str := r.URL.Query().Get("order")
	if str == "" {
		return list.OrderPublishedAtDesc, nil
	}

	ord := list.Order(strings.Trim(str, " "))

	switch ord {
	case list.OrderPublishedAtDesc,
		list.OrderSalaryDesc,
		list.OrderSalaryAsc:
		return ord, nil
	default:
		return "", common.ErrUnsupportedListOrder
	}
}

func ParseVacByCompListOrderQuery(r *http.Request) listbycomp.Order {
	str := r.URL.Query().Get("order")
	ord := listbycomp.Order(strings.Trim(str, " "))

	switch ord {
	case listbycomp.OrderCreatedAtDesc:
		return ord
	default:
		return listbycomp.OrderCreatedAtDesc
	}
}

func ListVacByCompRequestFromQuery(r *http.Request, compID uuid.UUID) (*listbycomp.Request, error) {
	q := r.URL.Query()

	order := ParseVacByCompListOrderQuery(r)
	req := &listbycomp.Request{
		CompID:        compID,
		Limit:         ParseLimit(r, "limit", 20),
		Order:         order,
		EncodedCursor: q.Get("cursor"),
		Requirements:  &list.Requirements{},
	}

	req.Requirements.Salary = parseRangeInt(q, "salary_min", "salary_max")
	req.Requirements.HoursPerWeek = parseRangeInt(q, "hours_min", "hours_max")
	req.Requirements.Duration = parseRangeInt(q, "duration_min", "duration_max")

	parseBoolToPtr(q, "is_paid", &req.Requirements.IsPaid)
	parseBoolToPtr(q, "internship_to_offer", &req.Requirements.InternshipToOffer)
	parseBoolToPtr(q, "flexible_schedule", &req.Requirements.FlexibleSchedule)

	applyRequirementsWorkFormat(q, req.Requirements)
	applyRequirementsCity(q, req.Requirements)

	if statusStr := strings.TrimSpace(q.Get("status")); statusStr != "" {
		status := vacancy.Status(statusStr)
		req.Status = &status
	}

	return req, nil
}

func applyRequirementsWorkFormat(q url.Values, req *list.Requirements) {
	wrapper := &list.Request{Requirements: req}
	applyWorkFormat(q, wrapper)
}

func applyRequirementsCity(q url.Values, req *list.Requirements) {
	wrapper := &list.Request{Requirements: req}
	applyCity(q, wrapper)
}
