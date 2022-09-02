package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func TestHospital(t *testing.T) {
	store, cleanup := helperConnect(t)
	defer cleanup()

	ctx := context.Background()

	var err error
	var hospital, hospitalSame, hospitalOther *models.Hospital

	t.Run("CreateHospital", func(t *testing.T) {
		hospital, err = store.CreateHospital(ctx, &dto.Hospital{
			Name:        "foo",
			DisplayName: "foo bar hospital",
		})
		assert.NoError(t, err)
		assert.Equal(t, "foo", hospital.Name)
		assert.Equal(t, "foo bar hospital", hospital.DisplayName)
		assert.Greater(t, hospital.ID, int64(0))
	})

	t.Run("CreateAnotherHospital", func(t *testing.T) {
		hospitalOther, err = store.CreateHospital(ctx, &dto.Hospital{
			Name:        "foo-other",
			DisplayName: "foo other hospital",
		})
		assert.NoError(t, err)
		assert.Equal(t, "foo-other", hospitalOther.Name)
		assert.Equal(t, "foo other hospital", hospitalOther.DisplayName)
		assert.Greater(t, hospitalOther.ID, int64(0))
	})

	t.Run("GetHospital", func(t *testing.T) {
		hospitalSame, err = store.GetHospital(ctx, hospital.ID)
		assert.NoError(t, err)
		assert.Equal(t, hospital.Name, hospitalSame.Name)
		assert.Equal(t, hospital.DisplayName, hospitalSame.DisplayName)
		assert.Equal(t, hospital.ID, hospitalSame.ID)
	})

	t.Run("FindHospitals", func(t *testing.T) {
		hospitals, err := store.FindHospitals(ctx, 0, 10)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(hospitals))
		assert.Equal(t, hospital.ID, hospitals[0].ID)
		assert.Equal(t, hospital.Name, hospitals[0].Name)
		assert.Equal(t, hospital.DisplayName, hospitals[0].DisplayName)
		assert.Equal(t, hospital.ID, hospitals[0].ID)

		assert.Equal(t, hospitalOther.ID, hospitals[1].ID)
		assert.Equal(t, hospitalOther.Name, hospitals[1].Name)
		assert.Equal(t, hospitalOther.DisplayName, hospitals[1].DisplayName)
		assert.Equal(t, hospitalOther.ID, hospitals[1].ID)
	})

	t.Run("FindHospitalsWithLimit", func(t *testing.T) {
		total, err := store.CountHosptials(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint(2), total)

		hospitals, err := store.FindHospitals(ctx, 0, 1)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(hospitals))
		assert.Equal(t, hospital.ID, hospitals[0].ID)
		assert.Equal(t, hospital.Name, hospitals[0].Name)
		assert.Equal(t, hospital.DisplayName, hospitals[0].DisplayName)
		assert.Equal(t, hospital.ID, hospitals[0].ID)

		hospitals, err = store.FindHospitals(ctx, 1, 1)
		assert.NoError(t, err)

		assert.Equal(t, hospitalOther.ID, hospitals[0].ID)
		assert.Equal(t, hospitalOther.Name, hospitals[0].Name)
		assert.Equal(t, hospitalOther.DisplayName, hospitals[0].DisplayName)
		assert.Equal(t, hospitalOther.ID, hospitals[0].ID)
	})

	t.Run("CreateHospitalIfExist", func(t *testing.T) {
		_, err = store.CreateHospital(ctx, &dto.Hospital{
			Name:        "foo",
			DisplayName: "foo bar hospital",
		})
		assert.True(t, IsErrDuplicateEntry(err))
	})

	t.Run("UpdateHospital", func(t *testing.T) {
		n, err := store.UpdateHospital(ctx, &dto.Hospital{
			ID:          hospital.ID,
			Name:        "bar",
			DisplayName: hospital.DisplayName,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), n)

		h, err := store.GetHospital(ctx, hospital.ID)
		assert.NoError(t, err)
		assert.Equal(t, "bar", h.Name)
		assert.Equal(t, hospital.DisplayName, h.DisplayName)
	})
}
