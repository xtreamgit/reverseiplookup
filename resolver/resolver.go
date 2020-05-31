package resolver

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Resolver struct {
	storage ResolverStorage
}

func NewResolver(stor ResolverStorage) Resolver {
	return Resolver{storage: stor}
}

type Domains struct {
	Domain string    `json:"domain"`
	Valid  time.Time `json:"validFrom"`
}

// IPLookup takes a ip string and returns domains assocatiated.
func (r *Resolver) IPLookup(ip string) (domains []Domains, err error) {
	var errCount = 0
	addr, err := net.LookupAddr(ip)
	if err != nil {
		errCount++
	}

	for _, a := range addr {
		clean := strings.TrimSuffix(a, ".")
		domains = append(domains, Domains{Domain: clean, Valid: time.Now()})
	}

	//do a merge of the domains for newer ones from the net lookup
	dbList, err := r.storage.GetByIP(ip)
	if err != nil {
		errCount++
	}

	for _, domain := range dbList {
		domains = append(domains, Domains{Domain: domain.Domain, Valid: domain.Valid})
	}

	if errCount == 2 {
		return nil, fmt.Errorf("IPLookup error - no returned domains")
	}

	return unique(domains), err
}

// HostLookup takes a domain string and returns ips assocatiated.
func (r *Resolver) HostLookup(domain string) (domains []Domains, err error) {
	addr, err := net.LookupHost(domain)
	if err != nil {
		return nil, fmt.Errorf("HostLookup error - no returned domains")
	}
	domains = append(domains, Domains{Domain: domain, Valid: time.Now()})

	for _, a := range addr {
		doms, err := r.IPLookup(a)
		if err != nil {
			log.Println(err)
			continue
		}

		domains = append(domains, doms...)
		err = r.storage.InsertOrUpdate(DomainRecord{Domain: domain, IP: a, Valid: time.Now()})
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return unique(domains), err
}

// UpdateValid checks and updates valid field for host
func (r *Resolver) UpdateValid(limit int) error {
	fmt.Println("*** Starting Update ***")
	addr, err := r.storage.GetOldest(limit)
	if err != nil {
		return err
	}

	for _, a := range addr {
		doms, err := r.IPLookup(a.IP)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, d := range doms {
			if a.Domain == d.Domain {
				a.Valid = time.Now()
				err = r.storage.Update(a.ID, a)
				if err != nil {
					log.Println(err)
				}
				continue
			}

			err = r.storage.InsertOrUpdate(DomainRecord{Domain: d.Domain, IP: a.IP, Valid: time.Now()})
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}

	return err
}

// Make sure more recent entries go first
func unique(domSlice []Domains) []Domains {
	keys := make(map[string]bool)
	list := []Domains{}
	for _, entry := range domSlice {
		if _, value := keys[entry.Domain]; !value {
			keys[entry.Domain] = true
			list = append(list, entry)
		}
	}
	return list
}
