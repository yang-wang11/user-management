package discovery

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/resolver"
)

type Server struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Weight  int64  `json:"weight"`
}

func BuildRegisteredPath(s *Server) string {
	pathPrefix := fmt.Sprintf("%s/", s.Name)
	if s.Version != "" {
		pathPrefix = fmt.Sprintf("%s/%s/", s.Name, s.Version)
	}
	return fmt.Sprintf("%s%s", pathPrefix, s.Addr)
}

func MarshalRegisteredServer(s *Server) ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalRegisteredServer(value []byte) (*Server, error) {
	s := Server{}
	if err := json.Unmarshal(value, &s); err != nil {
		return &s, err
	}
	return &s, nil
}

func IsRegisteredPathExist(addrs []resolver.Address, addr resolver.Address) bool {
	for _, add := range addrs {
		if add.Equal(addr) {
			return true
		}
	}
	return false
}
