package image

import (
	"github.com/rs/xid"
)

type Service struct {
	StatusRepo      *StatusRepository
	ProcessProducer *ProcessProducer
}

func (s *Service) Receive(req ReceiveRequest) (ReceiveResponse, error) {
	// Generate UUID
	guid := xid.New()

	// Create status
	if err := s.StatusRepo.Create(guid.String(), Status{
		Status:      Processing,
		OriginalURL: req.URL,
	}); err != nil {
		return ReceiveResponse{}, nil
	}

	// Send process message
	if err := s.ProcessProducer.SendProcess(ProcessMessage{
		UUID: guid.String(),
		URL:  req.URL,
	}, Topic); err != nil {
		return ReceiveResponse{}, nil
	}

	return ReceiveResponse{
		UUID: guid.String(),
	}, nil
}
