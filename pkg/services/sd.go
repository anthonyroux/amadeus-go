/*
	SD implements Service Discovery using **Consul** (you can implement different SD later if you feel like it)
*/

package services

import (
	"log"
	"time"

	consul "github.com/hashicorp/consul/api"
)

type serviceReg struct {
	Name        string
	TTL         time.Duration
	ConsulAgent *consul.Agent
}

func registerService(addr string, port int, ttl time.Duration) (*serviceReg, error) {
	s := serviceReg{
		Name: addr,
		TTL:  ttl,
	}

	c, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}

	s.ConsulAgent = c.Agent()

	serviceDef := consul.AgentServiceRegistration{
		Name: s.Name,
		// TODO: change this address later
		Address: "localhost",
		Port:    port,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}
	if err := s.ConsulAgent.ServiceRegister(&serviceDef); err != nil {
		return nil, err
	}

	go s.updateTTL()

	return &s, nil
}

func (s *serviceReg) updateTTL() {
	ticker := time.NewTicker(s.TTL / 2)
	for range ticker.C {
		if agentErr := s.ConsulAgent.FailTTL("service:"+s.Name, ""); agentErr != nil {
			log.Println(agentErr)
		}

		if agentErr := s.ConsulAgent.PassTTL("service:"+s.Name, ""); agentErr != nil {
			log.Fatalln(agentErr)
		}
	}
}
