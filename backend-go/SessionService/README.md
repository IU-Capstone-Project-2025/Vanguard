Session service

При создании сессии:
1) Создает сессию с полями, WsConnection остается пустым, отправляет это в редис
{
   ID    string
   Code  string 
   State string
   ServerWsConnection string
}
2) Запрашивает данные по квизу от другого сервиса, отправляет в брокер сообщение и данные по квизу

При Добавлении игроков
1) принимает от клиента код и проверяет его
2) берет с редиса WsConnection, отправляет новуму юзеру токен с такими полями
{
   UserId             string 
   UserType           string 
   ServerWsConnection string 
   CurrentQuiz        string 
   Exp                int64  
}
При загрузке вопросов 
1) отправляет в раббит сообщение о новом вопросе
2) при завершении квиза удаляет его из редиса и отправляет в брокер сообщение о завершении


Была реализованна первая версия Session Service.
Реализованно создание сессии:
1) Генерация уникального ID сессии (UUID).
2) Генерация 6-значного буквенно-цифрового кода для подключения игроков.
3) Генерация токена для администратора сессии (включает nickname, session_code, exp).
4) Сохранение информации о сессии в Redis (ключ session:{code}).

Проверка кода нового участника:
1) Обработка запроса от клиента с кодом квиза.
2) Проверка существования сессии в Redis.
3) Генерация JWT токена для игрока (включает session_code, nickname, exp).

Запуск сервера:
Сервер запускается на localhost:8000. На данный момент доступно 3 гет запроса:
1) /validate проверяет код нового участника
2) /create создает новую сессию
3) /Next отправляет запрос на следующий вопрос

следующие шаги:
1) реализовать корректную отсновку сервера
2) написание юнит тестов


1. Добавлены юнит-тесты
Покрытие ключевых компонентов логики создания и валидации сессий.
Проверка генерации токенов, кодов подключения и взаимодействия с Redis.

2. Интеграция Swagger
Использована библиотека swaggo/swag для генерации документации.

Swagger UI доступен по эндпоинту /swagger/index.html.
Все доступные маршруты задокументированы с примерами запросов и ответов.

3. Изменены эндпоинты
/create → POST /sessions

/validate → POST /join

/next → POST /session/{id}/nextQuestion

4. бновлены схемы запросов и ответов
Добавлены модели CreateSessionRequest, ValidateSessionRequest, SessionResponse и др.
Ответы теперь возвращаются в едином формате с status, message.

5.  Обновлён подход к генерации токенов
Добавлены зарегистрированные поля в JWT (exp, iat, iss).
