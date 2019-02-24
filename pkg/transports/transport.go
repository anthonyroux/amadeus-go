package transports

import (
	pbComn "amadeus-go/api/amadeus/comn"
	pbFunc "amadeus-go/api/amadeus/func"
	pbType "amadeus-go/api/amadeus/type"
	"amadeus-go/pkg/endpoints"
	sv "amadeus-go/pkg/services"
	"errors"
	"net/http"

	"context"
	"github.com/go-kit/kit/log"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	FlightLowFareSearchHandler     grpcTransport.Handler
	FlightInspirationSearchHandler grpcTransport.Handler
}

func (s *grpcServer) FlightLowFareSearch(ctx context.Context, req *pbFunc.FlightLowFareSearchRequest) (*pbType.FlightLowFareSearchResponse, error) {
	_, resp, err := s.FlightLowFareSearchHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	response := resp.(*pbType.FlightLowFareSearchResponse)
	return response, nil
}

func (s *grpcServer) FlightInspirationSearch(ctx context.Context, req *pbFunc.FlightInspirationSearchRequest) (*pbType.FlightInspirationSearchResponse, error) {
	_, resp, err := s.FlightLowFareSearchHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	response := resp.(*pbType.FlightInspirationSearchResponse)
	return response, nil
}

func NewGRPCServer(endpoints *endpoints.AmadeusEndpointSet, logger log.Logger) (s pbFunc.AmadeusServiceServer) {
	s = &grpcServer{
		FlightLowFareSearchHandler: grpcTransport.NewServer(
			endpoints.FlightLowFareSearchEndpoint,
			decodeFlightLowFareSearchRequest,
			encodeFlightLowFareSearchResponse,
		),
		FlightInspirationSearchHandler: grpcTransport.NewServer(
			endpoints.FlightInspirationSearchEndpoint,
			decodeFlightInspirationSearchRequest,
			encodeFlightInspirationSearchResponse,
		),
	}

	return
}

// =============================================================================
func decodeFlightLowFareSearchRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pbFunc.FlightLowFareSearchRequest)
	if !ok {
		return nil, errors.New("your request is not of type <FlightLowFareSearchRequest>")
	}
	return &sv.FlightLowFareSearchRequest{
		Origin:        req.Origin,
		DepartureDate: req.DepartureDate,
		Destination:   req.Destination,
		ReturnDate:    req.ReturnDate,
	}, nil
}

func encodeFlightLowFareSearchResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*sv.FlightLowFareSearchResponse)
	if !ok {
		return &pbComn.Exception{
			Code: http.StatusInternalServerError,
			Description: "couldn't convert response to <FlightLowFareSearchResponse>",
		}, nil
	}

	var datas []*pbType.Data
	for _, data := range resp.Data {

		var offerItems []*pbType.OfferItem
		for _, offers := range data.OfferItems {

			var services []*pbType.Service
			for _, service := range offers.Services {

				var segments []*pbType.Segment
				for _, segment := range service.Segments {

					segments = append(segments, &pbType.Segment{
						FlightSegment: &pbType.FlightSegment{
							Duration: segment.FlightSegment.Duration,
							Number:   segment.FlightSegment.Number,
							Aircraft: &pbType.Aircraft{
								Code: segment.FlightSegment.Aircraft.Code,
							},
							Arrival: &pbType.DepartureArrival{
								At:       segment.FlightSegment.Arrival.At,
								IataCode: segment.FlightSegment.Arrival.IataCode,
								Terminal: segment.FlightSegment.Arrival.Terminal,
							},
							Departure: &pbType.DepartureArrival{
								At:       segment.FlightSegment.Departure.At,
								IataCode: segment.FlightSegment.Departure.IataCode,
								Terminal: segment.FlightSegment.Departure.Terminal,
							},
							CarrierCode: segment.FlightSegment.CarrierCode,
							Operating: &pbType.Operating{
								CarrierCode: segment.FlightSegment.Operating.CarrierCode,
								Number:      segment.FlightSegment.Operating.Number,
							},
						},
						PricingDetailPerAdult: &pbType.PricingDetailPerAdult{
							Availability: segment.PricingDetailPerAdult.Availability,
							FareBasis:    segment.PricingDetailPerAdult.FareBasis,
							FareClass:    segment.PricingDetailPerAdult.FareClass,
							TravelClass:  segment.PricingDetailPerAdult.TravelClass,
						},
					})
				}

				services = append(services, &pbType.Service{
					Segments: segments,
				})
			}

			offerItems = append(offerItems, &pbType.OfferItem{
				Services: services,
				Price: &pbType.Price{
					Total:      offers.Price.Total,
					TotalTaxes: offers.Price.TotalTaxes,
				},
				PricePerAdult: &pbType.Price{
					Total:      offers.PricePerAdult.Total,
					TotalTaxes: offers.PricePerAdult.TotalTaxes,
				},
			})
		}

		datas = append(datas, &pbType.Data{
			Id:         data.Id,
			Type:       data.Type,
			OfferItems: offerItems,
		})
	}

	var dictionaries pbType.Dictionaries
	dictionaries.Aircrafts = make(map[string]string)
	dictionaries.Locations = make(map[string]*pbType.LocationDetail)
	dictionaries.Carriers = make(map[string]string)
	dictionaries.Currencies = make(map[string]string)
	for k, v := range resp.Dictionaries.Aircrafts {
		dictionaries.Aircrafts[k] = v
	}
	for k, v := range resp.Dictionaries.Locations {
		if _, ok := dictionaries.Locations[k]; !ok {
			dictionaries.Locations[k] = &pbType.LocationDetail{
				Detail: make(map[string]string),
			}
		}
		for subK, subV := range v {
			dictionaries.Locations[k] = &pbType.LocationDetail{
				Detail: map[string]string{subK: subV},
			}
		}
	}
	for k, v := range resp.Dictionaries.Carriers {
		dictionaries.Aircrafts[k] = v
	}
	for k, v := range resp.Dictionaries.Currencies {
		dictionaries.Aircrafts[k] = v
	}

	meta := pbType.Meta{
		Links: &pbType.Links{
			Self: resp.Meta.Links.Self,
		},
		Currency: resp.Meta.Currency,
		Defaults: &pbType.Defaults{
			Adults:  resp.Meta.Defaults.Adults,
			NonStop: resp.Meta.Defaults.NonStop,
		},
	}

	return &pbType.FlightLowFareSearchResponse{
		Data:         datas,
		Dictionaries: &dictionaries,
		Meta:         &meta,
	}, nil
}

func decodeFlightInspirationSearchRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*pbFunc.FlightInspirationSearchRequest)
	if !ok {
		return nil, errors.New("your request is not of type <FlightInspirationSearchRequest>")
	}
	return &sv.FlightInspirationSearchRequest{
		Origin:   req.Origin,
		MaxPrice: req.MaxPrice,
	}, nil
}

func encodeFlightInspirationSearchResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*sv.FlightInspirationSearchResponse)
	if !ok {
		return &pbComn.Exception{
			Code: http.StatusInternalServerError,
			Description: "couldn't convert response to <FlightInspirationSearchResponse>",
		}, nil
	}

	var datas []*pbType.InspirationData
	for _, data := range resp.InspirationData {
		datas = append(datas, &pbType.InspirationData{
			Type:          data.Type,
			Origin:        data.Origin,
			Destination:   data.Destination,
			DepartureDate: data.DepartureDate,
			ReturnDate:    data.ReturnDate,
			Price: &pbType.Price{
				Total: data.Price.Total,
			},
			Links: &pbType.InspirationLinks{
				FlightDates: data.Links.FlightDates,
				FlightOffers: data.Links.FlightOffers,
			},
		})
	}

	var dictionaries pbType.Dictionaries
	dictionaries.Locations = make(map[string]*pbType.LocationDetail)
	dictionaries.Currencies = make(map[string]string)
	for k, v := range resp.Dictionaries.Locations {
		if _, ok := dictionaries.Locations[k]; !ok {
			dictionaries.Locations[k] = &pbType.LocationDetail{
				Detail: make(map[string]string),
			}
		}
		for subK, subV := range v {
			dictionaries.Locations[k] = &pbType.LocationDetail{
				Detail: map[string]string{subK: subV},
			}
		}
	}
	for k, v := range resp.Dictionaries.Currencies {
		dictionaries.Aircrafts[k] = v
	}

	meta := pbType.Meta{
		Links: &pbType.Links{
			Self: resp.Meta.Links.Self,
		},
		Currency: resp.Meta.Currency,
		Defaults: &pbType.Defaults{
			DepartureDate: resp.Meta.Defaults.DepartureDate,
			OneWay: resp.Meta.Defaults.OneWay,
			Duration: resp.Meta.Defaults.Duration,
			NonStop: resp.Meta.Defaults.NonStop,
			ViewBy: resp.Meta.Defaults.ViewBy,
		},
	}

	return &pbType.FlightInspirationSearchResponse{
		Data:         datas,
		Dictionaries: &dictionaries,
		Meta:         &meta,
	}, nil
}
