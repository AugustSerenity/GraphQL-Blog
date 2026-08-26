package service

type PostSaver interface {
	SavePost()
}

type Service struct {
	ps PostSaver
}

func NewService(ps PostSaver) *Service {
	return &Service{
		ps: ps,
	}
}
