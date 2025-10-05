package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/IamOnah/storefronthq/internal/app/auth"
	"github.com/IamOnah/storefronthq/internal/config"
	"github.com/IamOnah/storefronthq/internal/domain/users"
	mockdb "github.com/IamOnah/storefronthq/internal/domain/users/userstore/mock"
	"github.com/IamOnah/storefronthq/internal/infra/database"
	"github.com/IamOnah/storefronthq/internal/sdk/authz"
	"github.com/IamOnah/storefronthq/internal/sdk/jobs"
	"github.com/IamOnah/storefronthq/internal/transport/http/router"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	number := "123456789"
	createUser := users.UserCreate{
		Password:    &number,
		FirstName:   "victor",
		LastName:    "onah",
		Email:       "onahvictor@gamil.com",
		PhoneNumber: "08164528525",
		Country:     "Nigeria",
	}

	user, err := users.NewUser(createUser)
	require.NoError(t, err)

	type TestCases []struct {
		name          string
		body          any
		buildStubs         func(repo *mockdb.MockUserRepository)
		checkResponse func(r *httptest.ResponseRecorder)
	}

	var testCases = TestCases{
		{
			name: "OK: User Created",
			body: map[string]string{
				"fist_name":    "victor",
				"last_name":    "onah",
				"password":     "123456789",
				"phone_number": "08164528525",
				"email":        "onahvictor@gmail.com",
				"country":      "Nigeria",
			},

			//check text input parameter
			buildStubs: func(repo *mockdb.MockUserRepository) {
				repo.EXPECT().CreateUser(gomock.Any(), user).Times(1).Return(&user, nil)
			},
			checkResponse: func(r *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, r.Code)
			},
		},
	}

	for _, value := range testCases {
		t.Run(value.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg, err := config.LoadConfig(".")
			require.NoError(t, err)
			log := zerolog.New(os.Stdout)

			userService, err := auth.NewUserService(
				auth.WithAuth(authz.NewJWTMaker(cfg.Auth.TokenSymmetricKey)),
				auth.WithUserRepository(mockdb.NewMockUserRepository(ctrl)),
				auth.WithTrxManager(database.NewTRXManager(&pgxpool.Pool{}, &log)),
				auth.WithJob(jobs.NewJobClient(cfg.Redis, &log)),
				auth.WithLog(&log),
			)
			require.NoError(t, err)

			router := router.SetupRouter(userService)
			data, err := json.Marshal(value.body)
			require.NoError(t, err)

			reader := bytes.NewReader(data)
			req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			value.buildStubs(mockdb.NewMockUserRepository(ctrl))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			value.checkResponse(rec)

		})
	}

}
