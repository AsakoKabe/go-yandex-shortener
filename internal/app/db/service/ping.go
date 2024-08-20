package service

import "context"

// PingService сервис для проверки состояния БД
type PingService interface {
	PingDB(ctx context.Context) error
}
