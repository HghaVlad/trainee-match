package repository

import (
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/HghaVlad/trainee-match/backend/company/internal/usecase/vacancy/list"
)

const (
	publishedAtDescOrderBy string = "ORDER BY v.published_at DESC, v.id DESC"
	salaryDescOrderBy      string = "ORDER BY v.salary_from DESC NULLS LAST, v.salary_to DESC NULLS LAST, v.id DESC"
	salaryAscOrderBy       string = "ORDER BY v.salary_from ASC, v.salary_to ASC NULLS LAST, v.id ASC"
)

const (
	andSalaryNotNull string = " AND v.salary_from IS NOT NULL AND v.salary_to IS NOT NULL"
)

func listVacRequirementsToSQL(requirements *list_vacancy.Requirements) (string, []any) {
	if requirements == nil {
		return "1=1", nil
	}

	conditions := make([]string, 0)
	args := make([]any, 0)

	// ---- Salary ----
	if requirements.Salary != nil {
		if requirements.Salary.Min != nil {
			args = append(args, *requirements.Salary.Min)
			conditions = append(conditions,
				fmt.Sprintf("v.salary_to >= $%d", len(args)))
		}
		if requirements.Salary.Max != nil {
			args = append(args, *requirements.Salary.Max)
			conditions = append(conditions,
				fmt.Sprintf("v.salary_from <= $%d", len(args)))
		}
	}

	// ---- HoursPerWeek ----
	if requirements.HoursPerWeek != nil {
		if requirements.HoursPerWeek.Min != nil {
			args = append(args, *requirements.HoursPerWeek.Min)
			conditions = append(conditions,
				fmt.Sprintf("v.hours_per_week_to >= $%d", len(args)))
		}
		if requirements.HoursPerWeek.Max != nil {
			args = append(args, *requirements.HoursPerWeek.Max)
			conditions = append(conditions,
				fmt.Sprintf("v.hours_per_week_from <= $%d", len(args)))
		}
	}

	// ---- Duration ----
	if requirements.Duration != nil {
		if requirements.Duration.Min != nil {
			args = append(args, *requirements.Duration.Min)
			conditions = append(conditions,
				fmt.Sprintf("v.duration_to_days >= $%d", len(args)))
		}
		if requirements.Duration.Max != nil {
			args = append(args, *requirements.Duration.Max)
			conditions = append(conditions,
				fmt.Sprintf("v.duration_from_days <= $%d", len(args)))
		}
	}

	// ---- WorkFormat ----
	if requirements.WorkFormat != nil && len(*requirements.WorkFormat) > 0 {
		args = append(args, pq.Array(*requirements.WorkFormat))
		conditions = append(conditions,
			fmt.Sprintf("v.work_format = ANY($%d)", len(args)))
	}

	// ---- Companies ----
	if requirements.Companies != nil && len(*requirements.Companies) > 0 {
		args = append(args, pq.Array(*requirements.Companies))
		conditions = append(conditions,
			fmt.Sprintf("v.company_id = ANY($%d)", len(args)))
	}

	// ---- City ----
	if requirements.City != nil && len(*requirements.City) > 0 {
		args = append(args, pq.Array(*requirements.City))
		conditions = append(conditions,
			fmt.Sprintf("v.city = ANY($%d)", len(args)))
	}

	// ---- IsPaid ----
	if requirements.IsPaid != nil {
		args = append(args, *requirements.IsPaid)
		conditions = append(conditions,
			fmt.Sprintf("v.is_paid = $%d", len(args)))
	}

	// ---- InternshipToOffer ----
	if requirements.InternshipToOffer != nil {
		args = append(args, *requirements.InternshipToOffer)
		conditions = append(conditions,
			fmt.Sprintf("v.internship_to_offer = $%d", len(args)))
	}

	// ---- FlexibleSchedule ----
	if requirements.FlexibleSchedule != nil {
		args = append(args, *requirements.FlexibleSchedule)
		conditions = append(conditions,
			fmt.Sprintf("v.flexible_schedule = $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "1=1", args
	}

	return strings.Join(conditions, " AND "), args
}

func listVacCursorToSQL(order list_vacancy.Order, cursor any, args []any) (condition string, newArgs []any) {
	switch c := cursor.(type) {
	case *list_vacancy.PublishedAtCursor:
		return publishedAtCursorToSQL(*c, args)
	case *list_vacancy.SalaryCursor:
		return salaryCursorToSQL(order, *c, args)
	}

	return
}

func publishedAtCursorToSQL(cursor list_vacancy.PublishedAtCursor, args []any) (string, []any) {
	condition := fmt.Sprintf(
		"(v.published_at, v.id) < ($%d, $%d)",
		len(args)+1, len(args)+2)

	args = append(args, cursor.PublishedAt, cursor.Id)
	return condition, args
}

func salaryCursorToSQL(order list_vacancy.Order, cursor list_vacancy.SalaryCursor, args []any) (string, []any) {
	var condition string

	if order == list_vacancy.OrderSalaryDesc {
		condition = fmt.Sprintf(
			"(v.salary_from, v.salary_to, v.id) < ($%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3)
	} else {
		condition = fmt.Sprintf(
			"(v.salary_from, v.salary_to, v.id) > ($%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3)
	}

	args = append(args, cursor.SalaryFrom, cursor.SalaryTo, cursor.Id)
	return condition, args
}

func listVacOrderToSQL(order list_vacancy.Order) string {
	switch order {
	case list_vacancy.OrderPublishedAtDesc:
		return publishedAtDescOrderBy
	case list_vacancy.OrderSalaryDesc:
		return salaryDescOrderBy
	case list_vacancy.OrderSalaryAsc:
		return salaryAscOrderBy
	}

	return ""
}
