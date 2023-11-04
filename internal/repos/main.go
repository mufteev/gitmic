package repos

import (
	"fmt"
	"gitmic/api"
	"gitmic/internal/env"
)

const org = "microsoft"

func Run() error {
	// Получаем объект GitApi с указанным токеном
	ga, err := api.NewGitApi(api.WithGitToken(env.GIT_TOKEN))
	if err != nil {
		return fmt.Errorf("new git api: %w", err)
	}

	// Получаем список репозиториев
	repos, err := ga.GetReposByOrg(org)
	if err != nil {
		return fmt.Errorf("get repos by org `%s`: %w", org, err)
	}

	// Проходим по массиву репозиториев
	// Почитать про указатели - https://www.digitalocean.com/community/conceptual-articles/understanding-pointers-in-go-ru
	for _, repo := range *repos {
		// Выводим текущий репозиторий
		fmt.Printf("repo `%s` contrib users:\n", repo.HtmlUrl)

		// Получаем список участников
		contribs, err := ga.GetContributorsByRepo(repo.FullName)
		if err != nil {
			return fmt.Errorf("get contribs for `%s`: %w", repo.FullName, err)
		}

		// Проходим по массиву участников
		for _, contrib := range *contribs {
			// По логину участника получам информацию о пользователе
			user, err := ga.GetUserByLogin(contrib.Login)
			if err != nil {
				return fmt.Errorf("get user `%s`: %w", contrib.Login, err)
			}

			// Выводим информацию об участнике с пользовательской информацией
			fmt.Printf("\t• %s\n", user.GetDescribe())
		}

	}

	return nil
}
