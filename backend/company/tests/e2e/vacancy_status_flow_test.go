package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/value_types"
	"github.com/HghaVlad/trainee-match/backend/company/tests/e2e/helpers"
)

func Test_Vacancy_StatusFlow(t *testing.T) {
	api := helpers.NewCompanyAPI(baseURL, AuthClient)

	compID := api.CreateCompany(t,
		helpers.CreateCompanyParams{Name: "1comp" + uuid.New().String()},
	).ID

	vacID1 := api.CreateVacancy(t, compID,
		helpers.CreateVacancyParams{
			Description:    "desc",
			Title:          "title",
			WorkFormat:     value_types.WorkFormatHybrid,
			EmploymentType: value_types.EmploymentTypeInternship,
		}).ID

	vacD := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, vacD.Status, string(value_types.VacancyStatusDraft))

	comp0 := api.GetCompany(t, compID)
	assert.Equal(t, comp0.OpenVacanciesCnt, 0)

	api.PublishVacancy(t, compID, vacID1)

	vacP := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, vacP.Status, string(value_types.VacancyStatusPublished))

	comp1 := api.GetCompany(t, compID)
	assert.Equal(t, comp1.OpenVacanciesCnt, 1)

	// idempotency
	api.PublishVacancy(t, compID, vacID1)

	comp11 := api.GetCompany(t, compID)
	assert.Equal(t, comp11.OpenVacanciesCnt, 1)

	vacID2 := api.CreateVacancy(t, compID,
		helpers.CreateVacancyParams{
			Description:    "desc2",
			Title:          "title2",
			WorkFormat:     value_types.WorkFormatHybrid,
			EmploymentType: value_types.EmploymentTypeInternship,
		}).ID

	comp111 := api.GetCompany(t, compID)
	assert.Equal(t, comp111.OpenVacanciesCnt, 1)

	api.PublishVacancy(t, compID, vacID2)

	comp2 := api.GetCompany(t, compID)
	assert.Equal(t, comp2.OpenVacanciesCnt, 2)

	api.ArchiveVacancy(t, compID, vacID1)

	vacA := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, vacA.Status, string(value_types.VacancyStatusArchived))

	// check that can't get published
	api.RequirePublishedVacancyNotFound(t, vacID1)

	comp1Rem := api.GetCompany(t, compID)
	assert.Equal(t, comp1Rem.OpenVacanciesCnt, 1)

	// idempotency
	api.ArchiveVacancy(t, compID, vacID1)

	comp11Rem := api.GetCompany(t, compID)
	assert.Equal(t, comp11Rem.OpenVacanciesCnt, 1)

	api.ArchiveVacancy(t, compID, vacID2)
	comp0Rem := api.GetCompany(t, compID)
	assert.Equal(t, comp0Rem.OpenVacanciesCnt, 0)
}
