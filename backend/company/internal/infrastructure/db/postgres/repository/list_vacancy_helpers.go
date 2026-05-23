package repository

import (
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/HghaVlad/trainee-match/backend/company/internal/domain/vacancy"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/common"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listcompsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/listsearch"
	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/views"
)

const (
	andSalaryNotNull string = " AND v.salary_from IS NOT NULL AND v.salary_to IS NOT NULL"
)

func listVacRequirementsToSQL(requirements *listsearch.Requirements) (string, []any) {
	if requirements == nil {
		return "", nil
	}

	conditions := make([]string, 0)
	args := make([]any, 0)

	applySalary(requirements, &conditions, &args)
	applyHours(requirements, &conditions, &args)
	applyDuration(requirements, &conditions, &args)
	applyWorkFormat(requirements, &conditions, &args)
	applyCompanies(requirements, &conditions, &args)
	applyCity(requirements, &conditions, &args)
	applyFlags(requirements, &conditions, &args)

	if len(conditions) == 0 {
		return "", args
	}

	return strings.Join(conditions, " AND "), args
}

func applySalary(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.Salary == nil {
		return
	}

	if r.Salary.Min != nil {
		addCondition(conds, args, "v.salary_to >= $%d", *r.Salary.Min)
	}

	if r.Salary.Max != nil {
		addCondition(conds, args, "v.salary_from <= $%d", *r.Salary.Max)
	}
}

func applyHours(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.HoursPerWeek == nil {
		return
	}

	if r.HoursPerWeek.Min != nil {
		addCondition(conds, args, "v.hours_per_week_to >= $%d", *r.HoursPerWeek.Min)
	}

	if r.HoursPerWeek.Max != nil {
		addCondition(conds, args, "v.hours_per_week_from <= $%d", *r.HoursPerWeek.Max)
	}
}

func applyDuration(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.Duration == nil {
		return
	}

	if r.Duration.Min != nil {
		addCondition(conds, args, "v.duration_to_days >= $%d", *r.Duration.Min)
	}

	if r.Duration.Max != nil {
		addCondition(conds, args, "v.duration_from_days <= $%d", *r.Duration.Max)
	}
}

func applyWorkFormat(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.WorkFormat != nil && len(*r.WorkFormat) > 0 {
		addCondition(conds, args, "v.work_format = ANY($%d)", pq.Array(*r.WorkFormat))
	}
}

func applyCompanies(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.Companies != nil && len(*r.Companies) > 0 {
		addCondition(conds, args, "v.company_id = ANY($%d)", pq.Array(*r.Companies))
	}
}

func applyCity(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.City != nil && len(*r.City) > 0 {
		addCondition(conds, args, "v.city = ANY($%d)", pq.Array(*r.City))
	}
}

func applyFlags(r *listsearch.Requirements, conds *[]string, args *[]any) {
	if r.IsPaid != nil {
		addCondition(conds, args, "v.is_paid = $%d", *r.IsPaid)
	}

	if r.InternshipToOffer != nil {
		addCondition(conds, args, "v.internship_to_offer = $%d", *r.InternshipToOffer)
	}

	if r.FlexibleSchedule != nil {
		addCondition(conds, args, "v.flexible_schedule = $%d", *r.FlexibleSchedule)
	}
}

func addCondition(conditions *[]string, args *[]any, query string, arg any) {
	*args = append(*args, arg)
	*conditions = append(*conditions, fmt.Sprintf(query, len(*args)))
}

//--------------------------
// Cursors conditions to SQL
//--------------------------

// returns SQL condition, updated slice of args
func listVacCursorToSQL(order listsearch.Order, cursor any, args []any) (string, []any) {
	switch c := cursor.(type) {
	case *listsearch.PublishedAtCursor:
		return publishedAtCursorToSQL(*c, args)
	case *listsearch.SalaryCursor:
		return salaryCursorToSQL(order, *c, args)
	}

	return "", args
}

func publishedAtCursorToSQL(cursor listsearch.PublishedAtCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"(v.published_at, v.id) < ($%d, $%d)",
		len(args)+1, len(args)+2)

	args = append(args, cursor.PublishedAt, cursor.ID)
	return condition, args
}

func salaryCursorToSQL(order listsearch.Order, cursor listsearch.SalaryCursor, args []any) (string, []any) {
	var condition string

	if order == listsearch.OrderSalaryDesc {
		condition = fmt.Sprintf(
			"(v.salary_from, v.salary_to, v.id) < ($%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3)
	} else {
		condition = fmt.Sprintf(
			"(v.salary_from, v.salary_to, v.id) > ($%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3)
	}

	args = append(args, cursor.SalaryFrom, cursor.SalaryTo, cursor.ID)
	return condition, args
}

func listByCompStatusToSQL(status *vacancy.Status, args []any) (string, []any) {
	if status == nil {
		return "", args
	}

	condition := fmt.Sprintf("v.status = $%d", len(args)+1)
	args = append(args, *status)
	return condition, args
}

func listByCompCreatedAtCursorToSQL(cursor listcompsearch.CreatedAtCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"(v.created_at, v.id) < ($%d, $%d)",
		len(args)+1, len(args)+2)

	args = append(args, cursor.CreatedAt, cursor.ID)
	return condition, args
}

//--------------
// Orders to SQL
//--------------

const (
	publishedAtDescOrderBy string = "v.published_at DESC, v.id DESC"
	salaryDescOrderBy      string = "v.salary_from DESC NULLS LAST, v.salary_to DESC NULLS LAST, v.id DESC"
	salaryAscOrderBy       string = "v.salary_from ASC, v.salary_to ASC NULLS LAST, v.id ASC"
)

func listVacOrderToSQL(order listsearch.Order) string {
	switch order {
	case listsearch.OrderPublishedAtDesc:
		return publishedAtDescOrderBy
	case listsearch.OrderSalaryDesc:
		return salaryDescOrderBy
	case listsearch.OrderSalaryAsc:
		return salaryAscOrderBy
	case listsearch.OrderRelevance:
		return ""
	}

	return ""
}

func buildVacSearchResult(
	vacancies []views.PublishedVacSummary,
	order listsearch.Order,
	limit int,
) (*listsearch.SearchResult, error) {
	result := &listsearch.SearchResult{
		Vacancies: vacancies,
	}

	result.HasNext = len(vacancies) > limit

	if !result.HasNext {
		return result, nil
	}

	cursorVac := vacancies[limit]
	result.Vacancies = result.Vacancies[:limit]

	switch order {
	case listsearch.OrderPublishedAtDesc:
		result.NextCursor = &listsearch.PublishedAtCursor{
			PublishedAt: cursorVac.PublishedAt,
			ID:          cursorVac.ID,
		}

	case listsearch.OrderSalaryAsc, listsearch.OrderSalaryDesc:
		if cursorVac.SalaryFrom == nil || cursorVac.SalaryTo == nil {
			result.HasNext = false
		} else {
			result.NextCursor = &listsearch.SalaryCursor{
				SalaryFrom: *cursorVac.SalaryFrom,
				SalaryTo:   *cursorVac.SalaryTo,
				ID:         cursorVac.ID,
			}
		}

	case listsearch.OrderRelevance:
		return nil, common.ErrUnsupportedListOrder

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	return result, nil
}

func buildCompVacSearchResult(
	vacancies []views.MemberVacSummary,
	order listcompsearch.Order,
	limit int,
) (*listcompsearch.SearchResult, error) {
	result := &listcompsearch.SearchResult{
		Vacancies: vacancies,
	}

	result.HasNext = len(vacancies) > limit

	if !result.HasNext {
		return result, nil
	}

	cursorVac := vacancies[limit]
	result.Vacancies = result.Vacancies[:limit]

	switch order {
	case listcompsearch.OrderCreatedAtDesc:
		result.NextCursor = &listcompsearch.CreatedAtCursor{
			CreatedAt: cursorVac.CreatedAt,
			ID:        cursorVac.ID,
		}

	case listcompsearch.OrderRelevance:
		return nil, common.ErrUnsupportedListOrder

	default:
		return nil, common.ErrUnsupportedListOrder
	}

	return result, nil
}
