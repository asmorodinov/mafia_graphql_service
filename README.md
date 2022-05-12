# mafia_graphql_service

## GraphQL
Непосредственно GraphQL сервер и клиент находятся в [папке](golang/GraphQL).

Собрать их можно либо используя стандартные комманды golang (например ```cd golang && go mod download && cd GraphQL/server && go run .```), 
либо, сервер можно собрать с помощью докера: ```cd golang && docker build . -t graphqlserver -f=DockerfileGraphQLServer``` (для клиента я не создавал докерфайл, но его можно написать по аналогии с сервером).

Докер контейнер с GraphQL сервером собирать самому не обязательно, можно воспользоваться готовым образом ```asmorodinov/graphqlserver```.

Например, так можно быстро протестировать клиент и сервер: 

```docker run -i -p 8090:8090 asmorodinov/graphqlserver -addr="[::]:8090"```

```cd golang && go mod download && cd GraphQL/client && go run . --help```, или запуск любой другой команды, например ```go run . -getGamesList```

Клиент довольно примитивный, он просто принимает в виде аргументов командной строки информацию о действии, которое мы хотим выполнить (добавление комментария, просмотр списка игр и т.д.), отправляет запрос к серверу, печатает результат и завершает работу.

## Остальная часть
Все остальные файлы с небольшими изменениями были скопированы из [прошлого задания](https://github.com/asmorodinov/mafia_rest_service).

### Изменения относительно репозитория mafia_rest_service
- Python часть приложения и golang часть были перемещены в отдельные папки.
- mafia_server.py при запуске теперь просит ввести адрес GraphQL сервиса (например http://localhost:8090, либо http://{имя контейнера}:8090, по аналогии с адресом REST сервиса).
- Когда в игре происходит какое-то событие, mafia_server.py отправляет POST запрос к GraphQL серверу с информацией о текущем состоянии игры.
- В папке golang (которой раньше не было, но тем не менее) появилась папка GraphQL и файл DockerfileGraphQLServer (довольно очевидное изменение, но всё равно решил его упомянуть).

## Запуск всех компонент (GraphQL, mafia_server, RestServer, PDFWorker, server_tcp)
Нужно смотреть readme прошлой задачи, здесь я упомяну лишь изменения относительно него.

- Если запускаем без докера всё кроме rabbitmq, то нужно дополнительно запустить GraphQL сервер ```cd golang && go mod download && cd GraphQL/server && go run .```, а также при запуске mafia_server.py указать адрес GraphQL сервера (http://localhost:8090).
- Если запускаем серверную часть в докере, то нужно запустить GraphQL сервер командой ```docker run -i --network test-net --name graphqlserver -p 8090:8090 asmorodinov/graphqlserver -addr=[::]:8090```, а также при запуске mafia_server.py указать адрес GraphQL сервера (http://graphqlserver:8090).
- В местах где делаем ```cd``` может быть необходимо дополнительно указать python, либо golang (например ```cd basic_mafia``` -> ```cd python/basic_mafia```).
- Так как mafia_server.py обновился, то нужно использовать asmorodinov/mafiaserverwithgraphql вместо asmorodinov/mafiaserver, если mafia_server.py запускался через докер.
