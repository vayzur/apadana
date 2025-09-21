package register

import (
	"context"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"

	apadana "github.com/vayzur/apadana/pkg/client"
)

type RegisterManager struct {
	apadanaClient *apadana.Client
}

func NewRegisterManager(
	apadanaClient *apadana.Client,
) *RegisterManager {
	return &RegisterManager{
		apadanaClient: apadanaClient,
	}
}

func (r *RegisterManager) RegisterWithAPIServer(ctx context.Context, node *corev1.Node) error {
	step := 100 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(step):
			step = step * 2
			if step >= 7*time.Second {
				step = 7 * time.Second
			}

			_, err := r.apadanaClient.CreateNode(node)
			if err != nil {
				continue
			}
			return nil
		}
	}
}
