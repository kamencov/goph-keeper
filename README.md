[![codecov](https://codecov.io/gh/kamencov/goph-keeper/branch/iter3/graph/badge.svg?token=1ZHJfDLN5k)](https://codecov.io/gh/kamencov/goph-keeper)
# goph-keeper

Шаги:

1. Сначала требуется запустить сервер

        go run cmd/service/main.go

2. Запускаем клиента

         go run cmd/client/main.go
___
В сервисе клиенте реализованы метода online:

- **Register** - регистрирует клиента через сервис gRPC.
- **Auth** - авторизация клиента через сервис gRPC, получает token и хранит у клиента.
- **Quite** - выходит из клиента.

Offline режим:

![img_5.png](img_5.png)
- **Auth** - авторизация клиента, который хранится у клиента.

![img_6.png](img_6.png)
- **Quite** - выходит из клиента.
___
После авторизации открывается возможность сохранять, искать, удалять:

![img_7.png](img_7.png)
___
1. **Find all data** - выводит все данные в виде таблицы.

![img_8.png](img_8.png)
___
2. **Creadentials** - хранит данные в виде resource, login, password:

![img_9.png](img_9.png)
- **Save** - реализован.
- **Deleted** - реализован.
- **Quite** - реализован.
___
3. **Text** - хранит данные в виде text:

![img_10.png](img_10.png)
- **Save** - реализован.
- **Deleted** - реализован.
- **Quite** - реализован.
___
4. **Binary** - хранит данные в виде binary:

![img_11.png](img_11.png)
- **Save** - реализован.
- **Deleted** - реализован.
- **Quite** - реализован.
___
5. **Cards** - хранит данные в виде card number:

![img_12.png](img_12.png)
- **Save** - реализован.
- **Deleted** - реализован.
- **Quite** - реализован.
___
Пример работы sync клиента:
![img.png](img.png)