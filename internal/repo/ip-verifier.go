package repo

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type IPVerifierRepo struct {
	db *geoip2.Reader
}

func NewIPVerifierRepo(db *geoip2.Reader) *IPVerifierRepo {
	return &IPVerifierRepo{
		db: db,
	}
}

func (r *IPVerifierRepo) GetCountryByIP(ipAddress string) (string, error) {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", ipAddress)
	}

	record, err := r.db.Country(ip)
	if err != nil {
		return "", err
	}

	return record.Country.IsoCode, nil
}
