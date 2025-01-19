[![codecov](https://codecov.io/gh/kamencov/goph-keeper/branch/iter3/graph/badge.svg?token=1ZHJfDLN5k)](https://codecov.io/gh/kamencov/goph-keeper)
# goph-keeper

Шаги:

1. Сначала требуется запустить сервер

        go run cmd/service/main.go

2. Запускаем клиента

         go run cmd/client/main.go
___
В сервисе клиенте реализованы метода:

- **Register** - регистрирует клиента через сервис gRPC.
- **Auth** - авторизация клиента через сервис gRPC, получает token и хранит у клиента.
- **Quite** - выходит из клиента.
___
После авторизации открывается возможность сохранять, искать, удалять:
___
1. **Find all data** - выводит все данные.
___
2. **Creadentials** - хранит данные в виде resource, login, password:

**Save** - реализован.

**Deleted** - нет реализации.

**Quite** - выходит из клиента.
___
3. **Text** - хранит данные в виде text:

**Save** - нет реализации.

**Deleted** - нет реализации.

**Quite** - выходит из клиента.
___
4. **Binary** - хранит данные в виде binary:

**Save** - нет реализации.

**Deleted** - нет реализации.

**Quite** - выходит из клиента.
___
5. **Cards** - хранит данные в виде card number:

**Save** - нет реализации.

**Deleted** - нет реализации.

**Quite** - выходит из клиента.