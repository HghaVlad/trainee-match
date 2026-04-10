package e2e_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
	"github.com/HghaVlad/trainee-match/backend/company/tests/e2e/helpers"
)

type publishedVacancySeed struct {
	CompanyID uuid.UUID
	Title     string
	Params    helpers.CreateVacancyParams
}

func Test_Vacancy_List(t *testing.T) {
	api := helpers.NewCompanyAPI(baseURL, authServiceBaseURL, AuthClient)

	compA := api.CreateCompany(t, helpers.CreateCompanyParams{
		Name: "list-comp-a-" + uuid.NewString(),
	}).ID
	compB := api.CreateCompany(t, helpers.CreateCompanyParams{
		Name: "list-comp-b-" + uuid.NewString(),
	}).ID

	seeds := []publishedVacancySeed{
		{
			CompanyID: compA,
			Title:     "Go Remote Paid Junior",
			Params: helpers.CreateVacancyParams{
				Title:             "Go Remote Paid Junior",
				Description:       "remote paid junior",
				WorkFormat:        vacancy.WorkFormatRemote,
				City:              helpers.Ptr("Moscow"),
				DurationFromDays:  helpers.Ptr(60),
				DurationToDays:    helpers.Ptr(120),
				EmploymentType:    vacancy.EmploymentTypeInternship,
				HoursPerWeekFrom:  helpers.Ptr(20),
				HoursPerWeekTo:    helpers.Ptr(30),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(1000),
				SalaryTo:          helpers.Ptr(1500),
				InternshipToOffer: true,
			},
		},
		{
			CompanyID: compA,
			Title:     "Go Hybrid Mid",
			Params: helpers.CreateVacancyParams{
				Title:             "Go Hybrid Mid",
				Description:       "hybrid mid",
				WorkFormat:        vacancy.WorkFormatHybrid,
				City:              helpers.Ptr("Moscow"),
				DurationFromDays:  helpers.Ptr(90),
				DurationToDays:    helpers.Ptr(180),
				EmploymentType:    vacancy.EmploymentTypeInternship,
				HoursPerWeekFrom:  helpers.Ptr(30),
				HoursPerWeekTo:    helpers.Ptr(40),
				FlexibleSchedule:  false,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(2000),
				SalaryTo:          helpers.Ptr(2600),
				InternshipToOffer: false,
			},
		},
		{
			CompanyID: compA,
			Title:     "QA Onsite Unpaid",
			Params: helpers.CreateVacancyParams{
				Title:             "QA Onsite Unpaid",
				Description:       "onsite unpaid",
				WorkFormat:        vacancy.WorkFormatOnSite,
				City:              helpers.Ptr("Kazan"),
				DurationFromDays:  helpers.Ptr(30),
				DurationToDays:    helpers.Ptr(60),
				EmploymentType:    vacancy.EmploymentTypeInternship,
				HoursPerWeekFrom:  helpers.Ptr(15),
				HoursPerWeekTo:    helpers.Ptr(20),
				FlexibleSchedule:  false,
				IsPaid:            false,
				InternshipToOffer: true,
			},
		},
		{
			CompanyID: compA,
			Title:     "Frontend Remote Senior",
			Params: helpers.CreateVacancyParams{
				Title:             "Frontend Remote Senior",
				Description:       "remote senior",
				WorkFormat:        vacancy.WorkFormatRemote,
				City:              helpers.Ptr("Saint Petersburg"),
				DurationFromDays:  helpers.Ptr(120),
				DurationToDays:    helpers.Ptr(180),
				EmploymentType:    vacancy.EmploymentTypeFullTime,
				HoursPerWeekFrom:  helpers.Ptr(35),
				HoursPerWeekTo:    helpers.Ptr(40),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(3000),
				SalaryTo:          helpers.Ptr(3800),
				InternshipToOffer: false,
			},
		},
		{
			CompanyID: compA,
			Title:     "Data Analyst Hybrid",
			Params: helpers.CreateVacancyParams{
				Title:             "Data Analyst Hybrid",
				Description:       "analyst hybrid",
				WorkFormat:        vacancy.WorkFormatHybrid,
				City:              helpers.Ptr("Novosibirsk"),
				DurationFromDays:  helpers.Ptr(45),
				DurationToDays:    helpers.Ptr(90),
				EmploymentType:    vacancy.EmploymentTypePartTime,
				HoursPerWeekFrom:  helpers.Ptr(20),
				HoursPerWeekTo:    helpers.Ptr(25),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(1700),
				SalaryTo:          helpers.Ptr(2200),
				InternshipToOffer: true,
			},
		},
		{
			CompanyID: compA,
			Title:     "DevOps Remote Part Time",
			Params: helpers.CreateVacancyParams{
				Title:             "DevOps Remote Part Time",
				Description:       "devops remote",
				WorkFormat:        vacancy.WorkFormatRemote,
				City:              helpers.Ptr("Yerevan"),
				DurationFromDays:  helpers.Ptr(75),
				DurationToDays:    helpers.Ptr(100),
				EmploymentType:    vacancy.EmploymentTypePartTime,
				HoursPerWeekFrom:  helpers.Ptr(10),
				HoursPerWeekTo:    helpers.Ptr(20),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(2500),
				SalaryTo:          helpers.Ptr(3200),
				InternshipToOffer: false,
			},
		},
		{
			CompanyID: compB,
			Title:     "Backend Onsite Lead",
			Params: helpers.CreateVacancyParams{
				Title:             "Backend Onsite Lead",
				Description:       "backend onsite lead",
				WorkFormat:        vacancy.WorkFormatOnSite,
				City:              helpers.Ptr("Moscow"),
				DurationFromDays:  helpers.Ptr(180),
				DurationToDays:    helpers.Ptr(365),
				EmploymentType:    vacancy.EmploymentTypeFullTime,
				HoursPerWeekFrom:  helpers.Ptr(40),
				HoursPerWeekTo:    helpers.Ptr(40),
				FlexibleSchedule:  false,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(4500),
				SalaryTo:          helpers.Ptr(5500),
				InternshipToOffer: false,
			},
		},
		{
			CompanyID: compB,
			Title:     "Support Hybrid Unpaid",
			Params: helpers.CreateVacancyParams{
				Title:             "Support Hybrid Unpaid",
				Description:       "support hybrid unpaid",
				WorkFormat:        vacancy.WorkFormatHybrid,
				City:              helpers.Ptr("Samara"),
				DurationFromDays:  helpers.Ptr(30),
				DurationToDays:    helpers.Ptr(45),
				EmploymentType:    vacancy.EmploymentTypeInternship,
				HoursPerWeekFrom:  helpers.Ptr(20),
				HoursPerWeekTo:    helpers.Ptr(25),
				FlexibleSchedule:  true,
				IsPaid:            false,
				InternshipToOffer: false,
			},
		},
		{
			CompanyID: compB,
			Title:     "ML Remote Research",
			Params: helpers.CreateVacancyParams{
				Title:             "ML Remote Research",
				Description:       "ml remote research",
				WorkFormat:        vacancy.WorkFormatRemote,
				City:              helpers.Ptr("Tbilisi"),
				DurationFromDays:  helpers.Ptr(120),
				DurationToDays:    helpers.Ptr(240),
				EmploymentType:    vacancy.EmploymentTypeInternship,
				HoursPerWeekFrom:  helpers.Ptr(25),
				HoursPerWeekTo:    helpers.Ptr(35),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(2800),
				SalaryTo:          helpers.Ptr(3600),
				InternshipToOffer: true,
			},
		},
		{
			CompanyID: compB,
			Title:     "Product Remote Associate",
			Params: helpers.CreateVacancyParams{
				Title:             "Product Remote Associate",
				Description:       "product remote associate",
				WorkFormat:        vacancy.WorkFormatRemote,
				City:              helpers.Ptr("Moscow"),
				DurationFromDays:  helpers.Ptr(60),
				DurationToDays:    helpers.Ptr(90),
				EmploymentType:    vacancy.EmploymentTypePartTime,
				HoursPerWeekFrom:  helpers.Ptr(20),
				HoursPerWeekTo:    helpers.Ptr(30),
				FlexibleSchedule:  true,
				IsPaid:            true,
				SalaryFrom:        helpers.Ptr(1300),
				SalaryTo:          helpers.Ptr(1800),
				InternshipToOffer: true,
			},
		},
	}

	publishedTitles := make([]string, 0, len(seeds))
	for _, seed := range seeds {
		vacID := api.CreateVacancy(t, seed.CompanyID, seed.Params).ID
		api.PublishVacancy(t, seed.CompanyID, vacID)
		publishedTitles = append(publishedTitles, seed.Title)

		time.Sleep(5 * time.Millisecond)
	}

	draftID := api.CreateVacancy(t, compA, helpers.CreateVacancyParams{
		Title:          "Draft Hidden Vacancy",
		Description:    "draft should not be listed",
		WorkFormat:     vacancy.WorkFormatRemote,
		EmploymentType: vacancy.EmploymentTypeInternship,
	}).ID

	archivedID := api.CreateVacancy(t, compA, helpers.CreateVacancyParams{
		Title:             "Archived Hidden Vacancy",
		Description:       "archived should not be listed",
		WorkFormat:        vacancy.WorkFormatRemote,
		EmploymentType:    vacancy.EmploymentTypeInternship,
		IsPaid:            true,
		SalaryFrom:        helpers.Ptr(900),
		SalaryTo:          helpers.Ptr(1100),
		InternshipToOffer: true,
	}).ID
	api.PublishVacancy(t, compA, archivedID)
	api.ArchiveVacancy(t, compA, archivedID)

	t.Run("lists only published vacancies", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order: list.OrderPublishedAtDesc,
			Limit: helpers.Ptr(20),
		})

		require.Len(t, resp.Vacancies, len(seeds))

		gotTitles := make([]string, 0, len(resp.Vacancies))
		gotIDs := make([]uuid.UUID, 0, len(resp.Vacancies))
		for _, v := range resp.Vacancies {
			gotTitles = append(gotTitles, v.Title)
			gotIDs = append(gotIDs, v.ID)
		}

		assert.ElementsMatch(t, publishedTitles, gotTitles)
		assert.NotContains(t, gotIDs, draftID)
		assert.NotContains(t, gotIDs, archivedID)
	})

	t.Run("filters by company work format and paid flag", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order:      list.OrderPublishedAtDesc,
			CompanyIDs: []uuid.UUID{compA},
			WorkFormat: []vacancy.WorkFormat{vacancy.WorkFormatRemote},
			IsPaid:     helpers.Ptr(true),
			Limit:      helpers.Ptr(20),
		})

		require.Len(t, resp.Vacancies, 3)
		assert.Equal(t, "DevOps Remote Part Time", resp.Vacancies[0].Title)
		assert.Equal(t, "Frontend Remote Senior", resp.Vacancies[1].Title)
		assert.Equal(t, "Go Remote Paid Junior", resp.Vacancies[2].Title)
	})

	t.Run("filters by city internship offer and flexible schedule", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order:             list.OrderPublishedAtDesc,
			City:              []string{"Moscow"},
			InternshipToOffer: helpers.Ptr(true),
			FlexibleSchedule:  helpers.Ptr(true),
			Limit:             helpers.Ptr(20),
		})

		require.Len(t, resp.Vacancies, 2)
		assert.Equal(t, "Product Remote Associate", resp.Vacancies[0].Title)
		assert.Equal(t, "Go Remote Paid Junior", resp.Vacancies[1].Title)
	})

	t.Run("filters by intersecting hours and duration ranges", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order: list.OrderSalaryDesc,
			HoursPerWeek: &helpers.RangeIntFilter{
				Min: helpers.Ptr(24),
				Max: helpers.Ptr(36),
			},
			Duration: &helpers.RangeIntFilter{
				Min: helpers.Ptr(80),
				Max: helpers.Ptr(160),
			},
			IsPaid: helpers.Ptr(true),
			Limit:  helpers.Ptr(20),
		})

		require.Len(t, resp.Vacancies, 6)
		assert.Equal(t, "Frontend Remote Senior", resp.Vacancies[0].Title)
		assert.Equal(t, "ML Remote Research", resp.Vacancies[1].Title)
		assert.Equal(t, "Go Hybrid Mid", resp.Vacancies[2].Title)
		assert.Equal(t, "Data Analyst Hybrid", resp.Vacancies[3].Title)
		assert.Equal(t, "Product Remote Associate", resp.Vacancies[4].Title)
		assert.Equal(t, "Go Remote Paid Junior", resp.Vacancies[5].Title)
	})

	t.Run("orders by published at desc", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order: list.OrderPublishedAtDesc,
			Limit: helpers.Ptr(4),
		})

		require.Len(t, resp.Vacancies, 4)
		require.NotNil(t, resp.NextCursor)
		assert.Equal(t, "Product Remote Associate", resp.Vacancies[0].Title)
		assert.Equal(t, "ML Remote Research", resp.Vacancies[1].Title)
		assert.Equal(t, "Support Hybrid Unpaid", resp.Vacancies[2].Title)
		assert.Equal(t, "Backend Onsite Lead", resp.Vacancies[3].Title)
	})

	t.Run("orders by salary desc", func(t *testing.T) {
		resp := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order:  list.OrderSalaryDesc,
			IsPaid: helpers.Ptr(true),
			Limit:  helpers.Ptr(5),
		})

		require.Len(t, resp.Vacancies, 5)
		assert.Equal(t, "Backend Onsite Lead", resp.Vacancies[0].Title)
		assert.Equal(t, "Frontend Remote Senior", resp.Vacancies[1].Title)
		assert.Equal(t, "ML Remote Research", resp.Vacancies[2].Title)
		assert.Equal(t, "DevOps Remote Part Time", resp.Vacancies[3].Title)
		assert.Equal(t, "Go Hybrid Mid", resp.Vacancies[4].Title)
		require.NotNil(t, resp.NextCursor)
	})

	t.Run("paginates published at desc without overlap", func(t *testing.T) {
		firstPage := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order: list.OrderPublishedAtDesc,
			Limit: helpers.Ptr(3),
		})

		require.Len(t, firstPage.Vacancies, 3)
		require.NotNil(t, firstPage.NextCursor)
		assert.Equal(t, "Product Remote Associate", firstPage.Vacancies[0].Title)
		assert.Equal(t, "ML Remote Research", firstPage.Vacancies[1].Title)
		assert.Equal(t, "Support Hybrid Unpaid", firstPage.Vacancies[2].Title)

		secondPage := api.ListVacancies(t, helpers.ListVacanciesParams{
			Order:  list.OrderPublishedAtDesc,
			Limit:  helpers.Ptr(3),
			Cursor: firstPage.NextCursor,
		})

		require.Len(t, secondPage.Vacancies, 3)
		assert.Equal(t, "Backend Onsite Lead", secondPage.Vacancies[0].Title)
		assert.Equal(t, "DevOps Remote Part Time", secondPage.Vacancies[1].Title)
		assert.Equal(t, "Data Analyst Hybrid", secondPage.Vacancies[2].Title)

		for _, first := range firstPage.Vacancies {
			for _, second := range secondPage.Vacancies {
				assert.NotEqual(t, first.ID, second.ID)
			}
		}
	})
}
