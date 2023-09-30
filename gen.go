package gen

//go:generate mockgen -destination=internal/server/mocks/server.go -package=mocks github.com/kotkovdev/pow/internal/server QuotesService,Challenger
