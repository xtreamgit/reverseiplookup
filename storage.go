package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "reverseiplookup/resolver"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

var create = `
    CREATE DATABASE IF NOT EXISTS domains;
`

var schema = `
CREATE TABLE IF NOT EXISTS records (
    id INT(11) NOT NULL AUTO_INCREMENT,
    domain VARCHAR(255) NOT NULL,
    ip VARCHAR(16) NOT NULL,
    valid DATETIME NOT NULL,
    PRIMARY KEY (id) USING BTREE
)
`

type Storage struct {
    Record resolver.DomainRecord
    DB     *sqlx.DB
}

// SetupDB brings back the storage object
func SetupDB() Storage {
    fmt.Println(os.Getenv("DB_CONN"))
    db, err := sqlx.Connect("mysql", os.Getenv("DB_CONN"))
    if err != nil {
        log.Fatalln(err)
    }

    db.MustExec(create)
    db.MustExec(schema)

    return Storage{DB: db}
}

// GetByIP ...
func (s Storage) GetByIP(IP string) ([]resolver.DomainRecord, error) {
    rec := []resolver.DomainRecord{}
    err := s.DB.Select(&rec, "SELECT * FROM records WHERE ip = ?", IP)
    return rec, err
}

// GetOldest ...
func (s Storage) GetOldest(limit int) ([]resolver.DomainRecord, error) {
    rec := []resolver.DomainRecord{}
    err := s.DB.Select(&rec, "SELECT * FROM records ORDER BY `valid` ASC LIMIT ?", limit)
    return rec, err
}

// Insert ...
func (s Storage) Insert(dr resolver.DomainRecord) (ID int64, err error) {
    r, err := s.DB.NamedExec(`INSERT INTO records (domain,ip,valid) VALUES (:domain,:ip,:valid)`,
        map[string]interface{}{
            "domain": dr.Domain,
            "ip":     dr.IP,
            "valid":  time.Now(),
        })
    if err != nil {
        return 0, err
    }
    id, _ := r.LastInsertId()

    return id, nil
}

// InsertOrUpdate ...
func (s Storage) InsertOrUpdate(dr resolver.DomainRecord) (err error) {
    res, err := s.DB.NamedExec(`UPDATE records SET valid=:valid WHERE domain=:domain AND ip=:ip`,
        map[string]interface{}{
            "domain": dr.Domain,
            "ip":     dr.IP,
            "valid":  dr.Valid,
        })
    if err != nil && err != sql.ErrNoRows {
        log.Println(err, "hello")
        return err
    }

    rows, _ := res.RowsAffected()
    if rows == 0 || err == sql.ErrNoRows {
        _, err := s.DB.NamedExec(`INSERT INTO records (domain,ip,valid) VALUES (:domain,:ip,:valid)`,
            map[string]interface{}{
                "domain": dr.Domain,
                "ip":     dr.IP,
                "valid":  dr.Valid,
            })
        if err != nil {
            return err
        }
        return nil
    }

    return nil
}

//Update ...
func (s Storage) Update(ID int64, record resolver.DomainRecord) (err error) {
    _, err = s.DB.NamedExec(`UPDATE records SET domain=:domain, ip=:ip, valid=:valid WHERE id=:id`,
        map[string]interface{}{
            "id":     ID,
            "domain": record.Domain,
            "ip":     record.IP,
            "valid":  record.Valid,
        })
    if err != nil {
        return err
    }

    return err
}
