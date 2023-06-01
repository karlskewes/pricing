package pricing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type (
	AddBrandRequest struct {
		Name string `json:"name"`
	}
	AddBrandResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	GetBrandRequest struct {
		Name string `json:"name"`
	}
	GetBrandResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	AddPriceRequest struct {
		BrandID   int       `json:"brand_id"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		ProductID int       `json:"product_id"`
		Priority  int       `json:"priority"`
		Price     int       `json:"price"`
		Curr      string    `json:"curr"`
	}
	GetPriceRequest struct {
		BrandID   int       `json:"brand_id"`
		ProductID int       `json:"product_id"`
		Date      time.Time `json:"date"`
		StringID  string    `json:"string_id"`
	}
	GetPriceResponse struct {
		BrandID   int    `json:"brand_id"`
		ProductID int    `json:"product_id"`
		Price     string `json:"price"`
		Curr      string `json:"curr"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		StringID  string `json:"string_id"`
	}
)

// Handler will expose our service via an "open host service"
type Handler struct {
	svc *Service
	// logger zerolog.Logger // log http queries, etc
}

func NewHandler(svc *Service) (*Handler, error) {
	if svc == nil {
		return nil, errors.New("service cannot be empty")
	}

	return &Handler{svc: svc}, nil
}

func (h Handler) AddBrand(w http.ResponseWriter, req *http.Request) {
	// TODO, not implemented.
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) GetBrand(w http.ResponseWriter, req *http.Request) {
	// Consider restricting method
	// if req.Method != "GET" {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// }

	name := req.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		// TODO write error json response
		return
	}

	brand, err := h.svc.GetBrand(req.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(GetBrandResponse(brand))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(res)
	if err != nil {
		// TODO: log error
		return // silences staticcheck
	}
}

func (h Handler) AddPrice(w http.ResponseWriter, req *http.Request) {
	// TODO, not implemented.
	w.WriteHeader(http.StatusNotImplemented)
}

func (h Handler) GetPrice(w http.ResponseWriter, req *http.Request) {
	// Consider restricting method
	// if req.Method != "GET" {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// }

	brandID := req.URL.Query().Get("brand_id")
	if brandID == "" {
		w.WriteHeader(http.StatusBadRequest)
		// TODO write error json response
		return
	}

	productID := req.URL.Query().Get("product_id")
	if productID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	date := req.URL.Query().Get("date")
	if date == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	formattedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// custom string identifier mentioned in pdf
	stringID := req.URL.Query().Get("string_id")
	if stringID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bid, err := strconv.Atoi(brandID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pid, err := strconv.Atoi(productID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	price, err := h.svc.GetPrice(req.Context(), bid, pid, formattedDate)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res, err := json.Marshal(GetPriceResponse{
		BrandID:   price.BrandID,
		ProductID: price.ProductID,
		Price:     fmt.Sprintf("%d.%02d", price.Price/100, price.Price%100), // TODO: use money.Money
		Curr:      price.Curr,
		StartDate: price.StartDate.String(),
		EndDate:   price.EndDate.String(),
		StringID:  stringID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(res)
	if err != nil {
		// TODO: log error
		return // silences staticcheck
	}
}
