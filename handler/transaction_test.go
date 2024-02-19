package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rhuantac/rinha-concurrency/config"
	"github.com/rhuantac/rinha-concurrency/internal"
)

func setupServer() *gin.Engine {
	mongoClient := config.SetupMongo()
	internal.ClearDb(mongoClient)
	internal.SeedDb(mongoClient)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	api := router.Group("/clientes/:id")

	api.POST("/transacoes", TransactionHandler(mongoClient.Database(os.Getenv("MONGO_DATABASE")), config.SetupRedis()))
	api.GET("/extrato", HistoryHandler(mongoClient.Database(os.Getenv("MONGO_DATABASE"))))
	return router
}

func TestTransactions(t *testing.T) {
	godotenv.Load(filepath.Join("../", ".env"))
	t.Run("credits increase amount", func(t *testing.T) {
		router := setupServer()
		accId := 1
		endpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           500,
			TransactionType: "c",
			Description:     "Test transaction",
		})
		wantedResponse := TransactionResponse{
			Limit:   100000,
			Balance: 500,
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request didn't returned the right status code. Got %d expected %d", w.Code, http.StatusOK)
		}

		var got TransactionResponse
		json.NewDecoder(w.Body).Decode(&got)

		if !reflect.DeepEqual(wantedResponse, got) {
			t.Errorf("got %v want %v", got, wantedResponse)
		}

	})

	t.Run("debits decrease amount", func(t *testing.T) {
		router := setupServer()
		accId := 1
		endpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           500,
			TransactionType: "d",
			Description:     "Test transaction",
		})
		wantedResponse := TransactionResponse{
			Limit:   100000,
			Balance: -500,
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request didn't returned the right status code. Got %d expected %d", w.Code, http.StatusOK)
		}

		var got TransactionResponse
		json.NewDecoder(w.Body).Decode(&got)

		if !reflect.DeepEqual(wantedResponse, got) {
			t.Errorf("got %v want %v", got, wantedResponse)
		}

	})

	t.Run("debits won't exceed limit", func(t *testing.T) {
		router := setupServer()
		accId := 1
		endpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           100001,
			TransactionType: "d",
			Description:     "Test transaction",
		})

		wantedCode := http.StatusUnprocessableEntity

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(w, req)

		if w.Code != wantedCode {
			t.Errorf("request didn't returned the right status code. Got %d expected %d", w.Code, wantedCode)
		}
	})

	t.Run("unknown user returns 404", func(t *testing.T) {
		router := setupServer()
		accId := 999
		endpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
		testPayload, _ := json.Marshal(TransactionRequest{
			Value:           100001,
			TransactionType: "d",
			Description:     "Test transaction",
		})

		wantedCode := http.StatusNotFound

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(w, req)

		if w.Code != wantedCode {
			t.Errorf("request didn't returned the right status code. Got %d expected %d", w.Code, wantedCode)
		}
	})
}

func BenchmarkTransactions(b *testing.B) {
	router := setupServer()
	accId := 1
	endpoint := fmt.Sprintf("/clientes/%d/transacoes", accId)
	testPayload, _ := json.Marshal(TransactionRequest{
		Value:           500,
		TransactionType: "d",
		Description:     "Test transaction",
	})

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(string(testPayload)))
		router.ServeHTTP(w, req)
	}	
}
