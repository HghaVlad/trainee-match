package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/HghaVlad/trainee-match/backend/candidate/internal/domain"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_candidate_by_user_id"
	"github.com/HghaVlad/trainee-match/backend/candidate/internal/usecase/get_resume"
	candidatepb "github.com/HghaVlad/trainee-match/backend/contracts/go/candidate/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	candidatepb.UnimplementedCandidateServiceServer

	getByUserUC *get_candidate_by_user_id.UseCase
	getResumeUC *get_resume.UseCase
}

func New(getByUser *get_candidate_by_user_id.UseCase, getResume *get_resume.UseCase) *Server {
	return &Server{getByUserUC: getByUser, getResumeUC: getResume}
}

func (s *Server) GetCandidateByUserId(
	ctx context.Context,
	req *candidatepb.GetCandidateByUserIdRequest,
) (*candidatepb.GetCandidateResponse, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	out, err := s.getByUserUC.Execute(ctx, uid)
	if err != nil {
		if errors.Is(err, domain.ErrCandidateNotFound) {
			return nil, status.Error(codes.NotFound, "candidate not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &candidatepb.GetCandidateResponse{
		Id:       out.ID.String(),
		UserId:   out.UserID.String(),
		Phone:    out.Phone,
		Telegram: out.Telegram,
		City:     out.City,
		Birthday: out.Birthday.String(),
	}, nil
}

func (s *Server) GetResumeById(
	ctx context.Context,
	req *candidatepb.GetResumeByIdRequest,
) (*candidatepb.GetResumeResponse, error) {
	if req == nil || req.ResumeId == "" || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "resume_id and user_id is required")
	}

	rid, err := uuid.Parse(req.ResumeId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid resume_id format")
	}

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	r, err := s.getResumeUC.GetById(ctx, rid, uid)
	if err != nil {
		if errors.Is(err, domain.ErrResumeNotFound) {
			return nil, status.Error(codes.NotFound, "resume not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &candidatepb.GetResumeResponse{
		Id:          r.ID.String(),
		CandidateId: r.CandidateID.String(),
		Name:        r.Name,
		Status:      fmt.Sprint(r.Status),
	}, nil
}
