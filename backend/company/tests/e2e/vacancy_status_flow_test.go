//go:build e2e

package e2e_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/tests/e2e/helpers"
)

func Test_Vacancy_StatusFlow(t *testing.T) {
	api := helpers.NewCompanyAPI(baseURL, authServiceBaseURL, AuthClient)

	compID := api.CreateCompany(t,
		helpers.CreateCompanyParams{Name: "1comp" + uuid.New().String()},
	).ID

	vacID1 := api.CreateVacancy(t, compID,
		helpers.CreateVacancyParams{
			Description:    "desc",
			Title:          "title",
			WorkFormat:     vacancy.WorkFormatHybrid,
			EmploymentType: vacancy.EmploymentTypeInternship,
		}).ID

	vacD := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, string(vacancy.StatusDraft), vacD.Status)

	comp0 := api.GetCompany(t, compID)
	assert.Equal(t, 0, comp0.OpenVacanciesCnt)

	api.PublishVacancy(t, compID, vacID1)

	vacP := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, vacP.Status, string(vacancy.StatusPublished))

	comp1 := api.GetCompany(t, compID)
	assert.Equal(t, 1, comp1.OpenVacanciesCnt)

	// idempotency
	api.PublishVacancy(t, compID, vacID1)

	comp11 := api.GetCompany(t, compID)
	assert.Equal(t, 1, comp11.OpenVacanciesCnt)

	vacID2 := api.CreateVacancy(t, compID,
		helpers.CreateVacancyParams{
			Description:    "desc2",
			Title:          "title2",
			WorkFormat:     vacancy.WorkFormatHybrid,
			EmploymentType: vacancy.EmploymentTypeInternship,
		}).ID

	comp111 := api.GetCompany(t, compID)
	assert.Equal(t, 1, comp111.OpenVacanciesCnt)

	api.PublishVacancy(t, compID, vacID2)

	comp2 := api.GetCompany(t, compID)
	assert.Equal(t, 2, comp2.OpenVacanciesCnt)

	api.ArchiveVacancy(t, compID, vacID1)

	vacA := api.GetVacancy(t, compID, vacID1)
	assert.Equal(t, vacA.Status, string(vacancy.StatusArchived))

	// check that can't get published
	api.RequirePublishedVacancyNotFound(t, vacID1)

	comp1Rem := api.GetCompany(t, compID)
	assert.Equal(t, 1, comp1Rem.OpenVacanciesCnt)

	// idempotency
	api.ArchiveVacancy(t, compID, vacID1)

	comp11Rem := api.GetCompany(t, compID)
	assert.Equal(t, 1, comp11Rem.OpenVacanciesCnt)

	api.ArchiveVacancy(t, compID, vacID2)
	comp0Rem := api.GetCompany(t, compID)
	assert.Equal(t, 0, comp0Rem.OpenVacanciesCnt)
}
