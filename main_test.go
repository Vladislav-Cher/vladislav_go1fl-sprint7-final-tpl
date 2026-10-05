package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := "/cafe?city=moscow&count=" + strconv.Itoa(v.count)
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		bodyString := strings.TrimSpace(response.Body.String())
		bodySlice := strings.Split(bodyString, ",")
		if bodyString == "" {
			bodySlice = []string{}
		}

		if v.count > 0 && bodySlice == nil {
			t.Error("response body is empty")
			return
		}

		assert.Equal(t, v.want, len(bodySlice))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := "/cafe?city=moscow&search=" + v.search
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		bodyString := strings.TrimSpace(strings.ToLower(response.Body.String()))
		bodySlice := strings.Split(bodyString, ",")
		if bodyString == "" {
			bodySlice = []string{}
		}
		if bodySlice == nil {
			t.Error("response body is empty")
			return
		}

		for _, s := range bodySlice {
			assert.Contains(t, s, v.search)
		}
	}
}
