package logger

import "github.com/google/uuid"

func generateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return u.String()
}
