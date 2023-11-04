package repos

import (
	"fmt"
	"net/http"
	"net/url"
)

// Структура данных Репозитория GitHub
// Прочитать про теги
// https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go-ru
// https://pkg.go.dev/encoding/json#:~:text=The%20encoding%20of%20each%20struct%20field%20can%20be%20customized%20by%20the%20format%20string%20stored
type Repo struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	HtmlUrl  string `json:"html_url"`
	Private  bool   `json:"private"`
}

// Подготовка HTTP-запроса к Репозиториям организации в GitHub
func MakeRequest(host, org string, page, perPage int) (*http.Request, error) {
	// Собираем первоначальный URL
	queryUrl := fmt.Sprintf("%s/orgs/%s/repos", host, org)

	// Парсим через пакет `url` и проверяем его корректность
	u, err := url.Parse(queryUrl)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	// Из URL достаём объект query-параметров
	q := u.Query()
	// Устанавливаем query-параметры
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("per_page", fmt.Sprintf("%d", perPage))

	// Модифицированные query-параметры кодируем и присваиваем обратно к URL
	u.RawQuery = q.Encode()

	// Объект URL приводим к строке - это и будет итоговый URL-запрос
	// В качестве Body нам нечего передать в метод GET
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	return req, nil
}
