package watcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	preChangeData = []byte(`
[info]
name = "agora"
version = "1.2.3"
env = "qa"

[api]
port = 9123
pathPrefix = "prefix"

[logging]
level = "trace"
outputPaths = ["stdout", "/var/log/app.log"]
[logging.fields]
from = "toml"

[db]
vendor = "postgres"
user = "test"
password = "test"
host = "localhost"
port = 6000
name = "test"
[db.args]
sslmode = "disable"
`)
	postChangeData = []byte(`
[info]
name = "agora"
version = "1.2.4"
env = "prod"

[api]
port = 9124
pathPrefix = "differentPrefix"

[logging]
level = "trace"
outputPaths = ["stdout", "/var/log/app.log"]
[logging.fields]
from = "toml"

[db]
vendor = "postgres"
user = "test"
password = "test"
host = "localhost"
port = 6000
name = "test"
[db.args]
sslmode = "disable"
`)
)

func TestWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watcher Suite")
}
