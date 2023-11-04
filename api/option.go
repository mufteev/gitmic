package api

import (
	"fmt"
	"net/url"
)

// Реализация паттерна функциональных опций

// https://habr.com/ru/articles/575316
// https://www.sohamkamani.com/golang/options-pattern

// Сокращённый тип функции, которая принимает указатель на GitApi
// Внутри эта функция должна преобразовать исходный объект
// Если возникает ошибка - вернуть её
type ApiOption func(*GitApi) error

// Опция указания хоста
func WithAnotherHost(host string) ApiOption {
	return func(ga *GitApi) error {
		// Прочитать про пустой идентификатор ( _ ) - https://go.dev/ref/spec#Blank_identifier
		_, err := url.Parse(host)
		if err != nil {
			return fmt.Errorf("host url not parsed: %w", err)
		}

		ga.Host = host

		return nil
	}
}

// Опция указания токена
func WithGitToken(token string) ApiOption {
	return func(ga *GitApi) error {
		// Передаём указатель на полученную переменную-токен
		ga.token = &token

		return nil
	}
}
