package voucher

import (
	"encoding/json"
	"fmt"
)

type Service struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FromJson(str string) (*Service, error) {
	err := json.Unmarshal([]byte(str), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Service) String() string {
	return fmt.Sprintf(`{"address":"%s","port":"%d"}`,
		s.Address,
		s.Port)
}

func (s *Service) Remote() string {
	return fmt.Sprintf("%s:%d", s.Address, s.Port)
}
