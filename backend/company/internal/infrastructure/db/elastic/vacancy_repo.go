package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listcompsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

type VacancyRepo struct {
	es *elasticsearch.Client
}

func NewVacancyRepo(es *elasticsearch.Client) *VacancyRepo {
	return &VacancyRepo{es: es}
}

func (r *VacancyRepo) Index(ctx context.Context, vac views.VacancySearch) error {
	doc := vacToDoc(vac)

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal vacancy doc to json: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      vacancyIndex,
		DocumentID: doc.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.es)
	if err != nil {
		return fmt.Errorf("elasticsearch index vacancy: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		return fmt.Errorf("elasticsearch index error: %s", res.Status())
	}

	return nil
}

func (r *VacancyRepo) ListPublishedSummaries(
	ctx context.Context,
	requirements *listsearch.Requirements,
	order listsearch.Order,
	cursor any,
	limit int,
) (*listsearch.SearchResult, error) {
	query, err := vacancyPublicReqToQuery(requirements, order, cursor, limit+1)
	if err != nil {
		return nil, fmt.Errorf("elastic search public vacancy: %w", err)
	}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return nil, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex(vacancyIndex),
		r.es.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, fmt.Errorf("elastic search public vacancy: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	result, err := searchRespToVacSums(res, order, limit)
	if err != nil {
		return nil, fmt.Errorf("elastic search public vacancy: decode resp: %w", err)
	}

	return result, nil
}

func (r *VacancyRepo) ListByCompanySummaries(
	ctx context.Context,
	requirements *listsearch.Requirements,
	status *vacancy.Status,
	order listcompsearch.Order,
	cursor any,
	limit int,
) (*listcompsearch.SearchResult, error) {
	query, err := vacancyListCompToQuery(requirements, status, order, cursor, limit+1)
	if err != nil {
		return nil, fmt.Errorf("elastic search company vacancy: %w", err)
	}

	var buf bytes.Buffer

	err = json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return nil, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex(vacancyIndex),
		r.es.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, fmt.Errorf("elastic search comp vacancy: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	result, err := searchRespToVacCompSums(res, order, limit)
	if err != nil {
		return nil, fmt.Errorf("elastic search comp vacancy: decode resp: %w", err)
	}

	return result, nil
}

func (r *VacancyRepo) UpdateCompanyName(ctx context.Context, compID uuid.UUID, newName string) error {
	query := map[string]any{
		"script": map[string]any{
			"source": "ctx._source.company_name = params.name",
			"params": map[string]any{
				"name": newName,
			},
		},
		"query": map[string]any{
			"term": map[string]any{
				"company_id": compID.String(),
			},
		},
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return fmt.Errorf("encode elastic update company name query: %w", err)
	}

	res, err := r.es.UpdateByQuery(
		[]string{vacancyIndex},
		r.es.UpdateByQuery.WithContext(ctx),
		r.es.UpdateByQuery.WithBody(&buf),
		r.es.UpdateByQuery.WithRefresh(true),
		r.es.UpdateByQuery.WithConflicts("proceed"),
	)
	if err != nil {
		return fmt.Errorf("elastic update company name: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("elastic update company name unexpected status %s: %s", res.Status(), string(body))
	}

	return nil
}

func (r *VacancyRepo) RemoveByCompanyID(ctx context.Context, compID uuid.UUID) error {
	query := map[string]any{
		"query": map[string]any{
			"term": map[string]any{
				"company_id": compID.String(),
			},
		},
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(query)
	if err != nil {
		return fmt.Errorf("encode elastic delete vacancies by company query: %w", err)
	}

	res, err := r.es.DeleteByQuery(
		[]string{vacancyIndex},
		&buf,
		r.es.DeleteByQuery.WithContext(ctx),
		r.es.DeleteByQuery.WithRefresh(true),
		r.es.DeleteByQuery.WithConflicts("proceed"),
	)
	if err != nil {
		return fmt.Errorf("elastic delete vacancies by company: %w", err)
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)

		return fmt.Errorf(
			"elastic delete vacancies by company unexpected status %s: %s",
			res.Status(),
			string(body),
		)
	}

	return nil
}

type searchResponse struct {
	Hits struct {
		Hits []vacHit `json:"hits"`
	} `json:"hits"`
}

type vacHit struct {
	Source VacancyDocument `json:"_source"`
	Sort   []any           `json:"sort"`
}

func searchRespToVacSums(res *esapi.Response, order listsearch.Order, limit int) (*listsearch.SearchResult, error) {
	var searchResp searchResponse
	err := json.NewDecoder(res.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	hits := searchResp.Hits.Hits

	result := &listsearch.SearchResult{
		Vacancies: make([]views.PublishedVacSummary, 0, limit+1),
	}

	if len(hits) == 0 {
		return result, nil
	}

	result.HasNext = len(hits) > limit

	for _, hit := range hits {
		doc := hit.Source

		id, err := uuid.Parse(doc.ID)
		if err != nil {
			return nil, err
		}

		compID, err := uuid.Parse(doc.CompanyID)
		if err != nil {
			return nil, err
		}

		var pub time.Time
		if doc.PublishedAt != nil {
			pub = *doc.PublishedAt
		}

		result.Vacancies = append(result.Vacancies, views.PublishedVacSummary{
			ID:             id,
			CompanyID:      compID,
			CompanyName:    doc.CompanyName,
			Title:          doc.Title,
			WorkFormat:     vacancy.WorkFormat(doc.WorkFormat),
			City:           doc.City,
			EmploymentType: vacancy.EmploymentType(doc.EmploymentType),
			IsPaid:         doc.IsPaid,
			SalaryFrom:     doc.SalaryFrom,
			SalaryTo:       doc.SalaryTo,
			PublishedAt:    pub,
		})
	}

	if !result.HasNext {
		return result, nil
	}

	result.Vacancies = result.Vacancies[:limit]
	lastVac := result.Vacancies[len(result.Vacancies)-1]
	lastHit := hits[limit-1]

	switch order {
	case listsearch.OrderRelevance:
		result.NextCursor, err = buildRelevanceCursor(lastHit.Sort)

	case listsearch.OrderPublishedAtDesc:
		result.NextCursor = &listsearch.PublishedAtCursor{
			PublishedAt: lastVac.PublishedAt,
			ID:          lastVac.ID,
		}

	case listsearch.OrderSalaryAsc, listsearch.OrderSalaryDesc:
		if lastVac.SalaryFrom == nil || lastVac.SalaryTo == nil {
			result.HasNext = false
		} else {
			result.NextCursor = &listsearch.SalaryCursor{
				SalaryFrom: *lastVac.SalaryFrom,
				SalaryTo:   *lastVac.SalaryTo,
				ID:         lastVac.ID,
			}
		}

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func searchRespToVacCompSums(
	res *esapi.Response,
	order listcompsearch.Order,
	limit int,
) (*listcompsearch.SearchResult, error) {
	var searchResp searchResponse
	err := json.NewDecoder(res.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	hits := searchResp.Hits.Hits

	result := &listcompsearch.SearchResult{
		Vacancies: make([]views.MemberVacSummary, 0, limit+1),
	}

	if len(hits) == 0 {
		return result, nil
	}

	result.HasNext = len(hits) > limit

	for _, hit := range hits {
		doc := hit.Source

		id, err := uuid.Parse(doc.ID)
		if err != nil {
			return nil, err
		}

		result.Vacancies = append(result.Vacancies, views.MemberVacSummary{
			ID:             id,
			Status:         vacancy.Status(doc.Status),
			Title:          doc.Title,
			WorkFormat:     vacancy.WorkFormat(doc.WorkFormat),
			City:           doc.City,
			EmploymentType: vacancy.EmploymentType(doc.EmploymentType),
			IsPaid:         doc.IsPaid,
			SalaryFrom:     doc.SalaryFrom,
			SalaryTo:       doc.SalaryTo,
			ModStatus:      vacancy.ModerationStatus(doc.ModerationStatus),
			CreatedAt:      doc.CreatedAt,
		})
	}

	if !result.HasNext {
		return result, nil
	}

	result.Vacancies = result.Vacancies[:limit]
	lastVac := result.Vacancies[len(result.Vacancies)-1]
	lastHit := hits[limit-1]

	switch order {
	case listcompsearch.OrderRelevance:
		result.NextCursor, err = buildVacCompRelevanceCursor(lastHit.Sort)

	case listcompsearch.OrderCreatedAtDesc:
		result.NextCursor = &listsearch.PublishedAtCursor{
			PublishedAt: lastVac.CreatedAt,
			ID:          lastVac.ID,
		}

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildRelevanceCursor(sort []any) (*listsearch.RelevanceCursor, error) {
	if len(sort) != 3 {
		return nil, errors.New("invalid sort length")
	}

	score, ok := sort[0].(float64)
	if !ok {
		return nil, errors.New("invalid score")
	}

	pubMs, ok := sort[1].(float64)
	if !ok {
		return nil, errors.New("invalid published_at")
	}

	pub := time.UnixMilli(int64(pubMs))

	idStr, ok := sort[2].(string)
	if !ok {
		return nil, errors.New("invalid id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &listsearch.RelevanceCursor{
		Relevance:   score,
		PublishedAt: pub,
		ID:          id,
	}, nil
}

func buildVacCompRelevanceCursor(sort []any) (*listcompsearch.RelevanceCursor, error) {
	if len(sort) != 3 {
		return nil, errors.New("invalid sort length")
	}

	score, ok := sort[0].(float64)
	if !ok {
		return nil, errors.New("invalid score")
	}

	crMs, ok := sort[1].(float64)
	if !ok {
		return nil, errors.New("invalid published_at")
	}

	created := time.UnixMilli(int64(crMs))

	idStr, ok := sort[2].(string)
	if !ok {
		return nil, errors.New("invalid id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return &listcompsearch.RelevanceCursor{
		Relevance: score,
		CreatedAt: created,
		ID:        id,
	}, nil
}

func vacToDoc(vac views.VacancySearch) *VacancyDocument {
	return &VacancyDocument{
		ID:                vac.ID.String(),
		CompanyID:         vac.CompanyID.String(),
		CompanyName:       vac.CompanyName,
		Title:             vac.Title,
		Description:       vac.Description,
		WorkFormat:        string(vac.WorkFormat),
		City:              vac.City,
		EmploymentType:    string(vac.EmploymentType),
		DurationFromDays:  vac.DurationFromDays,
		DurationToDays:    vac.DurationToDays,
		HoursPerWeekFrom:  vac.HoursPerWeekFrom,
		HoursPerWeekTo:    vac.HoursPerWeekTo,
		FlexibleSchedule:  vac.FlexibleSchedule,
		IsPaid:            vac.IsPaid,
		SalaryFrom:        vac.SalaryFrom,
		SalaryTo:          vac.SalaryTo,
		InternshipToOffer: vac.InternshipToOffer,
		Status:            string(vac.Status),
		ModerationStatus:  string(vac.ModerationStatus),
		PublishedAt:       vac.PublishedAt,
		CreatedAt:         vac.CreatedAt,
	}
}
