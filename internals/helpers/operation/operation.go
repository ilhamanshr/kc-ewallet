package operation

import "context"

type operationName struct{}

func SetOperationName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, operationName{}, name)
}

func GetOperationName(ctx context.Context) string {
	val, _ := ctx.Value(operationName{}).(string)
	return val
}
