package domain

import "context"

type PrometheusServer interface {
	Run()
	Stop(ctx context.Context)
}
