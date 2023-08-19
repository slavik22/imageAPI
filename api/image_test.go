package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	mockdb "github.com/slavik22/imageAPI/db/mock"
	db "github.com/slavik22/imageAPI/db/sqlc"
	"github.com/slavik22/imageAPI/token"
	"github.com/slavik22/imageAPI/util"
	"github.com/stretchr/testify/require"
)

var (
	fileName = "NYgf3QgEdUpHTR7YIYacJanBU3JEeDxmIKGOUKcD.jpg"
	imgUrl   = "https://fingers-site-production.s3.eu-central-1.amazonaws.com/uploads/images/" + fileName
)

func TestCreateImageAPI(t *testing.T) {
	user, _ := randomUser(t)
	image := randomImage(user.ID)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"image_url": image.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateImage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(image, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchImage(t, recorder.Body, image)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"image_url": image.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateImage(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"image_url": image.ImageUrl,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateImage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Image{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUrl",
			body: gin.H{
				"image_url": "googlecom",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateImage(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/upload-picture"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.jwtMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListImagesAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	images := make([]db.Image, n)
	for i := 0; i < n; i++ {
		images[i] = randomImage(user.ID)
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {

				store.EXPECT().
					GetImages(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(images, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchImages(t, recorder.Body, images)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetImages(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.JWTMaker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.ID, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetImages(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Image{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/images"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.jwtMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomImage(userId int64) db.Image {

	return db.Image{
		ID:       util.RandomInt(1, 1000),
		UserID:   userId,
		ImageUrl: imgUrl,
	}
}

func requireBodyMatchImage(t *testing.T, body *bytes.Buffer, image db.Image) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotImage db.Image
	err = json.Unmarshal(data, &gotImage)

	path, _ := filepath.Abs("../images")
	path += "/" + fileName

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		t.Error(err)
	}

	require.NoError(t, err)
	require.Equal(t, image.ImageUrl, gotImage.ImageUrl)
	require.Equal(t, image.UserID, gotImage.UserID)
}

func requireBodyMatchImages(t *testing.T, body *bytes.Buffer, images []db.Image) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotImages []db.Image
	err = json.Unmarshal(data, &gotImages)
	require.NoError(t, err)
	require.Equal(t, images, gotImages)
}
