package config

import (
	"github.com/hashicorp/errwrap"
)

type DB struct {
	vendor   DBVendor
	user     string
	password string
	host     string
	port     int
	name     string
	args     map[string]interface{}
}

func (d *DB) UnmarshalTOML(data interface{}) (err error) {
	dataMap := data.(map[string]interface{})

	if vendor, ok := dataMap["vendor"].(string); ok && vendor != "" {
		parsed, vendorErr := ParseDBVendor(vendor)
		if vendorErr != nil {
			err = errwrap.Wrap(vendorErr, err)
		} else {
			d.vendor = parsed
		}
	}

	if user, ok := dataMap["user"].(string); ok && user != "" {
		d.user = user
	}

	if password, ok := dataMap["password"].(string); ok && password != "" {
		d.password = password
	}

	if host, ok := dataMap["host"].(string); ok && host != "" {
		d.host = host
	}

	if name, ok := dataMap["name"].(string); ok && name != "" {
		d.name = name
	}

	if port, ok := dataMap["port"].(int64); ok && port > 0 {
		d.port = int(port)
	}

	if args, ok := dataMap["args"].(map[string]interface{}); ok {
		d.args = args
	}

	return err
}
