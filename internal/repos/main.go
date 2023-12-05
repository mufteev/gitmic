package repos

import (
	"context"
	"fmt"
	"gitmic/api"
	"gitmic/api/users"
	"gitmic/internal/env"

	"golang.org/x/sync/errgroup"
)

const org = "microsoft"

func RunSimple(isPrintData bool) error {
	ctx := context.TODO()

	// Получаем объект GitApi с указанным токеном
	ga, err := api.NewGitApi(nil, api.WithGitToken(env.GIT_TOKEN))
	if err != nil {
		return fmt.Errorf("new git api: %w", err)
	}

	// Получаем список репозиториев
	repos, err := ga.GetReposByOrg(ctx, org, 2, 10)
	if err != nil {
		return fmt.Errorf("get repos by org `%s`: %w", org, err)
	}

	// Проходим по массиву репозиториев
	// Почитать про указатели - https://www.digitalocean.com/community/conceptual-articles/understanding-pointers-in-go-ru
	for _, repo := range *repos {
		// Выводим текущий репозиторий
		if isPrintData {
			fmt.Printf("repo `%s` contrib users:\n", repo.HtmlUrl)
		}

		// Получаем список участников
		contribs, err := ga.GetContributorsByRepo(ctx, repo.FullName, 2, 10)
		if err != nil {
			return fmt.Errorf("get contribs for `%s`: %w", repo.FullName, err)
		}

		// Проходим по массиву участников
		for _, contrib := range *contribs {
			// По логину участника получам информацию о пользователе
			user, err := ga.GetUserByLogin(ctx, contrib.Login)
			if err != nil {
				return fmt.Errorf("get user `%s`: %w", contrib.Login, err)
			}

			// Выводим информацию об участнике с пользовательской информацией
			if isPrintData {
				fmt.Printf("\t• %s\n", user.GetDescribe())
			}
		}
	}

	return nil
}

func RunConcurrency(isPrintData bool) error {
	ctx := context.Background()

	// Получаем объект GitApi с указанным токеном
	ga, err := api.NewGitApi(nil, api.WithGitToken(env.GIT_TOKEN))
	if err != nil {
		return fmt.Errorf("new git api: %w", err)
	}

	// Получаем список репозиториев
	repos, err := ga.GetReposByOrg(ctx, org, 2, 10)
	if err != nil {
		return fmt.Errorf("get repos by org `%s`: %w", org, err)
	}

	// Создаём группу для горутин, если в какой-либо горутине возникнет ошибка
	// Все остальные должны будут завершить своё выполнение
	// Также обновляем информацию о контексте:
	//   Если хоть одна горутина будет завершена с ошибкой
	//   Весь контекст будет считаться завершённым
	//   Из-за чего обновлённые HTTP-запросы `http.NewRequestWithContext` будут завершаться до завершения HTTP-запроса
	//   Что обеспечит моментальное завершение программы без ожидания ответов на запросы
	g, ctx := errgroup.WithContext(ctx)

	// Инициируем слайс, чтобы на каждый репозиторий был массив пользователей
	repoUsers := make([][]*users.User, len(*repos))

	// Проходим по слайсу репозиториев
	for i := range *repos {
		// Так как дальнейшее выполнение будет происходить асинхронно
		// Необходимо запомнить для текущего контекста выполнения параметр `i`
		// Иначе при отложенном запуске горутины переменная `i` будет перезаписана и горутина может взять значение не в момент вызова, а после следующей итерации
		// По этой же причине мы не будем использовать второй параметр цикла содержащего в себе текущий
		ii := i

		// Запускаем асинхронную горутину при помощи вспомогательного метода
		g.Go(func() error {
			// Получаеем текущий репозиторий
			// Разыменовываем указатель на слайс и получаем доступ к элементу по индексу
			repo := (*repos)[ii]

			// Получаем список Участников
			contribs, err := ga.GetContributorsByRepo(ctx, repo.FullName, 2, 10)
			if err != nil {
				return fmt.Errorf("get contribs for `%s`: %w", repo.FullName, err)
			}

			// Создаём ещё одну асинхронную группу горутин
			g2, ctx := errgroup.WithContext(ctx)
			// Иницируем слайс по количеству Участников = Пользователей
			contribUsers := make([]*users.User, len(*contribs))

			// Проходим по слайсу Участиков
			for i2 := range *contribs {
				// Запоминаем `i2` в текущем контексте исполнения
				ii2 := i2

				// Запускаем асинхронную горутину
				g2.Go(func() error {
					// Получаем текущего Участника
					contrib := (*contribs)[ii2]

					// Получаем информацию о пользователе на текущего участника
					user, err := ga.GetUserByLogin(ctx, contrib.Login)
					if err != nil {
						return fmt.Errorf("get user `%s`: %w", contrib.Login, err)
					}

					// Присваиваем Пользователя по соответствующему индексу Участника
					contribUsers[ii2] = user

					return nil
				})
			}

			// Ожидаем завершения группы горутин
			// если возникает ошибка - завершаем дальнейшее выполнение и передаём её дальше
			if err := g2.Wait(); err != nil {
				return fmt.Errorf("goroutine 2: %w", err)
			}

			// Если список Пользователей успешно загружен для текущего репозитория
			// присваиваем этот слайс конкретному репозиторию по индексу
			repoUsers[ii] = contribUsers

			return nil
		})
	}

	// Ожидаем завершения группы горутин
	if err := g.Wait(); err != nil {
		return fmt.Errorf("goroutine: %w", err)
	}

	if isPrintData {
		for i, users := range repoUsers {
			repo := (*repos)[i]

			fmt.Printf("repo `%s` contrib users:\n", repo.HtmlUrl)

			for _, u := range users {
				fmt.Printf("\t• %s\n", u.GetDescribe())
			}
		}
	}

	return nil
}
