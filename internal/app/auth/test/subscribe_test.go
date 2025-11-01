package auth_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/iamonah/merchcore/internal/app/auth"
// 	"github.com/iamonah/merchcore/internal/domain/users"
// 	mockdb "github.com/iamonah/merchcore/internal/domain/users/userdb/mock"
// 	"github.com/iamonah/merchcore/internal/infra/database"
// 	"github.com/iamonah/merchcore/internal/sdk/authz"
// 	"github.com/iamonah/merchcore/internal/sdk/jobs"
// 	"github.com/iamonah/merchcore/internal/transport/http/router"

// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// )

// func TestRegisterUser(t *testing.T) {
// 	number := "123456789"
// 	createUser := users.UserCreate{
// 		Password:    number,
// 		FirstName:   "victor",
// 		LastName:    "onah",
// 		Email:       "onahvictor@gamil.com",
// 		PhoneNumber: "08164528525",
// 		Country:     "Nigeria",
// 	}

// 	user, err := users.NewUser(createUser)
// 	require.NoError(t, err)

// 	type TestCases []struct {
// 		name          string
// 		body          any
// 		buildStubs    func(repo *mockdb.MockUserRepository)
// 		checkResponse func(r *httptest.ResponseRecorder)
// 	}

// 	var testCases = TestCases{
// 		{
// 			name: "OK: User Created",
// 			body: map[string]string{
// 				"fist_name":    "victor",
// 				"last_name":    "onah",
// 				"password":     "123456789",
// 				"phone_number": "08164528525",
// 				"email":        "onahvictor@gmail.com",
// 				"country":      "Nigeria",
// 			},

// 			//check text input parameter
// 			buildStubs: func(repo *mockdb.MockUserRepository) {
// 				repo.EXPECT().CreateUser(gomock.Any(), user).Times(1).Return(&user, nil)
// 			},
// 			checkResponse: func(r *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, r.Code)
// 			},
// 		},
// 	}

// 	for _, value := range testCases {
// 		t.Run(value.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			repo := mockdb.NewMockUserRepository(ctrl)

// 			userService, err := auth.NewUserService(
// 				auth.WithAuth(authz.NewJWTMaker(cfg.Auth.TokenSymmetricKey)),
// 				auth.WithUserRepository(repo),
// 				auth.WithTrxManager(database.NewTRXManager(pool, &log)),
// 				auth.WithJob(jobs.NewJobClient(cfg.Redis, &log)),
// 				auth.WithLog(&log),
// 			)
// 			if err != nil {
// 				t.Fatalf("failed to create user service: %v", err)
// 			}

// 			router := router.SetupRouter(userService, &log)
// 			data, err := json.Marshal(value.body)
// 			require.NoError(t, err)

// 			reader := bytes.NewReader(data)
// 			req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
// 			req.Header.Set("Content-Type", "application/json")
// 			require.NoError(t, err)

// 			value.buildStubs(mockdb.NewMockUserRepository(ctrl))
// 			rec := httptest.NewRecorder()
// 			router.ServeHTTP(rec, req)
// 			value.checkResponse(rec)

// 		})
// 	}

// }
