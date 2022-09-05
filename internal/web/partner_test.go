package web_test

import (
	"customer-partner/internal/domain"
	"customer-partner/internal/entities"
	"customer-partner/internal/web"
	"customer-partner/internal/web/mocks"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery --name PartnerService

func TestPartnerAPI_GetPartner(t *testing.T) {
	type testCase struct {
		name           string
		serviceReturn1 entities.Partner
		serviceReturn2 error
		expStatus      int
		expBody        func() string
	}
	tests := []testCase{
		{
			name:           "Returns 200 with valid body",
			serviceReturn1: entities.Partner{},
			serviceReturn2: nil,
			expStatus:      http.StatusOK,
			expBody: func() string {
				body, _ := json.Marshal(entities.Partner{})
				return fmt.Sprintf("%s\n", body)
			},
		},
		{
			name:           "Returns 404 on ErrRecordNotExist",
			serviceReturn1: entities.Partner{},
			serviceReturn2: entities.ErrRecordNotExist,
			expStatus:      http.StatusNotFound,
			expBody: func() string {
				return "Not found"
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := &mocks.PartnerService{}
			service.On("GetPartner", "123").Return(tt.serviceReturn1, tt.serviceReturn2)
			api := web.NewPartnerAPI(service)

			assert.HTTPStatusCode(t, api.GetPartner, http.MethodGet, "/partners/123", nil, tt.expStatus)
			assert.HTTPBodyContains(t, api.GetPartner, http.MethodGet, "/partners/123", nil, tt.expBody())
			service.AssertExpectations(t)
		})
	}
}

func TestPartnerAPI_GetPartners(t *testing.T) {
	type testCase struct {
		name           string
		urlValues      url.Values
		serviceReturn  []entities.Partner
		expServiceCall bool
		expStatus      int
		expBody        func() string
	}
	tests := []testCase{
		{
			name:           "Returns 400 on missing query parameter 'material'",
			urlValues:      nil,
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: parameter material missing\n" },
		},
		{
			name: "Returns 400 on invalid input for query parameter 'material'",
			urlValues: url.Values{
				"material": []string{"dark matter"},
				"long":     []string{"80.123"},
				"lat":      []string{"42.125"},
			},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: invalid input for parameter material\n" },
		},
		{
			name:           "Returns 400 on missing query parameter 'long'",
			urlValues:      url.Values{"material": []string{"wood"}},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: parameter long missing\n" },
		},
		{
			name:           "Returns 400 on invalid input for query parameter 'long'",
			urlValues:      url.Values{"material": []string{"wood"}, "long": []string{"abc"}},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: invalid input for parameter long\n" },
		},
		{
			name: "Returns 400 on out of bounds for parameter 'long'",
			urlValues: url.Values{
				"material": []string{"wood"},
				"long":     []string{"200"},
				"lat":      []string{"42.125"},
			},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: invalid input for parameter long\n" },
		},
		{
			name:           "Returns 400 on missing query parameter 'lat'",
			urlValues:      url.Values{"material": []string{"wood"}, "long": []string{"80.123"}},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: parameter lat missing\n" },
		},
		{
			name:           "Returns 400 on invalid input for query parameter 'lat'",
			urlValues:      url.Values{"material": []string{"wood"}, "long": []string{"80.123"}, "lat": []string{"ab"}},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: invalid input for parameter lat\n" },
		},
		{
			name: "Returns 400 on out of bounds for query parameter 'lat'",
			urlValues: url.Values{
				"material": []string{"wood"},
				"long":     []string{"80.123"},
				"lat":      []string{"90.134"},
			},
			expServiceCall: false,
			expStatus:      http.StatusBadRequest,
			expBody:        func() string { return "Bad request: invalid input for parameter lat\n" },
		},
		{
			name: "Returns 200 with valid body on empty list",
			urlValues: url.Values{
				"material": []string{"wood"},
				"long":     []string{"80.123"},
				"lat":      []string{"42.125"},
			},
			serviceReturn:  []entities.Partner{},
			expServiceCall: true,
			expStatus:      http.StatusOK,
			expBody:        func() string { return "[]\n" },
		},
		{
			name: "Returns 200 with valid body on filled list",
			urlValues: url.Values{
				"material": []string{"wood"},
				"long":     []string{"80.123"},
				"lat":      []string{"42.125"},
			},
			serviceReturn:  []entities.Partner{{}},
			expServiceCall: true,
			expStatus:      http.StatusOK,
			expBody: func() string {
				body, _ := json.Marshal([]entities.Partner{{}})
				return fmt.Sprintf("%s\n", body)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			service := &mocks.PartnerService{}
			if tt.expServiceCall {
				service.On("GetPartners", domain.GetPartnersOpts{
					Material:            "wood",
					CustomerAddressLong: 80.123,
					CustomerAddressLat:  42.125,
				}).Return(tt.serviceReturn)
			}
			api := web.NewPartnerAPI(service)

			assert.HTTPStatusCode(t, api.GetPartners, http.MethodGet, "/partners", tt.urlValues, tt.expStatus)
			assert.HTTPBodyContains(t, api.GetPartners, http.MethodGet, "/partners", tt.urlValues, tt.expBody())
			service.AssertExpectations(t)
		})
	}
}
