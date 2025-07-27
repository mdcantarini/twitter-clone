package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

func NewSession(hosts []string, keyspace string) *gocql.Session {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("cannot connect to Cassandra: %v", err)
	}
	return session
}
