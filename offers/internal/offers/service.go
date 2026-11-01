// Package offers implements coolinary.offers.v1.OffersService.
package offers

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"

	commonv1 "coolinary/offers/gen/coolinary/common/v1"
	offersv1 "coolinary/offers/gen/coolinary/offers/v1"
	"coolinary/offers/gen/coolinary/offers/v1/offersv1connect"
	"coolinary/offers/internal/offers/store"
)

type Service struct {
	q   *store.Queries
	log *slog.Logger
}

var _ offersv1connect.OffersServiceHandler = (*Service)(nil)

func New(q *store.Queries, log *slog.Logger) *Service { return &Service{q: q, log: log} }

// ListOffers returns current offers for a country, paged by offer id.
// Ordering by distance happens on the phone; the server only knows the postcode area.
func (s *Service) ListOffers(ctx context.Context, req *connect.Request[offersv1.ListOffersRequest]) (*connect.Response[offersv1.ListOffersResponse], error) {
	size := req.Msg.GetPage().GetPageSize()
	if size <= 0 || size > 100 {
		size = 20
	}
	rows, err := s.q.ListCurrentOffers(ctx, store.ListCurrentOffersParams{
		CountryCode: strings.ToUpper(req.Msg.GetCountryCode()),
		ID:          req.Msg.GetPage().GetPageToken(),
		Limit:       size + 1,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	res := &offersv1.ListOffersResponse{Page: &commonv1.PageResponse{}}
	for i, r := range rows {
		if int32(i) == size {
			res.Page.NextPageToken = rows[i-1].ID
			break
		}
		res.Offers = append(res.Offers, &offersv1.Offer{
			Id: r.ID, ShopId: r.ShopID, ProductName: r.ProductName,
			Price:      &commonv1.Money{AmountMinor: r.PriceMinor, Currency: r.Currency},
			Conditions: r.Conditions,
			ValidFrom:  r.ValidFrom.Time.Format("2006-01-02"),
			ValidTo:    r.ValidTo.Time.Format("2006-01-02"),
			SourceUrl:  r.SourceUrl,
		})
	}
	return connect.NewResponse(res), nil
}
