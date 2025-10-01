package http

import (
	"encoding/json"
	"errors"
	"io"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shopspring/decimal"

	"github.com/example/payment-gateway/internal/application"
	"github.com/example/payment-gateway/internal/domain"
)

// Server wires HTTP routes to application services.
type Server struct {
	transferService *application.TransferCalculationService
	payoutService   *application.PayoutService
}

// NewServer constructs an HTTP server with the provided services.
func NewServer(transferService *application.TransferCalculationService, payoutService *application.PayoutService) *Server {
	return &Server{transferService: transferService, payoutService: payoutService}
}

// Router exposes the configured routes.
func (s *Server) Router() stdhttp.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/transfer/quote", s.handleTransferQuote)
	r.Post("/payouts", s.handlePayout)

	return r
}

type transferQuoteRequest struct {
	SourceCurrency string          `json:"sourceCurrency"`
	TargetCurrency string          `json:"targetCurrency"`
	Amount         decimal.Decimal `json:"amount"`
}

type moneyResponse struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

type transferQuoteResponse struct {
	CanTransact   bool          `json:"canTransact"`
	ExchangeRate  string        `json:"exchangeRate"`
	SendAmount    moneyResponse `json:"sendAmount"`
	ReceiveAmount moneyResponse `json:"receiveAmount"`
	Fee           moneyResponse `json:"fee"`
	TotalDebit    moneyResponse `json:"totalDebit"`
}

func (s *Server) handleTransferQuote(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req transferQuoteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, err)
		return
	}

	quote, err := s.transferService.Calculate(r.Context(), application.TransferCalculationInput{
		SourceCurrency: req.SourceCurrency,
		TargetCurrency: req.TargetCurrency,
		Amount:         req.Amount,
	})
	if err != nil {
		writeError(w, stdhttp.StatusBadRequest, err)
		return
	}

	resp := transferQuoteResponse{
		CanTransact:   quote.Allowed,
		ExchangeRate:  quote.ExchangeRate.String(),
		SendAmount:    toMoneyResponse(quote.SendAmount, 2),
		ReceiveAmount: toMoneyResponse(quote.ReceiveAmount, 2),
		Fee:           toMoneyResponse(quote.Fee, 2),
		TotalDebit:    toMoneyResponse(quote.TotalDebit, 2),
	}

	writeJSON(w, stdhttp.StatusOK, resp)
}

type payoutRequest struct {
	WalletID    string          `json:"walletId"`
	Amount      decimal.Decimal `json:"amount"`
	Currency    string          `json:"currency"`
	Reference   string          `json:"reference"`
	Destination struct {
		Type          string `json:"type"`
		AccountNumber string `json:"accountNumber"`
		BankCode      string `json:"bankCode"`
		BankName      string `json:"bankName"`
		Country       string `json:"country"`
		Currency      string `json:"currency"`
		RecipientName string `json:"recipientName"`
	} `json:"destination"`
}

type payoutResponse struct {
	PayoutID         string                `json:"payoutId"`
	Quote            transferQuoteResponse `json:"quote"`
	RemainingBalance moneyResponse         `json:"remainingBalance"`
}

func (s *Server) handlePayout(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	var req payoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, stdhttp.StatusBadRequest, err)
		return
	}

	result, err := s.payoutService.InitiatePayout(r.Context(), application.PayoutInput{
		WalletID:       req.WalletID,
		Amount:         req.Amount,
		SourceCurrency: req.Currency,
		TargetCurrency: req.Destination.Currency,
		Destination: domain.PayoutDestination{
			Type:          req.Destination.Type,
			AccountNumber: req.Destination.AccountNumber,
			BankCode:      req.Destination.BankCode,
			BankName:      req.Destination.BankName,
			Country:       req.Destination.Country,
			Currency:      req.Destination.Currency,
			RecipientName: req.Destination.RecipientName,
		},
		Reference: req.Reference,
	})
	if err != nil {
		writeError(w, stdhttp.StatusBadRequest, err)
		return
	}

	resp := payoutResponse{
		PayoutID: result.PayoutID,
		Quote: transferQuoteResponse{
			CanTransact:   result.Quote.Allowed,
			ExchangeRate:  result.Quote.ExchangeRate.String(),
			SendAmount:    toMoneyResponse(result.Quote.SendAmount, 2),
			ReceiveAmount: toMoneyResponse(result.Quote.ReceiveAmount, 2),
			Fee:           toMoneyResponse(result.Quote.Fee, 2),
			TotalDebit:    toMoneyResponse(result.Quote.TotalDebit, 2),
		},
		RemainingBalance: toMoneyResponse(result.RemainingBalance, 2),
	}

	writeJSON(w, stdhttp.StatusAccepted, resp)
}

func toMoneyResponse(m domain.Money, precision int32) moneyResponse {
	return moneyResponse{Currency: m.Currency, Amount: m.Amount.StringFixed(precision)}
}

func decodeJSON(r *stdhttp.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body is required")
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return errors.New("request body is empty")
	}
	if err := json.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w stdhttp.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
