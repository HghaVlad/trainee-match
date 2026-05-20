package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
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
