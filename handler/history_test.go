package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestHistory(t *testing.T) {
	godotenv.Load(filepath.Join("../", ".env"))

	t.Run("balance is returned correctly", func(t *testing.T) {
		router := setupServer()
		accId := 1
		txEndpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           500,
			TransactionType: "c",
			Description:     "Test transaction",
		})
		req, _ := http.NewRequest(http.MethodPost, txEndpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(httptest.NewRecorder(), req)

		historyEndpoint := fmt.Sprintf("/clientes/%d/extrato", accId)
		want := HistoryResponse{
			Balance: Balance{
				Total: 500,
				Limit: 100000,
			},
		}
		w := httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, historyEndpoint, nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request didn't return the right status code. Got %d expected %d", w.Code, http.StatusOK)
		}

		var got HistoryResponse
		json.NewDecoder(w.Body).Decode(&got)

		if got.Balance.Limit != want.Balance.Limit || got.Balance.Total != want.Balance.Total {
			t.Errorf("request didn't return balance correctly. Got %v expected %v", got.Balance, want.Balance)
		}

	})

	t.Run("transactions are returned correctly", func(t *testing.T) {
		router := setupServer()
		accId := 1
		txEndpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           500,
			TransactionType: "c",
			Description:     "Test transaction",
		})

		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest(http.MethodPost, txEndpoint, strings.NewReader(string(testPayload)))
			router.ServeHTTP(httptest.NewRecorder(), req)
		}

		historyEndpoint := fmt.Sprintf("/clientes/%d/extrato", accId)
		wantedTransactions := 5
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, historyEndpoint, nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request didn't return the right status code. Got %d expected %d", w.Code, http.StatusOK)
		}

		var got HistoryResponse
		json.NewDecoder(w.Body).Decode(&got)

		if len(got.Transactions) != wantedTransactions {
			t.Errorf("request didn't return transactions correctly. Got %d expected %d", len(got.Transactions), wantedTransactions)
		}
	})
}
