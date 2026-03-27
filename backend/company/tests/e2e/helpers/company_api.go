package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"
	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/company/internal/delivery/http/dto"
	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

type CompanyAPI struct {
	baseURL    string
	authClient *http.Client
}

func NewCompanyAPI(baseURL string, authClient *http.Client) *CompanyAPI {
	syncAuthCookies(baseURL, authClient)

	return &CompanyAPI{
		baseURL:    strings.TrimRight(baseURL, "/"),
		authClient: authClient,
	}
}

func Ptr[T any](v T) *T {
	return &v
}

func syncAuthCookies(baseURL string, authClient *http.Client) {
	if authClient == nil || authClient.Jar == nil {
		return
	}

	authURL, err := url.Parse(authServiceBaseUrl)
	if err != nil {
		return
	}

	targetURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	cookies := authClient.Jar.Cookies(authURL)
	if len(cookies) == 0 {
		return
	}

	authClient.Jar.SetCookies(targetURL, cookies)
}

type CreateCompanyParams struct {
	Name        string
	Description *string
	Website     *string
}

type CreateVacancyParams struct {
	Title       string
	Description string

	WorkFormat vacancy.WorkFormat
	City       *string

	DurationFromDays *int
	DurationToDays   *int

	EmploymentType   vacancy.EmploymentType
	HoursPerWeekFrom *int
	HoursPerWeekTo   *int

	FlexibleSchedule bool

	IsPaid     bool
	SalaryFrom *int
	SalaryTo   *int

	InternshipToOffer bool
}

type RangeIntFilter struct {
	Min *int
	Max *int
}

type ListVacanciesParams struct {
	Order  list.Order
	Cursor *string
	Limit  *int

	Salary       *RangeIntFilter
	HoursPerWeek *RangeIntFilter
	Duration     *RangeIntFilter

	WorkFormat []vacancy.WorkFormat
	City       []string
	CompanyIDs []uuid.UUID

	IsPaid            *bool
	InternshipToOffer *bool
	FlexibleSchedule  *bool
}

func (api *CompanyAPI) CreateCompany(t testing.TB, params CreateCompanyParams) dto.CompanyCreatedResponse {
	t.Helper()

	reqBody := dto.CompanyCreateRequest{
		Name:        params.Name,
		Description: params.Description,
		Website:     params.Website,
	}

	return doJSON[dto.CompanyCreatedResponse](
		t,
		api.authClient,
		http.MethodPost,
		api.url("/api/v1/companies/"),
		http.StatusCreated,
		reqBody,
	)
}

func (api *CompanyAPI) GetCompany(t testing.TB, companyID uuid.UUID) dto.CompanyResponse {
	t.Helper()

	return doJSON[dto.CompanyResponse](
		t,
		api.authClient,
		http.MethodGet,
		api.url("/api/v1/companies/%s/", companyID),
		http.StatusOK,
		nil,
	)
}

func (api *CompanyAPI) CreateVacancy(
	t testing.TB,
	companyID uuid.UUID,
	params CreateVacancyParams,
) dto.VacancyCreatedResponse {
	t.Helper()

	employmentType := string(params.EmploymentType)

	reqBody := dto.VacancyCreateRequest{
		Title:             params.Title,
		Description:       params.Description,
		WorkFormat:        string(params.WorkFormat),
		City:              params.City,
		DurationFromDays:  params.DurationFromDays,
		DurationToDays:    params.DurationToDays,
		EmploymentType:    &employmentType,
		HoursPerWeekFrom:  params.HoursPerWeekFrom,
		HoursPerWeekTo:    params.HoursPerWeekTo,
		FlexibleSchedule:  params.FlexibleSchedule,
		IsPaid:            params.IsPaid,
		SalaryFrom:        params.SalaryFrom,
		SalaryTo:          params.SalaryTo,
		InternshipToOffer: params.InternshipToOffer,
	}

	return doJSON[dto.VacancyCreatedResponse](
		t,
		api.authClient,
		http.MethodPost,
		api.url("/api/v1/companies/%s/vacancies/", companyID),
		http.StatusCreated,
		reqBody,
	)
}

func (api *CompanyAPI) GetVacancy(t testing.TB, companyID, vacancyID uuid.UUID) dto.VacancyFullResponse {
	t.Helper()

	return doJSON[dto.VacancyFullResponse](
		t,
		api.authClient,
		http.MethodGet,
		api.url("/api/v1/companies/%s/vacancies/%s/", companyID, vacancyID),
		http.StatusOK,
		nil,
	)
}

func (api *CompanyAPI) GetPublishedVacancy(t testing.TB, vacancyID uuid.UUID) dto.VacancyPublicResponse {
	t.Helper()

	return doJSON[dto.VacancyPublicResponse](
		t,
		api.authClient,
		http.MethodGet,
		api.url("/api/v1/vacancies/%s", vacancyID),
		http.StatusOK,
		nil,
	)
}

func (api *CompanyAPI) ListVacancies(t testing.TB, params ListVacanciesParams) dto.VacancyListResponse {
	t.Helper()

	query := make(url.Values)

	if params.Order != "" {
		query.Set("order", string(params.Order))
	}
	if params.Cursor != nil {
		query.Set("cursor", *params.Cursor)
	}
	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}

	addRangeFilter(query, "salary_min", "salary_max", params.Salary)
	addRangeFilter(query, "hours_min", "hours_max", params.HoursPerWeek)
	addRangeFilter(query, "duration_min", "duration_max", params.Duration)

	for _, workFormat := range params.WorkFormat {
		query.Add("work_format", string(workFormat))
	}
	for _, city := range params.City {
		query.Add("city", city)
	}
	for _, companyID := range params.CompanyIDs {
		query.Add("company_id", companyID.String())
	}

	addBoolFilter(query, "is_paid", params.IsPaid)
	addBoolFilter(query, "internship_to_offer", params.InternshipToOffer)
	addBoolFilter(query, "flexible_schedule", params.FlexibleSchedule)

	listURL := api.url("/api/v1/vacancies")
	if encoded := query.Encode(); encoded != "" {
		listURL += "?" + encoded
	}

	return doJSON[dto.VacancyListResponse](t, api.authClient, http.MethodGet, listURL, http.StatusOK, nil)
}

func (api *CompanyAPI) RequirePublishedVacancyNotFound(t testing.TB, vacancyID uuid.UUID) {
	t.Helper()

	doNoContent(t, api.authClient, http.MethodGet, api.url("/api/v1/vacancies/%s", vacancyID), http.StatusNotFound)
}

func (api *CompanyAPI) PublishVacancy(t testing.TB, companyID, vacancyID uuid.UUID) {
	t.Helper()

	doNoContent(t, api.authClient, http.MethodPost, api.url(
		"/api/v1/companies/%s/vacancies/%s/publish",
		companyID,
		vacancyID,
	), http.StatusNoContent)
}

func (api *CompanyAPI) ArchiveVacancy(t testing.TB, companyID, vacancyID uuid.UUID) {
	t.Helper()

	doNoContent(t, api.authClient, http.MethodPost, api.url(
		"/api/v1/companies/%s/vacancies/%s/archive",
		companyID,
		vacancyID,
	), http.StatusNoContent)
}

func (api *CompanyAPI) url(pattern string, args ...any) string {
	return api.baseURL + fmt.Sprintf(pattern, args...)
}

func addRangeFilter(query url.Values, minKey, maxKey string, filter *RangeIntFilter) {
	if filter == nil {
		return
	}

	if filter.Min != nil {
		query.Set(minKey, strconv.Itoa(*filter.Min))
	}
	if filter.Max != nil {
		query.Set(maxKey, strconv.Itoa(*filter.Max))
	}
}

func addBoolFilter(query url.Values, key string, value *bool) {
	if value == nil {
		return
	}

	query.Set(key, strconv.FormatBool(*value))
}

func doJSON[T any](t testing.TB, client *http.Client, method, url string, expectedStatus int, requestBody any) T {
	t.Helper()

	resp := doRequest(t, client, method, url, expectedStatus, requestBody)
	defer func() {
		_ = resp.Body.Close()
	}()

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response from %s %s: %v", method, url, err)
	}

	return result
}

func doNoContent(t testing.TB, client *http.Client, method, url string, expectedStatus int) {
	t.Helper()

	resp := doRequest(t, client, method, url, expectedStatus, nil)
	defer func() {
		_ = resp.Body.Close()
	}()
}

func doRequest(
	t testing.TB,
	client *http.Client,
	method, url string,
	expectedStatus int,
	requestBody any,
) *http.Response {
	t.Helper()

	var bodyReader io.Reader = http.NoBody
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("marshal request for %s %s: %v", method, url, err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("build request %s %s: %v", method, url, err)
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request %s %s failed: %v", method, url, err)
	}

	if resp.StatusCode != expectedStatus {
		defer func() {
			_ = resp.Body.Close()
		}()

		var errResp responds.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			t.Fatalf(
				"unexpected status for %s %s: got %d, want %d, error: %s",
				method,
				url,
				resp.StatusCode,
				expectedStatus,
				errResp.Error,
			)
		}

		rawBody, _ := io.ReadAll(resp.Body)
		t.Fatalf(
			"unexpected status for %s %s: got %d, want %d, body: %s",
			method,
			url,
			resp.StatusCode,
			expectedStatus,
			string(rawBody),
		)
	}

	return resp
}
