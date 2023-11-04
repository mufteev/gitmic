package api

import (
	"fmt"
	"gitmic/api/orgs/repos"
)

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories
func (ga *GitApi) GetReposByOrg(org string) (*[]*repos.Repo, error) {
	// Получаем подготовленный HTTP-запрос по указанным параметрам
	req, err := repos.MakeRequest(ga.Host, org, 1, 10)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	// Добавляем к запросу токен
	ga.prepareRequestToken(req)

	// Инициируем исходный тип результата - массив с указателями на репозитории
	//  с нулевым количеством
	repos := make([]*repos.Repo, 0)

	// Выполняем запрос результат которого запишется в слайс
	if err := doRequest(req, &repos); err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return &repos, nil
}
