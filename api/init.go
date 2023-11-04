package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// Основной объект API, который будет содержать в себе адрес хоста и токен
type GitApi struct {
	Host  string
	token *string
}

// Хост по умолчанию
const defaultHost = "https://api.github.com"

// Функция-конструктор для объекта API с указанием списка параметров
func NewGitApi(opts ...ApiOption) (*GitApi, error) {
	// Начальное объявление GitApi
	ga := &GitApi{
		Host: defaultHost,
	}

	// Проход по списку параметров, которые реализуют паттерн функциональных опций
	// https://habr.com/ru/articles/575316
	// https://www.sohamkamani.com/golang/options-pattern
	for _, opt := range opts {
		if err := opt(ga); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	// Если токен не был указан в опции, его указатель пустой - возвращаем ошибку
	if ga.token == nil {
		return nil, fmt.Errorf("expected token")
	}

	// Всё ок
	return ga, nil
}

// По указателю получаем Запрос и добавляем к нему заголовок "Authorization: Bearer <token>"
func (ga *GitApi) prepareRequestToken(r *http.Request) {
	bearerToken := fmt.Sprintf("Bearer %s", *ga.token)

	r.Header.Set("Authorization", bearerToken)
}

// Выполняем типовой запрос, большинство запросов шаблонные
//	в таком случае можно указать `initResponse` тип interface{}/any (C# - object, Pascal - Variant)
//  так как параметр был передан по указателю - его можно не возвращать из этой функции, изменения будут видны в исходной точке

func doRequest(r *http.Request, initResponse interface{}) error {
	// Проверяем через пакет reflect тип данных упакованный в interface{}/any
	rv := reflect.ValueOf(initResponse)
	// Если тип не указатель - возвращаем ошибку
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("init response not pointer")
	}

	// Выполняем HTTP-запрос через стандартный клиент (можно реализовывать кастомный)
	httpResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("client do request: %w", err)
	}
	// Прочитать про `defer` - https://go.dev/ref/spec#Defer_statements
	// Body реализует интерфейс `io.ReadCloser` - такой следует всегда закрывать, после использования
	// Проэкспериментируйте и удалите слово defer
	defer httpResponse.Body.Close()

	// Считываем всё тело результата в байты
	// (в дальнейшем следует относится осторожнее, т.к. тело может быть большим, а это пишется в ОЗУ)
	buff, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Декодируем байты в структуру или массив структур с готовыми json-тегами
	if err := json.Unmarshal(buff, initResponse); err != nil {
		return fmt.Errorf("response body unmarshal json: %w (%s)", err, string(buff))
	}

	return nil
}
