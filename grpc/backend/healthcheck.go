package grpc_backend

import (
	"context"

	"github.com/gogo/status"
	healthcheckpb "github.com/kei6u/dogfood/proto/v1/healthcheck"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
)

func (s *Server) LivenessProbe(context.Context, *healthcheckpb.LivenessProbeRequest) (*healthcheckpb.LivenessProbeResponse, error) {
	return &healthcheckpb.LivenessProbeResponse{}, nil
}

func (s *Server) ReadinessProbe(context.Context, *healthcheckpb.ReadinessProbeRequest) (*healthcheckpb.ReadinessProbeResponse, error) {
	if err := s.db.Ping(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ping database: %s", err)
	}
	return &healthcheckpb.ReadinessProbeResponse{}, nil
}

func (s *Server) StartupProbe(context.Context, *healthcheckpb.StartupProbeRequest) (*healthcheckpb.StartupProbeResponse, error) {
	if err := s.db.Ping(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ping database: %s", err)
	}
	return &healthcheckpb.StartupProbeResponse{}, nil
}
