package main

import (
    "flag"
    "bufio"
    "encoding/binary"
    "os"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "gopkg.in/ini.v1"
    "log"
)

type Config struct {
    Driver string
    Host string
    Port string
    User string
    Pass string
    Dbname string
    Dbargs string
    Table string
    UserField string
    PassField string
    ServerField string
}

func auth(conf *Config, db *sql.DB, user string,
    server string, passwd string) (bool, error) {
    var query string
    var value string
    var err error

    log.Printf("user: %s, server: %s, passwd: %s\n", user, server, passwd)

    if conf.ServerField != "" {
        query = fmt.Sprintf(
            "select %s from %s where %s = $1 and %s = $2 and %s = $3",
            conf.UserField, conf.Table, conf.UserField, conf.PassField,
            conf.ServerField)
        err = db.QueryRow(query, user, passwd, server).Scan(&value)
    } else {
        query = fmt.Sprintf(
            "select %s from %s where %s = $1 and %s = $2",
            conf.UserField, conf.Table, conf.UserField, conf.PassField)
        err = db.QueryRow(query, user, passwd).Scan(&value)
    }
    if err != nil || value == "" {
        return false, err
    }

    return true, nil
}

func isuser(conf *Config, db *sql.DB, user string,
    server string) (bool, error) {
    var query string
    var value string
    var err error

    if conf.ServerField != "" {
        query = fmt.Sprintf(
            "select %s from %s where %s = $1 and %s = $2",
            conf.UserField, conf.Table, conf.UserField, conf.ServerField)
        err = db.QueryRow(query, user, server).Scan(&value)
    } else {
        query = fmt.Sprintf(
            "select %s from %s where %s = $1",
            conf.UserField, conf.Table, conf.UserField)
        err = db.QueryRow(query, user).Scan(&value)
    }
    if err != nil || value == "" {
        return false, err
    }

    return true, nil
}

func GetSqlConnectionString(conf *Config) string {
    return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
        conf.Driver, conf.User, conf.Pass, conf.Host, conf.Port,
        conf.Dbname, conf.Dbargs)
}

func OpenSqlConnection(conf *Config) (*sql.DB, error) {
    var err error

    connectionString := GetSqlConnectionString(conf)
    db, err := sql.Open(conf.Driver, connectionString)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}

func AuthLoop(conf *Config) {
    db, err := OpenSqlConnection(conf)

    bioIn := bufio.NewReader(os.Stdin)
    bioOut := bufio.NewWriter(os.Stdout)

    var success bool
    var length uint16
    var result uint16

    for {
        _ = binary.Read(bioIn, binary.BigEndian, &length)

        buf := make([]byte, length)

        r, _ := bioIn.Read(buf)
        if r == 0 {
            continue
        }

        if err != nil {
            err = db.Ping()
        }

        if err == nil {
            data := strings.Split(string(buf), ":")
            if data[0] == "auth" {
                success, err = auth(conf, db, data[1], data[2], data[3])
            } else if data[0] == "isuser" {
                success, err = isuser(conf, db, data[1], data[2])
            } else {
                success = false
            }
        } else {
            success = false
        }

        length = 2
        binary.Write(bioOut, binary.BigEndian, &length)

        if success != true {
            result = 0
        } else {
            result = 1
        }

        binary.Write(bioOut, binary.BigEndian, &result)
        bioOut.Flush()
    }
}

func main() {
    filename := flag.String("conf", "/etc/ejabberd-pg-auth.ini",
        "Config file with all the connection infos needed.")
    flag.Parse()

    cfg, err := ini.Load(*filename)
    if err != nil {
        log.Fatal(err)
    }

    conf := new(Config)
    err = cfg.MapTo(conf)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(os.Stderr)

    AuthLoop(conf)
}
