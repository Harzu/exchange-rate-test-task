package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Harzu/exchange-rate-test-task/internal/entities"

	"github.com/Harzu/exchange-rate-test-task/internal/services/rates"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/chi/v5"
)

type Container struct {
	ratesService *rates.Service
}

func NewHTTPContainer(ratesService *rates.Service) *Container {
	return &Container{
		ratesService: ratesService,
	}
}

func (c *Container) Mux() *chi.Mux {
	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	mux.Get("/service/price", c.getRate)
	return mux
}

func (c *Container) getRate(w http.ResponseWriter, r *http.Request) {
	sourceSymbols := strings.Split(chi.URLParam(r, "fsyms"), ",")
	targetSymbols := strings.Split(chi.URLParam(r, "fsyms"), ",")
	if len(sourceSymbols) == 0 || len(targetSymbols) == 0 {
		http.Error(w, "invalid params", http.StatusBadRequest)
		return
	}

	pairRates, err := c.ratesService.GetRate(r.Context(), sourceSymbols, targetSymbols)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ratesMapping := map[string]map[string]entities.Rate{}
	for _, pairRate := range pairRates {
		if _, ok := ratesMapping[pairRate.SourceSymbol]; !ok {
			ratesMapping[pairRate.SourceSymbol][pairRate.TargetSymbol] = pairRate.Rate
			continue
		}

		ratesMapping[pairRate.SourceSymbol][pairRate.TargetSymbol] = pairRate.Rate
	}

	type response struct {
		RAW     map[string]map[string]entities.Rate
		DISPLAY map[string]map[string]entities.Rate
	}

	resp := &response{
		RAW:     ratesMapping,
		DISPLAY: ratesMapping,
	}

	responseBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to serialize response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(responseBytes); err != nil {
		//	todo: logger
		return
	}
}
