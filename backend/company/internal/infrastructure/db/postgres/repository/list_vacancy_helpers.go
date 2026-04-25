package repository

import (
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

const (
	andSalaryNotNull string = " AND v.salary_from IS NOT NULL AND v.salary_to IS NOT NULL"
)

func listVacRequirementsToSQL(requirements *list.Requirements) (string, []any) {
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

func applySalary(r *list.Requirements, conds *[]string, args *[]any) {
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

func applyHours(r *list.Requirements, conds *[]string, args *[]any) {
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

func applyDuration(r *list.Requirements, conds *[]string, args *[]any) {
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

func applyWorkFormat(r *list.Requirements, conds *[]string, args *[]any) {
	if r.WorkFormat != nil && len(*r.WorkFormat) > 0 {
		addCondition(conds, args, "v.work_format = ANY($%d)", pq.Array(*r.WorkFormat))
	}
}

func applyCompanies(r *list.Requirements, conds *[]string, args *[]any) {
	if r.Companies != nil && len(*r.Companies) > 0 {
		addCondition(conds, args, "v.company_id = ANY($%d)", pq.Array(*r.Companies))
	}
}

func applyCity(r *list.Requirements, conds *[]string, args *[]any) {
	if r.City != nil && len(*r.City) > 0 {
		addCondition(conds, args, "v.city = ANY($%d)", pq.Array(*r.City))
	}
}

func applyFlags(r *list.Requirements, conds *[]string, args *[]any) {
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
func listVacCursorToSQL(order list.Order, cursor any, args []any) (string, []any) {
	switch c := cursor.(type) {
	case *list.PublishedAtCursor:
		return publishedAtCursorToSQL(*c, args)
	case *list.SalaryCursor:
		return salaryCursorToSQL(order, *c, args)
	}

	return "", args
}

func publishedAtCursorToSQL(cursor list.PublishedAtCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"(v.published_at, v.id) < ($%d, $%d)",
		len(args)+1, len(args)+2)

	args = append(args, cursor.PublishedAt, cursor.ID)
	return condition, args
}

func salaryCursorToSQL(order list.Order, cursor list.SalaryCursor, args []any) (string, []any) {
	var condition string

	if order == list.OrderSalaryDesc {
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

//--------------
// Orders to SQL
//--------------

const (
	publishedAtDescOrderBy string = "v.published_at DESC, v.id DESC"
	salaryDescOrderBy      string = "v.salary_from DESC NULLS LAST, v.salary_to DESC NULLS LAST, v.id DESC"
	salaryAscOrderBy       string = "v.salary_from ASC, v.salary_to ASC NULLS LAST, v.id ASC"
)

func listVacOrderToSQL(order list.Order) string {
	switch order {
	case list.OrderPublishedAtDesc:
		return publishedAtDescOrderBy
	case list.OrderSalaryDesc:
		return salaryDescOrderBy
	case list.OrderSalaryAsc:
		return salaryAscOrderBy
	}

	return ""
}
