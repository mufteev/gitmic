package testinit

import (
	"os"
	"path"
	"runtime"
)

// Специальный Хак, который позволяет поднять уровень текущей директории выше
// Чтобы не изменять env-функцию получения данных из файла
// https://intellij-support.jetbrains.com/hc/en-us/community/posts/360009685279-Go-test-working-directory-keeps-changing-to-dir-of-the-test-file-instead-of-value-in-template?page=1#community_comment_360002183640
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")

	err := os.Chdir(dir)

	if err != nil {
		panic(err)
	}
}
