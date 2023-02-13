package app

import (
	"mergeApi/pkg/app/command"
	"mergeApi/pkg/app/query"
)

type Application struct {
	Commands Command
	Queries  Query
}

type Command struct {
	RefreshRate command.RefreshRateHandler
}

type Query struct {
	GetRates query.GetRatesHandler
}
