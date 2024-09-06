package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	mock_db "github.com/mikeheft/go-backend/db/mock"
	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/mikeheft/go-backend/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fmt.Printf("account.ID type: %T\n", account.ID)

	var accountID int64 = account.ID

	store := mock_db.NewMockStore(ctrl)
	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(accountID)).
		Times(1).
		Return(account, nil)

	// start test server and send req
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", accountID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)
	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
