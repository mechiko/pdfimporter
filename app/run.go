package app

import (
	"context"
)

// если потребуется заготовка запуска как задачи через групповой контекст
func (a *app) Run(ctx context.Context, cancel context.CancelFunc) error {
	<-ctx.Done()
	return nil
}
