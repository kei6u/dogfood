package grpc_backend

import (
	"context"
	"database/sql"
	"time"

	"github.com/gogo/status"
	"github.com/kei6u/dogfood/pkg/ddconfig"
	dogfoodpb "github.com/kei6u/dogfood/proto/v1/dogfood"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func (s *Server) CreateRecord(ctx context.Context, req *dogfoodpb.CreateRecordRequest) (*dogfoodpb.Record, error) {
	if s, ok := tracer.SpanFromContext(ctx); ok {
		span := tracer.StartSpan(ddconfig.GetService(ddconfig.WithServiceSuffix(".CreateRecord")), tracer.ChildOf(s.Context()))
		defer span.Finish()
	}
	eatenAt := time.Now()
	if err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO record (dogfood_name, gram, dog_name, eaten_at) VALUES ($1, $2, $3, $4)",
		req.GetDogfoodName(), req.GetGram(), req.GetDogName(), eatenAt,
	).Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create a record: %s", err)
	}
	return &dogfoodpb.Record{
		DogfoodName: req.GetDogfoodName(),
		Gram:        req.GetGram(),
		DogName:     req.GetDogName(),
		EatenAt:     timestamppb.New(eatenAt),
	}, nil
}

func (s *Server) ListRecords(ctx context.Context, req *dogfoodpb.ListRecordsRequest) (*dogfoodpb.ListRecordsResponse, error) {
	if s, ok := tracer.SpanFromContext(ctx); ok {
		span := tracer.StartSpan(ddconfig.GetService(ddconfig.WithServiceSuffix(".ListRecords")), tracer.ChildOf(s.Context()))
		defer span.Finish()
	}
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT * FROM record WHERE eaten_at >= $1 AND eaten_at < $2 LIMIT $3",
		req.GetFrom().AsTime(), req.GetTo().AsTime(), req.GetPageSize(),
	)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "no record exists")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list up records: %s", err)
	}
	var rs []*dogfoodpb.Record
	for rows.Next() {
		var r dogfoodpb.Record
		var t time.Time
		if err := rows.Scan(&r.DogfoodName, &r.Gram, &r.DogName, &t); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to list up records: %s", err)
		}
		r.EatenAt = timestamppb.New(t)
		rs = append(rs, &r)
	}
	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list data: %s", err)
	}
	return &dogfoodpb.ListRecordsResponse{
		Records: rs,
		To:      rs[len(rs)-1].GetEatenAt(),
	}, nil
}
