package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const fileEnvGit = "./secret/.git.env"

var GIT_TOKEN string

func init() {
	// Загружаем переменные среды из файла пакетом github.com/joho/godotenv
	if err := godotenv.Load(fileEnvGit); err != nil {
		log.Fatalf("load env from file `%s`: %v", fileEnvGit, err)
	}

	// Из переменных сред достаём по названию значение токена
	// Переменной присваиваем значение по указателю
	loadStrVar(&GIT_TOKEN, "GIT_TOKEN")
}

func loadStrVar(variable *string, name string) {
	var ok bool
	// Присваивание значения переменной по указателю
	// 	*(pointer) = значение переменной
	// семантика такого присваивания заставляет определять переменную Ok отдельно
	if *variable, ok = os.LookupEnv(name); !ok {
		log.Fatalf("`%s` environment not declare", name)
	}
}
