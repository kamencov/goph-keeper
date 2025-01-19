package service

import (
	"flag"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// Flags - структура флагов.
type Flags struct {
	log          *slog.Logger
	AddrGRPC     string
	TokenSalt    []byte
	PasswordSalt []byte
}

// NewFlags - конструктор флагов.
func NewFlags(log *slog.Logger) *Flags {
	return &Flags{
		log: log,
	}
}

// Parse - парсим флаги.
func (f *Flags) Parse() {
	f.parsFlags()
	f.initSaltFromEnv()
}

// parsFlags - парсим флаги.
func (f *Flags) parsFlags() {
	flag.StringVar(&f.AddrGRPC, "addr", ":8081", "gRPC address")
}

// initSaltFromEnv - инициализация соли.
func (f *Flags) initSaltFromEnv() {
	err := godotenv.Load("/Users/pavel/GolandProjects/goph-keeper/.env")
	if err != nil {
		f.log.Error("Fatal", "error loading .env file = ", err)
		return
	}

	f.TokenSalt = []byte(os.Getenv("TOKEN_SALT"))
	f.PasswordSalt = []byte(os.Getenv("PASSWORD_SALT"))
}
