package db

import (
	"context"
	"github.com/slavik22/bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomImage(t *testing.T) Image {
	user := createRandomUser(t)

	arg := CreateImageParams{
		UserID:    user.ID,
		ImagePath: util.RandomString(10),
		ImageUrl:  util.RandomString(10),
	}

	image, err := testStore.CreateImage(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserID, image.UserID)
	require.Equal(t, arg.ImagePath, image.ImagePath)
	require.Equal(t, arg.ImageUrl, image.ImageUrl)

	return image
}

func TestCreateImage(t *testing.T) {
	createRandomImage(t)
}

func TestGetImages(t *testing.T) {
	var img Image

	for i := 0; i < 10; i++ {
		img = createRandomImage(t)
	}

	images, err := testStore.GetImages(context.Background(), img.UserID)

	require.NoError(t, err)
	require.NotEmpty(t, images)

	for _, item := range images {
		require.NotEmpty(t, item)
		require.Equal(t, item.UserID, img.UserID)
	}
}
