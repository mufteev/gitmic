package api

import (
	"context"
	"fmt"
	"gitmic/api/users"
)

// https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28#get-a-user
func (ga *GitApi) GetUserByLogin(ctx context.Context, login string) (*users.User, error) {
	// Получаем подготовленный HTTP-запрос по указанным параметрам
	req, err := users.MakeRequest(ctx, ga.Host, login)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	// Добавляем к запросу токен
	ga.prepareRequestToken(req)

	// Здесь определяем переменную, где структура users.User будет определена со значениями по умолчанию
	var user users.User

	// Выполняем запрос результат которого запишется в структуру
	if _, err := doRequest(req, &user); err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return &user, nil
}
