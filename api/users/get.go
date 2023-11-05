package users

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
)

// Структура данных Пользователя GitHub
// Прочитать про теги
// https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go-ru
// https://pkg.go.dev/encoding/json#:~:text=The%20encoding%20of%20each%20struct%20field%20can%20be%20customized%20by%20the%20format%20string%20stored
type User struct {
	Id       int     `json:"id"`
	Login    string  `json:"login"`
	Name     *string `json:"name"`
	Location *string `json:"location"`
	Bio      *string `json:"bio"`
	Company  *string `json:"company"`
}

// Подготовка HTTP-запроса к Пользователю GitHub
func MakeRequest(ctx context.Context, host, login string) (*http.Request, error) {
	// Собираем первоначальный URL
	queryUrl := fmt.Sprintf("%s/users/%s", host, login)

	// Парсим через пакет `url` и проверяем его корректность
	u, err := url.Parse(queryUrl)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	// Объект URL приводим к строке - это и будет итоговый URL-запрос
	// В качестве Body нам нечего передать в метод GET
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	return req, nil
}

// Добавляем интерактивность для структуры данных User
// Определяем функцию, которая будет доступна для объекта User
//
//	будет возвращать описание о пользователе в формате
//	<login> (name=<Name>, loc=<Location>, bio=<Bio>, comp=<Company>)
//	Если дополнительных параметров нет, будет выведен логин без скобок
func (u *User) GetDescribe() string {
	var title = color.YellowString("%s", u.Login)

	// Резервируем массив на 4 строки
	descriptions := make([]string, 4)
	// Инициируем счётчик, определяющий последнее добавленное
	cnt := 0

	if u.Name != nil {
		// description[cnt] = fmt.Sprintf("name=%s", *u.Name)
		// Форматируем цвет вывода через пакет https://github.com/fatih/color
		descriptions[cnt] = color.BlueString("name=%s", *u.Name)
		cnt++
	}
	if u.Location != nil {
		// description[cnt] = fmt.Sprintf("loc=%s", *u.Location)
		descriptions[cnt] = color.GreenString("loc=%s", *u.Location)
		cnt++
	}
	if u.Bio != nil {
		// description[cnt] = fmt.Sprintf("bio=%s", *u.Bio)
		descriptions[cnt] = color.MagentaString("bio=%s", *u.Bio)
		cnt++
	}
	if u.Company != nil {
		// description[cnt] = fmt.Sprintf("comp=%s", *u.Company)
		descriptions[cnt] = color.CyanString("comp=%s", *u.Company)
		cnt++
	}

	// Если счётчик больше 0
	if cnt > 0 {
		// Нарезаем слайс до последнего добавленного параметра (чтобы не выводились пустые запятые)
		descriptions = descriptions[:cnt]
		// Собираем массив строк в строку по разделителю
		// Почитать про функционал пакета strings - https://pkg.go.dev/strings
		title = fmt.Sprintf("%s (%s)", title, strings.Join(descriptions, ", "))
	}

	return title
}
