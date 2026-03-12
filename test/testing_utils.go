package tests

/**
 * This file contains helper functions for setting up a test MongoDB database and inserting test data.
 */

import (
	"context"
	"os"
	"testing"

	"github.com/afuradanime/backend/cmd/api/app"
	"github.com/stretchr/testify/require"
)

func SetupTestApp(t *testing.T) (*app.Application, func()) {

	// Setup os
	// Starting the application with the "test" environment, which should use a separate MongoDB database for testing.
	os.Setenv("ENV", "test")
	os.Chdir("../../") // Change the working directory to the root of the project to ensure .env.test is loaded correctly

	application := app.New()

	cleanup := func() {
		err := application.Mongo.Client().Database(application.Config.MongoDatabase).Drop(context.Background())
		require.NoError(t, err)
	}

	return application, cleanup
}
