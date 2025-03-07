package account_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/VladimirKholomyanskyy/gym-api/internal/account"
	openapi "github.com/VladimirKholomyanskyy/gym-api/internal/api/go"
	"github.com/VladimirKholomyanskyy/gym-api/internal/common"
	"github.com/golang-migrate/migrate"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ProfileRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository account.ProfileRepository
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	ctx        context.Context
}

func TestProfileRepository(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	fmt.Println("Current working directory:", wd)
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProfileRepositoryTestSuite))
}

func (s *ProfileRepositoryTestSuite) SetupSuite() {
	var err error
	s.ctx = context.Background()

	// Setup docker
	s.pool, err = dockertest.NewPool("")
	require.NoError(s.T(), err, "Could not connect to Docker")

	s.pool.MaxWait = time.Minute * 2

	// Start PostgreSQL container
	s.resource, err = s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=testdb",
			"listen_addresses='*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	require.NoError(s.T(), err, "Could not start PostgreSQL container")

	// Get host and port
	hostAndPort := fmt.Sprintf(
		"host=%s port=%s user=postgres password=postgres dbname=testdb sslmode=disable",
		s.resource.GetBoundIP("5432/tcp"),
		s.resource.GetPort("5432/tcp"),
	)

	// Establish connection
	var sqlDB *sql.DB
	err = s.pool.Retry(func() error {
		var err error
		sqlDB, err = sql.Open("postgres", hostAndPort)
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	})
	require.NoError(s.T(), err, "Could not connect to PostgreSQL container")

	// Create GORM DB instance
	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(s.T(), err, "Could not create GORM DB instance")

	// Create schema
	dsn := "postgres://postgres:postgres@localhost:5432/testdb?sslmode=disable"
	err = s.createSchema(dsn)
	require.NoError(s.T(), err, "Could not create schema")

	// Create repository
	s.repository = account.NewProfileRepository(s.db)
}

func (s *ProfileRepositoryTestSuite) createSchema(dsn string) error {

	m, err := migrate.New(
		"file://C:/Development/GoProjects/gym-api/migrations", // Path to your migrations folder
		dsn,
	)

	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	// Apply all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

func (s *ProfileRepositoryTestSuite) TearDownSuite() {
	// Close DB connection
	sqlDB, err := s.db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}

	// Remove container
	if s.resource != nil {
		_ = s.pool.Purge(s.resource)
	}
}

func (s *ProfileRepositoryTestSuite) TestCreate() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.MALE
	now := time.Now()
	weight := 75.5
	height := 1.82
	avatarURL := "https://example.com/avatar.png"

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
		Birthday:   &now,
		Weight:     &weight,
		Height:     &height,
		AvatarURL:  &avatarURL,
	}

	// Execute
	err := s.repository.Create(s.ctx, profile)

	// Verify
	require.NoError(s.T(), err)
	assert.NotEmpty(s.T(), profile.ID, "Profile ID should be set")
	assert.NotZero(s.T(), profile.CreatedAt, "CreatedAt should be set")
	assert.NotZero(s.T(), profile.UpdatedAt, "UpdatedAt should be set")
}

func (s *ProfileRepositoryTestSuite) TestGetByID() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.FEMALE
	now := time.Now()
	weight := 65.0
	height := 1.70
	avatarURL := "https://example.com/avatar2.png"

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
		Birthday:   &now,
		Weight:     &weight,
		Height:     &height,
		AvatarURL:  &avatarURL,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Execute
	foundProfile, err := s.repository.GetByID(s.ctx, profile.ID)

	// Verify
	require.NoError(s.T(), err)
	assert.Equal(s.T(), profile.ID, foundProfile.ID)
	assert.Equal(s.T(), profile.ExternalID, foundProfile.ExternalID)
	assert.Equal(s.T(), *profile.Sex, *foundProfile.Sex)
	assert.Equal(s.T(), *profile.Weight, *foundProfile.Weight)
	assert.Equal(s.T(), *profile.Height, *foundProfile.Height)
	assert.Equal(s.T(), *profile.AvatarURL, *foundProfile.AvatarURL)
}

func (s *ProfileRepositoryTestSuite) TestGetByID_NotFound() {
	// Execute
	foundProfile, err := s.repository.GetByID(s.ctx, uuid.New().String())

	// Verify
	require.Error(s.T(), err)
	assert.Nil(s.T(), foundProfile)

	// Test with custom error type if implemented
	_, ok := err.(common.ErrEntityNotFound)
	assert.True(s.T(), ok, "Error should be ErrProfileNotFound")
}

func (s *ProfileRepositoryTestSuite) TestFindByExternalID() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.OTHER
	now := time.Now()
	weight := 70.0
	height := 1.75

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
		Birthday:   &now,
		Weight:     &weight,
		Height:     &height,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Execute
	foundProfile, err := s.repository.FindByExternalID(s.ctx, externalID)

	// Verify
	require.NoError(s.T(), err)
	assert.Equal(s.T(), profile.ID, foundProfile.ID)
	assert.Equal(s.T(), profile.ExternalID, foundProfile.ExternalID)
}

func (s *ProfileRepositoryTestSuite) TestUpdate() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.MALE
	now := time.Now()
	weight := 80.0
	height := 1.85

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
		Birthday:   &now,
		Weight:     &weight,
		Height:     &height,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Update
	newWeight := 82.5
	profile.Weight = &newWeight
	newHeight := 1.86
	profile.Height = &newHeight

	// Execute
	err = s.repository.Update(s.ctx, profile)

	// Verify
	require.NoError(s.T(), err)

	// Retrieve updated profile
	updatedProfile, err := s.repository.GetByID(s.ctx, profile.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), *profile.Weight, *updatedProfile.Weight)
	assert.Equal(s.T(), *profile.Height, *updatedProfile.Height)
}

func (s *ProfileRepositoryTestSuite) TestUpdatePartial() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.FEMALE
	now := time.Now()
	weight := 55.0
	height := 1.65

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
		Birthday:   &now,
		Weight:     &weight,
		Height:     &height,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Update
	newWeight := 57.5
	updates := map[string]interface{}{
		"weight": newWeight,
	}

	// Execute
	err = s.repository.UpdatePartial(s.ctx, profile.ID, updates)

	// Verify
	require.NoError(s.T(), err)

	// Retrieve updated profile
	updatedProfile, err := s.repository.GetByID(s.ctx, profile.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), newWeight, *updatedProfile.Weight)
	assert.Equal(s.T(), *profile.Height, *updatedProfile.Height) // Height should be unchanged
}

func (s *ProfileRepositoryTestSuite) TestDelete() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.MALE

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Execute
	err = s.repository.Delete(s.ctx, profile.ID)

	// Verify
	require.NoError(s.T(), err)

	// Try to retrieve deleted profile
	foundProfile, err := s.repository.GetByID(s.ctx, profile.ID)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), foundProfile)
}

func (s *ProfileRepositoryTestSuite) TestRestore() {
	// Setup
	externalID := uuid.New().String()
	sex := openapi.FEMALE

	profile := &account.Profile{
		ExternalID: externalID,
		Sex:        &sex,
	}

	err := s.repository.Create(s.ctx, profile)
	require.NoError(s.T(), err)

	// Delete
	err = s.repository.Delete(s.ctx, profile.ID)
	require.NoError(s.T(), err)

	// Execute restore
	err = s.repository.Restore(s.ctx, profile.ID)

	// Verify
	require.NoError(s.T(), err)

	// Try to retrieve restored profile
	foundProfile, err := s.repository.GetByID(s.ctx, profile.ID)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), foundProfile)
	assert.Equal(s.T(), profile.ID, foundProfile.ID)
}

func (s *ProfileRepositoryTestSuite) TestList() {
	// Clear existing data for this test
	s.db.Exec("DELETE FROM profiles WHERE external_id LIKE 'list-test-%'")

	// Setup - create multiple profiles
	for i := 0; i < 10; i++ {
		profile := &account.Profile{
			ExternalID: fmt.Sprintf("list-test-%d", i),
		}
		err := s.repository.Create(s.ctx, profile)
		require.NoError(s.T(), err)
	}

	// Execute
	profiles, count, err := s.repository.List(s.ctx, 5, 0)

	// Verify
	require.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), count, int64(10))
	assert.Len(s.T(), profiles, 5)

	// Test pagination
	secondPage, _, err := s.repository.List(s.ctx, 5, 5)
	require.NoError(s.T(), err)
	assert.Len(s.T(), secondPage, 5)

	// Make sure the pages are different
	assert.NotEqual(s.T(), profiles[0].ID, secondPage[0].ID)
}

func (s *ProfileRepositoryTestSuite) TestBatchCreate() {
	// Setup
	profiles := make([]*account.Profile, 0, 3)

	for i := 0; i < 3; i++ {
		sex := openapi.MALE
		weight := 70.0 + float64(i*5)

		profiles = append(profiles, &account.Profile{
			ExternalID: fmt.Sprintf("batch-test-%d", i),
			Sex:        &sex,
			Weight:     &weight,
		})
	}

	// Execute
	err := s.repository.BatchCreate(s.ctx, profiles)

	// Verify
	require.NoError(s.T(), err)

	// Check if all profiles have IDs assigned
	for _, profile := range profiles {
		assert.NotEmpty(s.T(), profile.ID)
	}

	// Verify all were created
	for _, profile := range profiles {
		found, err := s.repository.GetByID(s.ctx, profile.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), profile.ExternalID, found.ExternalID)
	}
}
