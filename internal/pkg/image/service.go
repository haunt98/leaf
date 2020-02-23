package image

import (
	"github.com/rs/xid"
)

type Service struct {
	StatusRepo      *StatusRepository
	ProcessProducer *ProcessProducer
}

func (s *Service) Receive(req ReceiveRequest) (Response, error) {
	// Generate UUID
	guid := xid.New()

	// Create status
	if err := s.StatusRepo.Set(guid.String(), Status{
		Status:      Processing,
		OriginalURL: req.URL,
	}); err != nil {
		return Response{}, err
	}

	// Send process message
	if err := s.ProcessProducer.SendProcess(ProcessMessage{
		UUID: guid.String(),
		URL:  req.URL,
	}, Topic); err != nil {
		return Response{}, err
	}

	return Response{
		Status: Successful,
		UUID:   guid.String(),
	}, nil
}
