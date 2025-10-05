package converter

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

func NullString(v null.String) *string {
	if !v.Valid {
		return nil
	}

	return &v.String
}

func NullStringValue(v null.String) string {
	if !v.Valid {
		return ""
	}

	return v.String
}

func NullInt(v null.Int) *int64 {
	if !v.Valid {
		return nil
	}

	return &v.Int64
}

func NullIntValue(v null.Int) int64 {
	if !v.Valid {
		return 0
	}

	return v.Int64
}

func NullFloat(v null.Float) *float64 {
	if !v.Valid {
		return nil
	}

	return &v.Float64
}

func NullFloatValue(v null.Float) float64 {
	if !v.Valid {
		return 0
	}

	return v.Float64
}

func NullBool(v null.Bool) *bool {
	if !v.Valid {
		return nil
	}

	return &v.Bool
}

func NullBoolValue(v null.Bool) bool {
	if !v.Valid {
		return false
	}

	return v.Bool
}

func NullTime(v null.Time) *time.Time {
	if !v.Valid {
		return nil
	}

	return &v.Time
}

func NullTimeValue(v null.Time) time.Time {
	if !v.Valid {
		return time.Time{}
	}

	return v.Time
}

func NullUUID(v uuid.NullUUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}

	return &v.UUID
}

func NullUUIDValue(v uuid.NullUUID) uuid.UUID {
	if !v.Valid {
		return uuid.Nil
	}

	return v.UUID
}

func NullUUIDToString(v uuid.NullUUID) *string {
	if !v.Valid {
		return nil
	}

	return String(v.UUID.String())
}

func NullUUIDToStringValue(v uuid.NullUUID) string {
	if !v.Valid {
		return ""
	}

	return v.UUID.String()
}
