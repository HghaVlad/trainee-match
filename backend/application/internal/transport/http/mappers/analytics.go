package mappers

import (
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/dynamics"
	"github.com/HghaVlad/trainee-match/backend/application/internal/usecase/analytics/summary"
)

func AnalyticsSummaryToHTTP(sum summary.Summary) oapi.AnalyticsSummary {
	return oapi.AnalyticsSummary{
		TotalApplications:     sum.TotalApplications,
		ActiveApplications:    sum.ActiveApplications,
		SubmittedCount:        sum.SubmittedCount,
		SeenCount:             sum.SeenCount,
		InterviewCount:        sum.InterviewCount,
		OfferCount:            sum.OfferCount,
		RejectedCount:         sum.RejectedCount,
		WithdrawnCount:        sum.WithdrawnCount,
		ConversionToInterview: sum.ConversionToInterview,
		ConversionToOffer:     sum.ConversionToOffer,
	}
}

func DynamicsCompPeriodToUC(request oapi.GetCompanyDynamicsRequestObject) dynamics.Period {
	period := dynamics.Period{
		From: request.Params.CreatedFrom,
		To:   request.Params.CreatedTo,
	}

	if request.Params.Interval == nil {
		period.Interval = dynamics.IntervalDay
		return period
	}

	switch *request.Params.Interval {
	case oapi.GetCompanyDynamicsParamsIntervalDay:
		period.Interval = dynamics.IntervalDay
	case oapi.GetCompanyDynamicsParamsIntervalWeek:
		period.Interval = dynamics.IntervalWeek
	case oapi.GetCompanyDynamicsParamsIntervalMonth:
		period.Interval = dynamics.IntervalMonth
	default:
		period.Interval = dynamics.Interval(*request.Params.Interval)
	}

	return period
}

func DynamicsVacPeriodToUC(request oapi.GetVacancyDynamicsRequestObject) dynamics.Period {
	period := dynamics.Period{
		From: request.Params.CreatedFrom,
		To:   request.Params.CreatedTo,
	}

	if request.Params.Interval == nil {
		period.Interval = dynamics.IntervalDay
		return period
	}

	switch *request.Params.Interval {
	case oapi.Day:
		period.Interval = dynamics.IntervalDay
	case oapi.Week:
		period.Interval = dynamics.IntervalWeek
	case oapi.Month:
		period.Interval = dynamics.IntervalMonth
	default:
		period.Interval = dynamics.Interval(*request.Params.Interval)
	}

	return period
}

func DynamicsDashboardToHTTP(buckets []dynamics.Bucket) []oapi.ApplicationDynamicsPoint {
	points := make([]oapi.ApplicationDynamicsPoint, len(buckets))

	for i := range buckets {
		points[i] = oapi.ApplicationDynamicsPoint{
			BucketStart:    buckets[i].Start,
			CreatedCount:   buckets[i].SubmittedCount,
			SeenCount:      buckets[i].SeenCount,
			InterviewCount: buckets[i].InterviewCount,
			OfferCount:     buckets[i].OfferCount,
			RejectedCount:  buckets[i].RejectedCount,
			WithdrawnCount: buckets[i].WithdrawnCount,
		}
	}

	return points
}
