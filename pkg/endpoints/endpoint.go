package endpoints

import (
	sv "amadeus-go/pkg/services"

	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type AmadeusEndpointSet struct {
	FlightLowFareSearchEndpoint            endpoint.Endpoint
	FlightInspirationSearchEndpoint        endpoint.Endpoint
	FlightMostTraveledDestinationsEndpoint endpoint.Endpoint
	FlightMostBookedDestinationsEndpoint   endpoint.Endpoint
	FlightBusiestTravelingPeriodEndpoint   endpoint.Endpoint
	AirportNearestRelevantEndpoint         endpoint.Endpoint
	AirportAndCitySearchEndpoint           endpoint.Endpoint
}

func (s AmadeusEndpointSet) FlightLowFareSearch(ctx context.Context, request *sv.FlightLowFareSearchRequest) (*sv.Response, error) {
	resp, err := s.FlightLowFareSearchEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) FlightInspirationSearch(ctx context.Context, request *sv.FlightInspirationSearchRequest) (*sv.Response, error) {
	resp, err := s.FlightInspirationSearchEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) FlightMostTraveledDestinations(ctx context.Context, request *sv.FlightMostTraveledDestinationsRequest) (*sv.Response, error) {
	resp, err := s.FlightMostTraveledDestinationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) FlightMostBookedDestinations(ctx context.Context, request *sv.FlightMostBookedDestinationsRequest) (*sv.Response, error) {
	resp, err := s.FlightMostBookedDestinationsEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) FlightBusiestTravelingPeriod(ctx context.Context, request *sv.FlightBusiestTravelingPeriodRequest) (*sv.Response, error) {
	resp, err := s.FlightBusiestTravelingPeriodEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) AirportNearestRelevant(ctx context.Context, request *sv.AirportNearestRelevantRequest) (*sv.Response, error) {
	resp, err := s.AirportNearestRelevantEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func (s AmadeusEndpointSet) AirportAndCitySearch(ctx context.Context, request *sv.AirportAndCitySearchRequest) (*sv.Response, error) {
	resp, err := s.AirportAndCitySearchEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}

	response := resp.(*sv.Response)
	return response, nil
}

func NewEndpointSet(srv sv.AmadeusService, logger log.Logger) *AmadeusEndpointSet {
	var (
		flightLowFareSearchEndpoint            endpoint.Endpoint
		flightInspirationSearchEndpoint        endpoint.Endpoint
		flightMostTraveledDestinationsEndpoint endpoint.Endpoint
		flightMostBookedDestinationsEndpoint   endpoint.Endpoint
		flightBusiestTravelingPeriodEndpoint   endpoint.Endpoint
		airportNearestRelevantEndpoint         endpoint.Endpoint
		airportAndCitySearchEndpoint           endpoint.Endpoint
	)

	flightLowFareSearchEndpoint = makeFlightLowFareSearchEndpoint(srv)
	flightLowFareSearchEndpoint = loggingMiddleware(logger, "FlightLowFareSearch")(flightLowFareSearchEndpoint)

	flightInspirationSearchEndpoint = makeFlightInspirationSearchEndpoint(srv)
	flightInspirationSearchEndpoint = loggingMiddleware(logger, "FlightInspirationSearch")(flightInspirationSearchEndpoint)

	flightMostTraveledDestinationsEndpoint = makeFlightMostTraveledDestinationsEndpoint(srv)
	flightMostTraveledDestinationsEndpoint = loggingMiddleware(logger, "FlightMostTraveledDestinations")(flightMostTraveledDestinationsEndpoint)

	flightMostBookedDestinationsEndpoint = makeFlightMostBookedDestinationsEndpoint(srv)
	flightMostBookedDestinationsEndpoint = loggingMiddleware(logger, "FlightMostBookedDestinations")(flightMostBookedDestinationsEndpoint)

	flightBusiestTravelingPeriodEndpoint = makeFlightBusiestTravelingPeriodEndpoint(srv)
	flightBusiestTravelingPeriodEndpoint = loggingMiddleware(logger, "FlightBusiestTravelingPeriod")(flightBusiestTravelingPeriodEndpoint)

	airportNearestRelevantEndpoint = makeAirportNearestRelevantEndpoint(srv)
	airportNearestRelevantEndpoint = loggingMiddleware(logger, "AirportNearestRelevant")(airportNearestRelevantEndpoint)

	airportAndCitySearchEndpoint = makeAirportAndCitySearchEndpoint(srv)
	airportAndCitySearchEndpoint = loggingMiddleware(logger, "AirportAndCitySearch")(airportAndCitySearchEndpoint)

	return &AmadeusEndpointSet{
		FlightLowFareSearchEndpoint:            flightLowFareSearchEndpoint,
		FlightInspirationSearchEndpoint:        flightInspirationSearchEndpoint,
		FlightMostTraveledDestinationsEndpoint: flightMostTraveledDestinationsEndpoint,
		FlightMostBookedDestinationsEndpoint:   flightMostBookedDestinationsEndpoint,
		FlightBusiestTravelingPeriodEndpoint:   flightBusiestTravelingPeriodEndpoint,
		AirportNearestRelevantEndpoint:         airportNearestRelevantEndpoint,
		AirportAndCitySearchEndpoint:           airportAndCitySearchEndpoint,
	}
}

// ====================================================================================================
func makeFlightLowFareSearchEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.FlightLowFareSearchRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <FlightLowFareSearchRequest>")
		}

		resp, err := srv.FlightLowFareSearch(ctx, req)
		return resp, err
	}
}

func makeFlightInspirationSearchEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.FlightInspirationSearchRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <FlightInspirationSearchRequest>")
		}

		resp, err := srv.FlightInspirationSearch(ctx, req)
		return resp, err
	}
}

func makeFlightMostTraveledDestinationsEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.FlightMostTraveledDestinationsRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <FlightMostTraveledDestinationsRequest>")
		}

		resp, err := srv.FlightMostTraveledDestinations(ctx, req)
		return resp, err
	}
}

func makeFlightMostBookedDestinationsEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.FlightMostBookedDestinationsRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <FlightMostBookedDestinationsRequest>")
		}

		resp, err := srv.FlightMostBookedDestinations(ctx, req)
		return resp, err
	}
}

func makeFlightBusiestTravelingPeriodEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.FlightBusiestTravelingPeriodRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <FlightBusiestTravelingPeriodRequest>")
		}

		resp, err := srv.FlightBusiestTravelingPeriod(ctx, req)
		return resp, err
	}
}

func makeAirportNearestRelevantEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.AirportNearestRelevantRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <AirportNearestRelevantRequest>")
		}

		resp, err := srv.AirportNearestRelevant(ctx, req)
		return resp, err
	}
}

func makeAirportAndCitySearchEndpoint(srv sv.AmadeusService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*sv.AirportAndCitySearchRequest)
		if !ok {
			return nil, errors.New("service did not fetch type <AirportAndCitySearchRequest>")
		}

		resp, err := srv.AirportAndCitySearch(ctx, req)
		return resp, err
	}
}

