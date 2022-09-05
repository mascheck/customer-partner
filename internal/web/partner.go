package web

import (
	"customer-partner/internal/domain"
	"customer-partner/internal/entities"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PartnerService interface {
	GetPartners(opts domain.GetPartnersOpts) []entities.Partner
	GetPartner(id string) (entities.Partner, error)
}

func NewPartnerAPI(service PartnerService) *PartnerAPI {
	return &PartnerAPI{service: service}
}

// PartnerAPI provides the functionality to host the Matching Customer & Partner api
type PartnerAPI struct {
	service PartnerService
}

// ListenAndServe starts serving the api.
// It is a blocking operation.
func (a *PartnerAPI) ListenAndServe() {
	http.HandleFunc("/partners", a.GetPartners)
	http.HandleFunc("/partners/", a.GetPartner)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (a *PartnerAPI) GetPartners(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Endpoint Hit: getPartners")
		err := validateGetPartnersRequest(r.URL.Query())
		if err != nil {
			http.Error(w, fmt.Sprintf("Bad request: %s", err.Error()), http.StatusBadRequest)
			return
		}
		opts, err := getPartnersOptsFromQuery(r.URL.Query())
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(a.service.GetPartners(opts))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *PartnerAPI) GetPartner(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Println("Endpoint Hit: getPartner")
		id := strings.TrimPrefix(r.URL.Path, "/partners/")
		partners, err := a.service.GetPartner(id)
		if errors.Is(err, entities.ErrRecordNotExist) {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(partners)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPartnersOptsFromQuery(params url.Values) (domain.GetPartnersOpts, error) {
	long, err := strconv.ParseFloat(params.Get("long"), 64)
	if err != nil {
		return domain.GetPartnersOpts{}, err
	}
	lat, err := strconv.ParseFloat(params.Get("lat"), 64)
	if err != nil {
		return domain.GetPartnersOpts{}, err
	}
	opts := domain.GetPartnersOpts{
		Material:            params.Get("material"),
		CustomerAddressLong: long,
		CustomerAddressLat:  lat,
	}
	return opts, nil
}

func validateGetPartnersRequest(params url.Values) error {
	if !params.Has("material") {
		return ErrMissingArgument("material")
	}
	if !stringInSlice(params.Get("material"), []string{"wood", "carpet", "tiles"}) {
		return ErrInvalidInput("material")
	}
	if !params.Has("long") {
		return ErrMissingArgument("long")
	}
	if long, err := strconv.ParseFloat(params.Get("long"), 64); err != nil || long < -180 || long > 180 {
		return ErrInvalidInput("long")
	}
	if _, err := strconv.ParseFloat(params.Get("long"), 64); err != nil {
		return ErrInvalidInput("long")
	}
	if !params.Has("lat") {
		return ErrMissingArgument("lat")
	}
	if lat, err := strconv.ParseFloat(params.Get("lat"), 64); err != nil || lat < -90 || lat > 90 {
		return ErrInvalidInput("lat")
	}
	return nil
}

func ErrMissingArgument(name string) error {
	return fmt.Errorf("parameter %s missing", name)
}

func ErrInvalidInput(name string) error {
	return fmt.Errorf("invalid input for parameter %s", name)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
