package service

import (
	"context"
)

type OnDataReceive func(ctx context.Context, station uint16, ioa uint, data any) error

type OnIecReceive func(station uint16, ioa uint, data any)
