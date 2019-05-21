package log

import (
	"github.com/favclip/ucon"
)

var _ (ucon.MiddlewareFunc) = LoggerMiddleware

func LoggerMiddleware(b *ucon.Bubble) error {
	ctx, logger, w := NewRequestLogger(b.Context, b.R, b.W)
	defer logger.Finish()
	b.W = w
	b.R = b.R.WithContext(ctx)
	b.Context = ctx

	return b.Next()
}
