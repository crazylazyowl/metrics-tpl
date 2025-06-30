package ping

import "context"

type Pinger interface {
	Ping(ctx context.Context) error
}

type PingUsecase struct {
	pinger Pinger
}

func New(pinger Pinger) *PingUsecase {
	return &PingUsecase{pinger: pinger}
}

func (u *PingUsecase) Ping(ctx context.Context) error {
	if err := u.pinger.Ping(ctx); err != nil {
		return err
	}
	return nil
}
