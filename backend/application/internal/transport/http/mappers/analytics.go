package mappers

import (
	"github.com/HghaVlad/trainee-match/backend/application/internal/transport/http/oapi"
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
