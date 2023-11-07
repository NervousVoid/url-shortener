package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shrotener/internal/storage"
	"url-shrotener/tools"
)

var getCases = []struct {
	name        string
	aliasURL    string
	aliasValue  string
	aliasAbsent bool
	respError   string
	dbErrors    []error
}{
	{
		name:       "Alias Exists",
		aliasURL:   "https://domain.com/alias",
		aliasValue: "alias",
	},
	{
		name:        "No Alias",
		aliasURL:    "https://domain.com/alias",
		aliasValue:  "alias",
		respError:   "not found",
		aliasAbsent: true,
		dbErrors:    []error{storage.ErrURLNotFound},
	},
}

func TestGetHandlerInMem(t *testing.T) {
	for _, tc := range getCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := storage.NewInMemStorage()

			if !tc.aliasAbsent {
				s.Mu.Lock()
				s.Aliases[tc.aliasValue] = "testurl"
				s.Urls["testurl"] = tc.aliasValue
				s.Mu.Unlock()
			}

			handler := Get(tools.NewMockLogger(), s)
			input := fmt.Sprintf(`{"alias": "%s"}`, tc.aliasURL)

			req, err := http.NewRequest(http.MethodGet, "/api/http", bytes.NewReader([]byte(input)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			var resp ResponseURL
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}

func TestGetHandlerDatabase(t *testing.T) {
	for _, tc := range getCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("cant create mock: %s", err)
			}
			defer db.Close()

			s := storage.Database{DB: db}

			if tc.aliasAbsent {
				mock.ExpectPrepare("^SELECT full_url FROM urls WHERE short_url = \\$1 LIMIT 1").
					ExpectQuery().
					WillReturnError(tc.dbErrors[0])
			} else {
				rows := sqlmock.NewRows([]string{"full_url"})
				rows = rows.AddRow("https://ozon.ru")
				mock.ExpectPrepare("^SELECT full_url FROM urls WHERE short_url = \\$1 LIMIT 1").
					ExpectQuery().
					WillReturnRows(rows)
			}

			handler := Get(tools.NewMockLogger(), &s)
			input := fmt.Sprintf(`{"alias": "%s"}`, tc.aliasURL)

			req, err := http.NewRequest(http.MethodGet, "/api/http", bytes.NewReader([]byte(input)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			var resp ResponseURL
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
