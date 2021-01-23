package config

import (
	"bytes"
	"fmt"
	"text/template"
)

type ConnectionStringer interface {
	ConnectionString() string
}

type DBVendor int8

const (
	VendorUnknown DBVendor = iota
	VendorPostgres
)

var (
	ErrUnknownDBVendor = func(in string) error {
		return fmt.Errorf("unknown db vendor `%s`", in)
	}
	vendorDisplay = []string{"unknown", "postgres"}
	vendorLookup  = map[string]DBVendor{
		"postgres": VendorPostgres,
		"unknown":  VendorUnknown,
	}
	pgConnStrTmpl     = `user={{.User}} password={{.Password}} host={{.Host}} port={{.Port}} dbname={{.DBName}}{{range $k, $v := .Args}} {{$k}}={{$v}}{{end}}`
	pgConnStrTmplName = "pg"
	tmpl              *template.Template
)

func init() {
	tmpl = template.Must(template.New(pgConnStrTmplName).Parse(pgConnStrTmpl))
}

func (d DBVendor) String() string {
	return vendorDisplay[d]
}

func ParseDBVendor(in string) (DBVendor, error) {
	vendor, ok := vendorLookup[in]
	if !ok {
		return VendorUnknown, ErrUnknownDBVendor(in)
	}

	return vendor, nil
}

func (d *DBVendor) UnmarshalText(data []byte) error {
	vendor, err := ParseDBVendor(string(data))
	if err != nil {
		return err
	}

	*d = vendor

	return nil
}

func ConnectionString(config DB) string {
	var connectionStringer ConnectionStringer

	switch config.vendor {
	case VendorPostgres:
		connectionStringer = newPostgresConnectionStringer(config)
	case VendorUnknown:
		connectionStringer = &unknownConnectionStringer{}
	}

	return connectionStringer.ConnectionString()
}

func DriverName(config DB) string {
	var driver string
	switch config.vendor {
	case VendorPostgres:
		driver = "pgx"
	}
	return driver
}

type postgreConnectionStringer struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	Args     map[string]interface{}
}

func newPostgresConnectionStringer(config DB) ConnectionStringer {
	return &postgreConnectionStringer{
		User:     config.user,
		Password: config.password,
		Host:     config.host,
		Port:     config.port,
		DBName:   config.name,
		Args:     config.args,
	}
}

func (p *postgreConnectionStringer) ConnectionString() string {
	var buf bytes.Buffer

	if err := tmpl.ExecuteTemplate(&buf, pgConnStrTmplName, &p); err != nil {
		panic(err)
	}

	return buf.String()
}

type unknownConnectionStringer struct{}

func (u *unknownConnectionStringer) ConnectionString() string {
	return "unknown"
}
