// Package kitchen implements coolinary.kitchen.v1.KitchenService.
package kitchen

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	commonv1 "coolinary/kitchen/gen/coolinary/common/v1"
	kitchenv1 "coolinary/kitchen/gen/coolinary/kitchen/v1"
	"coolinary/kitchen/gen/coolinary/kitchen/v1/kitchenv1connect"
	"coolinary/kitchen/internal/kitchen/store"
)

const cacheTTL = 24 * time.Hour

type Service struct {
	q     *store.Queries
	cache *redis.Client
	log   *slog.Logger
}

var _ kitchenv1connect.KitchenServiceHandler = (*Service)(nil)

func New(q *store.Queries, cache *redis.Client, log *slog.Logger) *Service {
	return &Service{q: q, cache: cache, log: log}
}

// ResolveBarcode walks the lookup chain. For now only the shared catalogue
// is implemented; Open Food Facts and GS1 come with the barcode feature.
func (s *Service) ResolveBarcode(ctx context.Context, req *connect.Request[kitchenv1.ResolveBarcodeRequest]) (*connect.Response[kitchenv1.ResolveBarcodeResponse], error) {
	res, err := s.resolveBarcode(ctx, req.Msg.GetEan())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(res), nil
}

func (s *Service) resolveBarcode(ctx context.Context, ean string) (*kitchenv1.ResolveBarcodeResponse, error) {
	if len(ean) != 8 && len(ean) != 12 && len(ean) != 13 && len(ean) != 14 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("ean must have 8, 12, 13 or 14 digits"))
	}

	key := "kitchen:ean:" + ean
	if b, err := s.cache.Get(ctx, key).Bytes(); err == nil {
		var p kitchenv1.Product
		if json.Unmarshal(b, &p) == nil {
			return &kitchenv1.ResolveBarcodeResponse{Found: true, Product: &p, Source: kitchenv1.ResolveSource_RESOLVE_SOURCE_SHARED_CATALOGUE}, nil
		}
	}

	row, err := s.q.GetProductByEAN(ctx, ean)
	if errors.Is(err, pgx.ErrNoRows) {
		return &kitchenv1.ResolveBarcodeResponse{Found: false}, nil
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	p := &kitchenv1.Product{Ean: row.Ean, Name: row.Name, Brand: row.Brand}
	if row.PackValue.Valid {
		p.PackSize = &commonv1.Quantity{Value: row.PackValue.Float64, Unit: commonv1.Unit(row.PackUnit)}
	}
	if b, err := json.Marshal(p); err == nil {
		s.cache.Set(ctx, key, b, cacheTTL)
	}
	return &kitchenv1.ResolveBarcodeResponse{Found: true, Product: p, Source: kitchenv1.ResolveSource_RESOLVE_SOURCE_SHARED_CATALOGUE}, nil
}
