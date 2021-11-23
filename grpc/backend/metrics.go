package grpc_backend

import (
	"context"

	dogfoodpb "github.com/kei6u/dogfood/proto/v1/dogfood"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	dogfoodGramGuage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "eaten_dogfood_gram",
			Help: "how much grams dog ate dogfood",
		},
		[]string{"dog", "dogfood"},
	)
	dogfoodNameCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eaten_dogfood_count",
			Help: "the number of eaten dogfood",
		},
		[]string{"dogfood"},
	)
)

func metricsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		switch r := req.(type) {
		case *dogfoodpb.CreateRecordRequest:
			dogfoodGramGuage.With(prometheus.Labels{
				"dog":     r.GetDogName(),
				"dogfood": r.GetDogfoodName(),
			}).Set(float64(r.GetGram()))
			dogfoodNameCount.With(prometheus.Labels{
				"dogfood": r.GetDogfoodName(),
			}).Inc()
		}
		return handler(ctx, req)
	}
}
