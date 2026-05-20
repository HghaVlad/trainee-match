package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type VacancyRepo struct {
	es *elasticsearch.Client
}

func NewVacancyRepo(es *elasticsearch.Client) *VacancyRepo {
	return &VacancyRepo{es: es}
}

func (r *VacancyRepo) Index(ctx context.Context, vac vacancy.Vacancy, compName string) error {
	doc := vacToDoc(vac, compName)

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
	requirements *list.Requirements,
	order list.Order,
	cursor any,
	limit int,
) ([]list.VacancySummary, error) {
	query, err := vacancyPublicReqToQuery(requirements, order, cursor, limit)
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

	sums, err := searchRespToVacSums(res)
	if err != nil {
		return nil, fmt.Errorf("elastic search public vacancy: decode resp: %w", err)
	}

	return sums, nil
}

func vacToDoc(vac vacancy.Vacancy, compName string) *VacancyDocument {
	return &VacancyDocument{
		ID:                vac.ID.String(),
		CompanyID:         vac.CompanyID.String(),
		CompanyName:       compName,
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

type searchResponse struct {
	Hits struct {
		Hits []struct {
			Source VacancyDocument `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func searchRespToVacSums(res *esapi.Response) ([]list.VacancySummary, error) {
	var searchResp searchResponse

	err := json.NewDecoder(res.Body).Decode(&searchResp)
	if err != nil {
		return nil, err
	}

	vacancies := make([]list.VacancySummary, 0)

	for _, hit := range searchResp.Hits.Hits {
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
		} else {
			pub = time.Now()
		}

		vacancies = append(vacancies, list.VacancySummary{
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

	return vacancies, nil
}
