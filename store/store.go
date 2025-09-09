package store

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

// Init ...
func Init() {

	var database Database
	dbDriver := viper.GetString("DB_DRIVER")

	switch dbDriver {
	case "mysql":
		database = &MysqlDB{}
	case "postgres":
		database = &PostgresDB{}
	default:
		panic(fmt.Sprintf("database driver not supported: %s", dbDriver))
	}

	database.GenerateDSN(viper.GetString("DB_NAME"))
	err := database.ConnectORM()

	if err != nil {
		panic(err)
	}

	dbConn = database.GetGorm()
}

// Conn struct holds the store connection
type Conn struct {
	DB           *gorm.DB
	IncidentConn Incident
	SLOConn      SLO
	// TODO: add other connection
}

// NewStore inits new store connection
func NewStore() *Conn {
	Init()
	conn := &Conn{
		DB: dbConn,
	}
	conn.IncidentConn = NewIncidentStore(conn)
	conn.SLOConn = NewSLOStore(conn)
	// TODO: Add other connections

	return conn
}

// Incident implements the store interface and it returns the Incident interface
func (s *Conn) Incident() Incident {
	return s.IncidentConn
}

// SLO implements the store interface and it returns the SLO interface
func (s *Conn) SLO() SLO {
	return s.SLOConn
}

func getCommonIndexes(tableName string) map[string]string {
	idx := fmt.Sprintf("idx_%s", tableName)
	return map[string]string{
		fmt.Sprintf("%s_created_at", idx): "created_at",
		fmt.Sprintf("%s_updated_at", idx): "updated_at",
	}
}

// recordExists should check if record is avail or not for particular table
// based on the given condition.
func recordExists(tableName, where string) (exists bool) {
	baseQ := fmt.Sprintf("select 1 from %s where %v", tableName, where)
	dbConn.Raw(fmt.Sprintf("select exists (%v)", baseQ)).Row().Scan(&exists)
	return
}
