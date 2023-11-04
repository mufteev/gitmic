package contributors

import (
	"fmt"
	"net/http"
	"net/url"
)

// Структура данных Участника GitHub
// Прочитать про теги
// https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go-ru
// https://pkg.go.dev/encoding/json#:~:text=The%20encoding%20of%20each%20struct%20field%20can%20be%20customized%20by%20the%20format%20string%20stored

type Contributor struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

// Подготовка HTTP-запроса к Участникам репозитория в GitHub
func MakeRequest(host, fullRepo string, page, perPage int) (*http.Request, error) {
	// Собираем первоначальный URL
	queryUrl := fmt.Sprintf("%s/repos/%s/contributors", host, fullRepo)

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
