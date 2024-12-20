# Панель Администратора Кафе

Добро пожаловать в проект Панели Администратора Кафе! Это приложение предоставляет удобный интерфейс для управления операциями кафе, включая аутентификацию пользователей, управление меню, обработку заказов и аналитику. Проект состоит из backend на Go и frontend на React, интегрированных с базой данных PostgreSQL.

## Содержание

- Особенности
- Предварительные требования
- Инструкции по настройке
  - Клонирование репозитория
  - Настройка backend
    - Установка зависимостей
    - Конфигурация базы данных
    - Запуск миграций и заполнение базы
    - Запуск backend-сервера
  - Настройка frontend
    - Установка зависимостей
    - Запуск frontend-сервера
- Запуск приложения
- Структура проекта
- Вклад в проект

## Особенности

- **Аутентификация пользователей:** Регистрация и вход с использованием JWT.
- **Управление меню:** Добавление, обновление и удаление позиций меню.
- **Обработка заказов:** Создание и управление заказами с обновлением статуса в реальном времени.
- **Аналитика:** Просмотр статистики по доходам и заказам за разные периоды.
- **Адаптивный интерфейс:** Удобный и адаптивный дизайн, созданный с помощью React.

## Предварительные требования

Перед началом убедитесь, что на вашем компьютере установлены следующие инструменты:

- **Git:** Для контроля версий.
- **Go:** Версия 1.18 или выше.
- **Node.js и npm:** Рекомендуется последняя LTS-версия.
- **PostgreSQL:** Версия 12 или выше.

## Инструкции по настройке

### Клонирование репозитория

```bash
git clone https://github.com/tneki22/cafe-admin-panel.git
cd cafe-admin-panel
```

### Настройка backend

#### Установка зависимостей

Перейдите в директорию backend и установите необходимые зависимости для Go.

```bash
cd backend
go mod download
```

#### Конфигурация базы данных

1. **Установите PostgreSQL:**
   - Если PostgreSQL не установлен, скачайте и установите его с [официального сайта](https://www.postgresql.org/download/).

2. **Создайте базу данных и пользователя:**

   Откройте PostgreSQL Shell и выполните следующие команды:

   ```sql
   CREATE DATABASE cafe-admin-users;
   CREATE USER postgres WITH PASSWORD 'postgres';
   GRANT ALL PRIVILEGES ON DATABASE cafe_admin_users TO postgres;
   ```

3. **Запустите миграции:**

   Создайте таблицы базы данных с помощью следующих SQL-команд:

   ```sql
   -- Таблица пользователей
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(32) NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password VARCHAR(255) NOT NULL
   );

   -- Таблица меню
   CREATE TABLE menu (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       description TEXT NOT NULL,
       price NUMERIC(10, 2) NOT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   );

   -- Таблица заказов
   CREATE TABLE orders (
       id SERIAL PRIMARY KEY,
       total NUMERIC(10, 2) NOT NULL,
       status VARCHAR(50) NOT NULL,
       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
   );

   -- Таблица позиций заказов
   CREATE TABLE order_items (
       id SERIAL PRIMARY KEY,
       order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
       menu_item_id INTEGER REFERENCES menu(id) ON DELETE CASCADE,
       quantity INTEGER NOT NULL
   );
   ```

#### Запуск миграций и заполнение базы

1. **Заполните таблицу меню:**

   В проекте есть скрипт для заполнения таблицы `menu` начальными данными.

   ```bash
   cd ../bd_data_menu
   go run main.go
   ```

2. **Заполните таблицу заказов:**

   Выполните скрипт для добавления случайных данных в таблицу заказов.

   ```bash
   cd ../bd_data_ordersrand
   go run main.go
   ```

#### Запуск backend-сервера

Вернитесь в директорию backend и запустите сервер.

```bash
cd ../backend
go run .
```

Сервер будет работать по адресу 

http://127.0.0.1:8885

.

### Настройка frontend

#### Установка зависимостей

Перейдите в директорию frontend и установите необходимые зависимости для Node.js.

```bash
cd ../src
npm install
```

#### Запуск frontend-сервера

Запустите React сервер для разработки.

```bash
npm start
```

Frontend будет доступен по адресу 

http://localhost:3000

.

## Запуск приложения

1. **Убедитесь, что PostgreSQL работает:** Проверьте, что сервер PostgreSQL запущен, а база `cafe_admin_users` доступна.

2. **Запустите backend-сервер:** Следуйте описанным выше шагам.

3. **Запустите frontend-сервер:** Следуйте описанным выше шагам.

4. **Доступ к приложению:** Откройте браузер и перейдите на 

http://localhost:3000

.

## Структура проекта

```
cafe-admin-panel/
├── backend/
│   ├── main.go
│   ├── config.go
│   ├── database.go
│   ├── usecase.go
│   ├── api.go
│   ├── auth.yaml
│   └── pkg/
│       └── vars/
│           └── jwt.go
├── bd_data_menu/
│   └── main.go
├── bd_data_ordersrand/
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── src/
│   ├── components/
│   │   ├── LoginPage.js
│   │   ├── MainPage.js
│   │   ├── OrdersPage.js
│   │   ├── MenuPage.js
│   │   └── AnalyticsPage.js
│   ├── styles/
│   │   ├── App.css
│   │   ├── LoginPage.css
│   │   ├── MainPage.css
│   │   ├── OrdersPage.css
│   │   ├── MenuPage.css
│   │   └── AnalyticsPage.css
│   ├── App.js
│   ├── App.test.js
│   ├── index.js
│   ├── index.css
│   ├── reportWebVitals.js
│   └── setupTests.js
├── public/
│   ├── index.html
│   ├── manifest.json
│   └── robots.txt
├── package.json
└── README.md
```

## Вклад в проект

Мы приветствуем ваш вклад! Следуйте этим шагам, чтобы внести изменения:

1. **Форк репозитория:** Нажмите кнопку **Fork** на странице репозитория.

2. **Клонируйте свой форк:**

   ```bash
   git clone https://github.com/tneki22/cafe-admin-panel.git
   cd cafe-admin-panel
   ```

3. **Создайте новую ветку:**

   ```bash
   git checkout -b feature/YourFeatureName
   ```

4. **Внесите изменения:** Реализуйте новую функцию или исправление ошибки.

5. **Зафиксируйте изменения:**

   ```bash
   git commit -m "Добавьте описание изменений"
   ```

6. **Отправьте изменения:**

   ```bash
   git push origin feature/YourFeatureName
   ```

7. **Создайте Pull Request:** Перейдите в оригинальный репозиторий и создайте Pull Request из вашего форка.
