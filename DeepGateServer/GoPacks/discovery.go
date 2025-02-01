// discovery/discovery.go
package discovery

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type ServiceDiscovery struct {
	ServiceName string
	Port        int
	Domain      string
	TTL         uint32
}

func NewServiceDiscovery(serviceName string, port int) *ServiceDiscovery {
	return &ServiceDiscovery{
		ServiceName: serviceName,
		Port:        port,
		Domain:      "local",
		TTL:         300,
	}
}

// Get local IP address
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no valid local IP address found")
}

// Broadcast service availability using mDNS
func (sd *ServiceDiscovery) Advertise(ctx context.Context) error {
	localIP, err := GetLocalIP()
	if err != nil {
		return err
	}

	fmt.Printf("Advertising service %s on IP %s\n", sd.ServiceName, localIP)

	mDNSServer := &dns.Server{
		Addr: ":5353",
		Net:  "udp",
	}

	dns.HandleFunc(sd.Domain+".", func(w dns.ResponseWriter, r *dns.Msg) {
		msg := new(dns.Msg)
		msg.SetReply(r)

		if len(r.Question) > 0 && strings.HasPrefix(r.Question[0].Name, sd.ServiceName) {
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   r.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    sd.TTL,
				},
				A: net.ParseIP(localIP),
			}
			msg.Answer = append(msg.Answer, rr)
			fmt.Printf("Responding to query for %s with IP %s\n", r.Question[0].Name, localIP)
		}
		w.WriteMsg(msg)
	})

	go func() {
		fmt.Println("Starting mDNS server...")
		if err := mDNSServer.ListenAndServe(); err != nil {
			fmt.Printf("mDNS server error: %v\n", err)
		}
	}()

	return nil
}

// Discover services
func (sd *ServiceDiscovery) Discover(ctx context.Context) ([]string, error) {
	var services []string
	// client := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(sd.ServiceName+"."+sd.Domain+".", dns.TypeA)
	m.RecursionDesired = true

	// Send multicast DNS query
	c, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	addr, err := net.ResolveUDPAddr("udp4", "224.0.0.251:5353")
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(time.Second * 2)
	c.SetReadDeadline(deadline)

	msg, err := m.Pack()
	if err != nil {
		return nil, err
	}
	if _, err := c.WriteTo(msg, addr); err != nil {
		return nil, err
	}

	fmt.Printf("Sent mDNS query for service %s\n", sd.ServiceName)

	buf := make([]byte, 1500)
	for time.Now().Before(deadline) {
		n, _, err := c.ReadFrom(buf)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				break
			}
			continue
		}

		var msg dns.Msg
		if err := msg.Unpack(buf[:n]); err != nil {
			continue
		}

		for _, answer := range msg.Answer {
			if a, ok := answer.(*dns.A); ok {
				services = append(services, a.A.String())
				fmt.Printf("Discovered service at IP %s\n", a.A.String())
			}
		}
	}

	return services, nil
}
