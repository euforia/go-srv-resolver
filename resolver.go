package resolver

/* This package implements an interface to perform service record dns lookups.  It
can used to communication with services using service discovery.
*/

import (
	"fmt"

	"github.com/miekg/dns"
)

// ServiceInfo contains information about a single instance of a service
type ServiceInfo struct {
	Hostname string
	Port     uint16
	IP       string
}

// Address returns the service address including port.
func (s *ServiceInfo) Address() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}

// Resolver implements a DNS resolver for SRV records
type Resolver struct {
	config dns.ClientConfig
	client *dns.Client
}

// NewResolver instantiates a new Resolver instance.
func NewResolver(port int, servers ...string) *Resolver {
	r := &Resolver{
		config: dns.ClientConfig{Servers: servers, Port: fmt.Sprintf("%d", port)},
		client: new(dns.Client),
	}

	return r
}

// Lookup performs an SRV lookup against the given name
func (resolver *Resolver) Lookup(name string) ([]ServiceInfo, error) {
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(name), dns.TypeSRV)
	m.SetEdns0(4096, true)
	r, _, err := resolver.client.Exchange(m, resolver.config.Servers[0]+":"+resolver.config.Port)
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("lookup failed: %d", r.Rcode)
	}

	info := map[string]ServiceInfo{}

	for _, k := range r.Extra {
		switch k.(type) {
		case *dns.A:
			a := k.(*dns.A)
			info[a.Hdr.Name] = ServiceInfo{IP: a.A.String(), Hostname: a.Hdr.Name}
		}
	}

	for _, k := range r.Answer {
		if key, ok := k.(*dns.SRV); ok {
			if v, ok := info[key.Target]; ok {
				v.Port = key.Port
				info[key.Target] = v
			}
		}
	}

	out := make([]ServiceInfo, 0, len(info))
	for _, v := range info {
		out = append(out, v)
	}

	return out, nil
}
