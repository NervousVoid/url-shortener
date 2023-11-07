package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shrotener/internal/storage"
	"url-shrotener/tools"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	aliasLen   = 10
	aliasAlpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_"
)

var saveCases = []struct {
	name         string
	url          string
	addDuplicate bool
	respError    string
	dbErrors     []error
}{
	{
		name: "Success",
		url:  "https://ozon.ru",
	},
	{
		name:      "Invalid url",
		url:       "htps:/www.ozon.ru",
		respError: "field URL is not a valid URL",
	},
	{
		name:      "Empty url",
		url:       "",
		respError: "field URL is a required field",
	},
	{
		name:         "Duplicate",
		addDuplicate: true,
		url:          "https://ozon.ru",
		dbErrors:     []error{storage.ErrURLExists},
	},
}

func TestSaveHandlerInMem(t *testing.T) {
	for _, tc := range saveCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := storage.NewInMemStorage()

			if tc.addDuplicate {
				s.Mu.Lock()
				s.Urls[tc.url] = "testalias"
				s.Aliases["testalias"] = tc.url
				s.Mu.Unlock()
			}

			handler := Save(tools.NewMockLogger(), s, aliasLen, aliasAlpha)
			input := fmt.Sprintf(`{"url": "%s"}`, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/api/http", bytes.NewReader([]byte(input)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			var resp ResponseAlias
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}

func TestSaveHandlerDatabase(t *testing.T) {
	for _, tc := range saveCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("cant create mock: %s", err)
			}
			defer db.Close()

			s := storage.Database{DB: db}

			if tc.addDuplicate {
				mock.ExpectPrepare("^INSERT INTO urls\\(full_url, short_url\\) VALUES\\(\\$1, \\$2\\)$").
					ExpectExec().
					WillReturnError(tc.dbErrors[0])

				rows := sqlmock.NewRows([]string{"short_url"})
				rows = rows.AddRow("123abc")
				mock.ExpectPrepare("^SELECT short_url FROM urls WHERE full_url = \\$1 LIMIT 1").
					ExpectQuery().
					WillReturnRows(rows)

			} else {
				mock.ExpectPrepare("^INSERT INTO urls\\(full_url, short_url\\) VALUES\\(\\$1, \\$2\\)$").
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			handler := Save(tools.NewMockLogger(), &s, aliasLen, aliasAlpha)
			input := fmt.Sprintf(`{"url": "%s"}`, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/api/http", bytes.NewReader([]byte(input)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			var resp ResponseAlias
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
