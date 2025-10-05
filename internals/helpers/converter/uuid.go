package converter

import "github.com/google/uuid"

func MustParseUUIDSliceFromStringSlice(uuidStrings []string) []uuid.UUID {
	if uuidStrings == nil {
		return nil
	}

	uuids := make([]uuid.UUID, len(uuidStrings))
	for i, uuidString := range uuidStrings {
		uuids[i] = uuid.MustParse(uuidString)
	}

	return uuids
}

func UUIDSliceFromStringSlice(uuidStrings []string) (uuids []uuid.UUID, err error) {
	if uuidStrings == nil {
		return nil, nil
	}

	uuids = make([]uuid.UUID, len(uuidStrings))
	for i, uuidString := range uuidStrings {
		uuids[i], err = uuid.Parse(uuidString)
		if err != nil {
			return nil, err
		}
	}

	return uuids, nil
}
