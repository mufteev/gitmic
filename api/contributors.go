package api

import (
	"fmt"
	"gitmic/api/repos/contributors"
)

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-contributors
func (ga *GitApi) GetContributorsByRepo(repo string) (*[]*contributors.Contributor, error) {
	// Получаем подготовленный HTTP-запрос по указанным параметрам
	req, err := contributors.MakeRequest(ga.Host, repo, 1, 10)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	// Добавляем к запросу токен
	ga.prepareRequestToken(req)

	// Инициируем исходный тип результата - массив с указателями на участников
	//  с нулевым количеством
	contribs := make([]*contributors.Contributor, 0)

	// Выполняем запрос результат которого запишется в слайс
	if err := doRequest(req, &contribs); err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return &contribs, nil
}
