package util

import (
	"context"
	"errors"
)

func GetContextValue(ctx context.Context, key any) (string, error) {
	contextValue, ok := ctx.Value(key).(string)
	if !ok || contextValue == "" {
		return "", errors.New("Invalid context key")
	}
	return contextValue, nil
}
