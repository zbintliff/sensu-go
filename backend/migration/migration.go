package migration

import (
	"github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
)

var logger = logrus.WithFields(logrus.Fields{
	"component": "migration",
})

// Run lauches the migration process
func Run(client *clientv3.Client) {
	environments(client)
}
