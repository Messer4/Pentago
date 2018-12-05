package db

import (
	"github.com/go-redis/redis"
    "github.com/jmoiron/sqlx"
    "fmt"
    "github.com/wedancedalot/squirrel"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "regexp"
)

var ErrDuplicate = fmt.Errorf("DB duplicate error")
var ErrNoRows = fmt.Errorf("DB no rows in resultset")
var errorsRegexp = regexp.MustCompile(`^Error (?P<code>\d+)`)

type DI struct {
    redisClient *redis.Client
    mysqlCli *Mysql
}

type Mysql struct {
    Db        *sqlx.DB
    DebugMode bool
}

var d *DI


func Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

    db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", "root", "localhost", 3306, "pentago"))
    if err != nil {
        panic(err)
    }

    db.SetMaxOpenConns(10)
    d = &DI{
        mysqlCli:&Mysql{db,false},
        redisClient:client,
    }
}

func GetClientRedis() *redis.Client {
	return d.redisClient
}

func GetClientMysql() *Mysql {
    return d.mysqlCli
}

// Insert from querybuilder
func (this *Mysql) Insert(b squirrel.InsertBuilder, tx ...*sqlx.Tx) (uint64, error) {
    q, args, err := b.ToSql()
    if err != nil {
        return 0, err
    }

    result, err := this.exec(q, args, tx...)
    if err != nil {
        return 0, this.parseError(err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return uint64(id), err
}

// Exec query
func (this *Mysql) Exec(q string, args []interface{}, err error, tx ...*sqlx.Tx) (uint64, error) {
    if err != nil {
        return 0, this.parseError(err)
    }

    result, err := this.exec(q, args, tx...)
    if err != nil {
        return 0, this.parseError(err)
    }

    affectedRows, err := result.RowsAffected()
    if err != nil {
        return 0, this.parseError(err)
    }

    return uint64(affectedRows), nil
}

// Find first row into dest from querybuilder
func (this *Mysql) FindFirst(dest interface{}, b squirrel.SelectBuilder, tx ...*sqlx.Tx) (err error) {
    q, params, err := b.ToSql()
    if err != nil {
        return
    }

   if len(tx) > 0 && tx[0] != nil {
        err = tx[0].Get(dest, q, params...)
    } else {
        err = this.Db.Get(dest, q, params...)
    }

    return this.parseError(err)
}

func (this *Mysql) exec(q string, params []interface{}, tx ...*sqlx.Tx) (sql.Result, error) {
    var result sql.Result
    var err error

    if len(tx) > 0 && tx[0] != nil {
        result, err = tx[0].Exec(q, params...)
    } else {
        result, err = this.Db.Exec(q, params...)
    }

    return result, err
}

func (this *Mysql) parseError(err error) error {
    if err == nil {
        return nil
    }

    // Just a wrapper not to use sql lib directly from code
    if err == sql.ErrNoRows {
        return ErrNoRows
    }

    matches := this.matchStringGroups(errorsRegexp, err.Error())
    code, ok := matches["code"]
    if !ok {
        return err
    }

    switch code {
    case "1062":
        return ErrDuplicate
    default:
        return err
    }
}

// matchStringGroups matches regexp with capture groups. Returns map string string
func (this *Mysql) matchStringGroups(re *regexp.Regexp, s string) map[string]string {
    m := re.FindStringSubmatch(s)
    n := re.SubexpNames()

    r := make(map[string]string, len(m))

    if len(m) > 0 {
        m, n = m[1:], n[1:]
        for i, _ := range n {
            r[n[i]] = m[i]
        }
    }

    return r
}
