//go:build wireinject
// +build wireinject

package api

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/wire"

	"github.com/liuerfire/boxpractice/internal/services"
	"github.com/liuerfire/boxpractice/pkg/store"
)

func InitAPIHandler(ctx context.Context, logger logr.Logger, sqlStore *store.SQLStore) (*API, error) {
	wire.Build(ProvideAPI, services.ProvideHospitalService, services.ProvideEmployeeService, services.ProvideTaskService)
	return &API{}, nil
}
