package helpers

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

//ConnectAndReturnSessionDB Crea una conexion con cassandra y returna la session para su posterior uso en consultas hacia cassandra
func ConnectAndReturnSessionDB() *gocql.Session {
	// connect to the cluster
	cluster := gocql.NewCluster("192.168.20.102", "192.168.20.103")
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 0
//	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = 60 * time.Second // Low deadline for testing purposes
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: "api_surnet", Password: "$@CassandraSurnetAPI246678yy01*"}
	session, err := cluster.CreateSession()
	if err != nil {
		ValidError(err)
		log.Fatal()
	}
	return session
}
