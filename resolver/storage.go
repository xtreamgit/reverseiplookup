package resolver

import "time"

type DomainRecord struct {
	ID     int64
	Domain string
	IP     string
	Valid  time.Time
}

type ResolverStorage interface {
	GetByIP(IP string) ([]DomainRecord, error)
	GetOldest(limit int) ([]DomainRecord, error)
	Insert(DomainRecord) (ID int64, err error)
	InsertOrUpdate(DomainRecord) error
	Update(ID int64, record DomainRecord) error
}
