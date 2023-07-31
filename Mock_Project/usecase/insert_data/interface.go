package insert_data

import (
	"Mock_Project/model"
	"context"
)

type IDB interface {
	StartDBProcess(ctx context.Context, collection *[]model.ObjectProcess) error
}
