 
## Сервис аутентификации

Интаграционные тесты  в cmd/client
проверяется взаимодействие по grpc 
создание пользователя и добавления его в бд

scp -P443 -r ./* root@94.159.99.214:/root/src/go/temp/lms/02/netevent
scp -P443 -r ./.* root@94.159.99.214:/root/src/go/temp/lms/02/netevent 
tar -xvf 



сервис авторизаций порт 5100
сервис уведомлений      5200
сервис событий          5300



## sql create manual

```sql
SELECT * FROM pg_catalog.pg_tables where pg_catalog.pg_tables.schemaname='public';
DROP TABLE IF EXISTS "tuser" CASCADE;
CREATE TABLE IF NOT EXISTS "tuser" (
    id INT PRIMARY KEY NOT NULL,
    name VARCHAR(30),
    password VARCHAR(30), 
    email   VARCHAR(30),
    interest TEXT,
    accesstkn TEXT,
    accessttl INT,
    refreshtkn TEXT,
    refreshttl INT
);

SELECT * FROM tuser;


```

```sql
SELECT * FROM pg_catalog.pg_tables where pg_catalog.pg_tables.schemaname='public';
DROP TABLE IF EXISTS "tuser" CASCADE;
DROP TABLE IF EXISTS "tevent" CASCADE;
DROP TABLE IF EXISTS "tparticipant";
CREATE TABLE IF NOT EXISTS "tuser" (
    id INT PRIMARY KEY NOT NULL,
    name VARCHAR(30) NOT NULL,
    password VARCHAR(30) NOT NULL, 
    email   VARCHAR(30),
    role VARCHAR(30) NOT NULL,
    interest TEXT,
    accesstkn TEXT,
    accessttl INT,
    refreshtkn TEXT,
    refreshttl INT
);
/*
CREATE TABLE IF NOT EXISTS  "tevent" (
    id INT PRIMARY KEY NOT NULL,
    createid INT,
    title VARCHAR(255),
    description TEXT,
    time TIME,
    place VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS "tparticipant" (
    userID INT REFERENCES tuser(id) ON DELETE CASCADE,
    eventID INT REFERENCES tevent(id) ON DELETE CASCADE
);
*/
SELECT * FROM tuser;
```