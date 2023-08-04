package fetch_db

import (
	"Mock_Project/model"
	"context"
)

type IFetchDB interface {
	StartFetchDB(ctx context.Context, collection *[]model.ObjectProcess) error
}
