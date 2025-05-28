package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"diprec_api/internal/config"
	"diprec_api/internal/domain"
	"diprec_api/internal/infrastructure/db/postgres"
)

/*────────────────────────── JSON helper ─────────────────────────*/

func j(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

/*────────────────────────── TRUNCATE helper ─────────────────────*/

func truncateAll(db *gorm.DB) error {
	return db.Exec(`TRUNCATE courses, questions, tests, test_questions 
	                RESTART IDENTITY CASCADE`).Error
}

/*────────────────────────── main ────────────────────────────────*/

func main() {
	cfg := config.MustLoad()
	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		User:            cfg.DB.User,
		Password:        cfg.DB.Password,
		DBName:          cfg.DB.DBName,
		SSLMode:         cfg.DB.SSLMode,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}

	// Миграция
	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Очистка
	if err := truncateAll(db); err != nil {
		log.Fatalf("Truncate failed: %v", err)
	}

	// Сидирование
	if err := seed_only_questions(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	fmt.Println("✓ Mock-данные успешно загружены!")
}

/*────────────────────────── Seeding ─────────────────────────────*/

func seed_only_questions(db *gorm.DB) error {
	ctx := context.Background()

	/* 2) 60 вопросов */
	qs := questions() // функция ниже
	if err := db.WithContext(ctx).Create(&qs).Error; err != nil {
		return err
	}

	return nil
}

/*────────────────────────── 60 вопросов ─────────────────────────*/

func questions() []domain.Question {
	return []domain.Question{
		/* ---------- Frontend: HTML & CSS ---------- */
		{Title: "Какой HTML-тег создаёт абзац?", Type: domain.Single,
			Variants: j(map[string]string{
				"a": "<p>",
				"b": "<div>",
				"c": "<br>",
				"d": "<span>",
			}), Answer: j("a")},

		{Title: "Сколько уровней заголовков (h1–h?) есть в HTML-стандарте?", Type: domain.Number,
			Variants: j(nil), Answer: j(6)},

		{Title: "CSS-свойство, задающее цвет текста", Type: domain.Text,
			Variants: j(nil), Answer: j("color")},

		{Title: "Что из списка делает элемент flex-контейнером?", Type: domain.Multiple,
			Variants: j(map[string]string{
				"a": "display: flex",
				"b": "display: grid",
				"c": "float: left",
				"d": "position: absolute",
			}), Answer: j([]string{"a"})},

		{Title: "Какой атрибут задаёт альтернативный текст картинки?", Type: domain.Single,
			Variants: j(map[string]string{
				"a": "src",
				"b": "alt",
				"c": "title",
				"d": "href",
			}), Answer: j("b")},

		{Title: "Расширение файла таблиц стилей", Type: domain.Single,
			Variants: j(map[string]string{
				"a": ".html",
				"b": ".css",
				"c": ".js",
				"d": ".png",
			}), Answer: j("b")},

		{Title: "Что означает аббревиатура «HTML»?", Type: domain.Single,
			Variants: j(map[string]string{
				"a": "HyperText Markup Language",
				"b": "Hyperlinks and Text Markup Language",
				"c": "Home Tool Markup Language",
				"d": "HighText Machine Language",
			}), Answer: j("a")},

		{Title: "Единица rem основывается на …", Type: domain.Single,
			Variants: j(map[string]string{
				"a": "текущий font-size элемента",
				"b": "font-size корневого (html) элемента",
				"c": "ширину viewport",
				"d": "абсолютные пиксели (px)",
			}), Answer: j("b")},

		{Title: "Какой meta-тег включает адаптивную вёрстку на мобильных устройствах?", Type: domain.Single,
			Variants: j(map[string]string{
				"a": `<meta http-equiv="Content-Type" content="text/html; charset=utf-8">`,
				"b": `<meta name="viewport" content="width=device-width, initial-scale=1">`,
				"c": `<meta name="theme-color" content="#ffffff">`,
				"d": `<meta name="format-detection" content="telephone=no">`,
			}), Answer: j("b")},

		{Title: "CSS-свойство, создающее внешнюю тень блока", Type: domain.Single,
			Variants: j(map[string]string{
				"a": "box-shadow",
				"b": "text-shadow",
				"c": "filter: drop-shadow",
				"d": "outline",
			}), Answer: j("a")},

		/* ---------- Frontend: JavaScript ---------- */

		{Title: "Как объявить переменную с блочной областью видимости?", Type: domain.Single,
			Variants: j(map[string]string{"a": "var", "b": "let", "c": "const", "d": "function"}), Answer: j("b")},
		{Title: "Какой результат у выражения `typeof null`?", Type: domain.Text,
			Variants: j(nil), Answer: j("object")},
		{Title: "Методы массива для добавления/удаления с конца", Type: domain.Multiple,
			Variants: j(map[string]string{"a": "push", "b": "pop", "c": "shift", "d": "unshift"}),
			Answer:   j([]string{"a", "b"})},
		{Title: "Что вернёт `Boolean([])`?", Type: domain.Single,
			Variants: j(map[string]string{"a": "true", "b": "false"}), Answer: j("a")},
		{Title: "Чему равняется `NaN === NaN`?", Type: domain.Text,
			Variants: j(nil), Answer: j("false")},
		{Title: "Метод для преобразования JSON-строки в объект", Type: domain.Text,
			Variants: j(nil), Answer: j("JSON.parse")},
		{Title: "Какая функция ставит колбэк в микротаски?", Type: domain.Text,
			Variants: j(nil), Answer: j("Promise.resolve().then")},
		{Title: "Сколько значений имеет булев тип?", Type: domain.Number,
			Variants: j(nil), Answer: j(2)},
		{Title: "Как называется корень DOM-дерева?", Type: domain.Text,
			Variants: j(nil), Answer: j("document")},
		{Title: "`0 == '0'` возвращает…", Type: domain.Text,
			Variants: j(nil), Answer: j("true")},

		/* ---------- Frontend: React ---------- */

		{Title: "Какой хук используется для локального состояния?", Type: domain.Text,
			Variants: j(nil), Answer: j("useState")},
		{Title: "Ключевой prop для списков в React", Type: domain.Text,
			Variants: j(nil), Answer: j("key")},
		{Title: "Как называется однонаправленный поток данных в React?", Type: domain.Text,
			Variants: j(nil), Answer: j("props")},
		{Title: "Что возвращает компонент React?", Type: domain.Single,
			Variants: j(map[string]string{"a": "string", "b": "HTML-текст", "c": "JSX", "d": "CSS-объект"}), Answer: j("c")},
		{Title: "Какой хук выполняется после каждого рендера?", Type: domain.Single,
			Variants: j(map[string]string{"a": "useEffect", "b": "useMemo", "c": "useRef", "d": "useCallback"}), Answer: j("a")},
		{Title: "Какой метод удалён в React 18?", Type: domain.Text,
			Variants: j(nil), Answer: j("ReactDOM.render")},
		{Title: "Компоненты с побочным эффектом — …", Type: domain.Single,
			Variants: j(map[string]string{"a": "pure components", "b": "class components", "c": "memo components", "d": "functional components"}),
			Answer:   j("b")},
		{Title: "Как получить children компонента?", Type: domain.Text,
			Variants: j(nil), Answer: j("props.children")},
		{Title: "Где хранится глобальное состояние в Context API?", Type: domain.Text,
			Variants: j(nil), Answer: j("Provider")},
		{Title: "Сколько корневых элементов допускает JSX?", Type: domain.Number,
			Variants: j(nil), Answer: j(1)},

		/* ---------- Backend: HTTP & REST ---------- */

		{Title: "Какой метод HTTP безопасный и идемпотентный?", Type: domain.Single,
			Variants: j(map[string]string{"a": "POST", "b": "PUT", "c": "GET", "d": "PATCH"}), Answer: j("c")},
		{Title: "Код 201 означает…", Type: domain.Text,
			Variants: j(nil), Answer: j("created")},
		{Title: "Выберите заголовки, относящиеся к кешированию", Type: domain.Multiple,
			Variants: j(map[string]string{"a": "Cache-Control", "b": "ETag", "c": "Accept", "d": "Content-Type"}),
			Answer:   j([]string{"a", "b"})},
		{Title: "Какой статус вернёт успешный DELETE?", Type: domain.Number,
			Variants: j(nil), Answer: j(204)},
		{Title: "Какой порт по умолчанию у HTTPS?", Type: domain.Number,
			Variants: j(nil), Answer: j(443)},
		{Title: "`Idempotent` значит…", Type: domain.Text,
			Variants: j(nil), Answer: j("повторяемый без побочных эффектов")},
		{Title: "Метод для частичного обновления ресурса", Type: domain.Text,
			Variants: j(nil), Answer: j("PATCH")},
		{Title: "Какой протокол лежит в основе HTTP/2?", Type: domain.Text,
			Variants: j(nil), Answer: j("TCP")},
		{Title: "Формат тела в REST чаще всего…", Type: domain.Text,
			Variants: j(nil), Answer: j("JSON")},
		{Title: "Заголовок для авторизации по JWT", Type: domain.Text,
			Variants: j(nil), Answer: j("Authorization")},

		/* ---------- Backend: Go Basics ---------- */
		{Title: "Ключевое слово запуска горутину", Type: domain.Text,
			Variants: j(nil), Answer: j("go")},
		{Title: "Какой интерфейс имеет метод `Error() string`?", Type: domain.Text,
			Variants: j(nil), Answer: j("error")},
		{Title: "Zero value для int", Type: domain.Number,
			Variants: j(nil), Answer: j(0)},
		{Title: "Пакет для конкурентной работы с каналами", Type: domain.Text,
			Variants: j(nil), Answer: j("sync")},
		{Title: "Сколько байт занимает int32?", Type: domain.Number,
			Variants: j(nil), Answer: j(4)},
		{Title: "Импорт формата времени RFC3339 находится в пакете…", Type: domain.Text,
			Variants: j(nil), Answer: j("time")},
		{Title: "`panic` останавливает…", Type: domain.Text,
			Variants: j(nil), Answer: j("goroutine")},
		{Title: "Как создать слайс длины 0 и cap 5?", Type: domain.Text,
			Variants: j(nil), Answer: j("make([]T,0,5)")},
		{Title: "Команда go для обновления зависимостей", Type: domain.Text,
			Variants: j(nil), Answer: j("go get -u")},
		{Title: "Функция для логирования фатальной ошибки", Type: domain.Text,
			Variants: j(nil), Answer: j("log.Fatal")},

		/* ---------- Backend: SQL ---------- */
		{Title: "Ключевое слово для сортировки", Type: domain.Text,
			Variants: j(nil), Answer: j("ORDER BY")},
		{Title: "Какой тип индекса ищет по подстроке?", Type: domain.Text,
			Variants: j(nil), Answer: j("GIN / FULLTEXT")},
		{Title: "Выберите DDL-операции", Type: domain.Multiple,
			Variants: j(map[string]string{"a": "CREATE", "b": "DROP", "c": "UPDATE", "d": "ALTER"}),
			Answer:   j([]string{"a", "b", "d"})},
		{Title: "Команда для удаления всей таблицы с данными", Type: domain.Text,
			Variants: j(nil), Answer: j("TRUNCATE")},
		{Title: "Сколько записей вернёт `LIMIT 1`?", Type: domain.Number,
			Variants: j(nil), Answer: j(1)},
		{Title: "Как называется соединение, выводящее пересечение двух наборов?", Type: domain.Text,
			Variants: j(nil), Answer: j("INNER JOIN")},
		{Title: "Функция подсчёта строк", Type: domain.Text,
			Variants: j(nil), Answer: j("COUNT()")},
		{Title: "Уровень изоляции по умолчанию в PostgreSQL", Type: domain.Text,
			Variants: j(nil), Answer: j("Read Committed")},
		{Title: "Тип данных для хранения JSON в PostgreSQL", Type: domain.Text,
			Variants: j(nil), Answer: j("jsonb")},
		{Title: "Оператор объединения результатов двух SELECT", Type: domain.Text,
			Variants: j(nil), Answer: j("UNION")},
	}
}
