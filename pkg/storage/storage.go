package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}
func checkDb(conn *sqlx.DB){
	conn.Exec("CREATE TABLE IF NOT EXISTS tokens(guid char(36) primary key, token varchar(255) not null, ip varchar(45) not null, email varchar(255))")
	log.Print("table was created")
}
func NewDbConn(cfg Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	checkDb(conn)
	return conn, nil
}

func SaveRefreshToken(guid string, refToken string, clientIp string) error {

	db, err := NewDbConn(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		log.Fatal("failed to initialize db: ", err)
	}

	_, err = db.Exec("INSERT INTO tokens(guid, token, ip, email) VALUES($1, $2, $3, '123456@mail.ru')", guid, refToken, clientIp)
	if err != nil {
		return err
	}

	return nil
}

type TokenInfo struct {
	Guid        string
	HashedToken string
	ClientIp    string
	Email       string
}

func GetRefreshTokenInfo(guid string) (*TokenInfo, error) {
	db, err := NewDbConn(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	defer db.Close()

	ti := new(TokenInfo)

	if err != nil {
		log.Fatal("failed to initialize db: ", err)
	}

	_ = db.QueryRow("SELECT guid, token, ip, email FROM tokens WHERE guid=$1", guid).Scan(&ti.Guid, &ti.HashedToken, &ti.ClientIp, &ti.Email)
	return ti, nil
}
