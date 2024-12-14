 
## Сервис аутентификации

scp -P443 -r ./* root@94.159.99.214:/root/src/go/temp/lms/02/netevent
scp -P443 -r ./.* root@94.159.99.214:/root/src/go/temp/lms/02/netevent 


сервис авторизаций порт 5100
сервис уведомлений      5200
сервис событий          5300

как обновить приложение 
make build
послать в контейнер экзе файл
make update


## sql create manual

```sql
CREATE TABLE IF NOT EXISTS "User" (
    userId INT PRIMARY KEY NOT NULL,
    name VARCHAR(30),
    secondName VARCHAR(30),
    passwdHash VARCHAR(30), //что бы не хранить пароль в открытом виде
    email   VARCHAR(30),
    role    int,
    interest TEXT 
);

CREATE TABLE IF NOT EXISTS  Interest (
    interestId INT NOT NULL,
    interest   VARCHAR(30)
    FOREIGN KEY InterestId REFERENCES User(InterestId) ON DELETE CASCADE,
);

CREATE TABLE IF NOT EXISTS  "Event" (
    eventId INT  PRIMARY KEY NOT NULL,
    creatorId INT,
    title VARCHAR(255),
    description TEXT,
    time TIME,
    place VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS "Participant" (
    userId INT NOT NULL,
    eventID INT NOT NULL,
    FOREIGN KEY userId REFERENCES User(UserId) ON DELETE,
    FOREIGN KEY eventID REFERENCES Event(eventId) ON DELETE CASCADE
);


```

```sql
SELECT * FROM pg_catalog.pg_tables where pg_catalog.pg_tables.schemaname='public';
DROP TABLE IF EXISTS "tuser" CASCADE;
DROP TABLE IF EXISTS "tevent" CASCADE;
DROP TABLE IF EXISTS "tparticipant";
CREATE TABLE IF NOT EXISTS "tuser" (
    userId INT PRIMARY KEY NOT NULL,
    name VARCHAR(30),
    secondName VARCHAR(30),
    passwdHash VARCHAR(30), 
    email   VARCHAR(30),
    role    int,
    interest TEXT 
);

CREATE TABLE IF NOT EXISTS  "tevent" (
    eventId INT PRIMARY KEY NOT NULL,
    creatorId INT,
    title VARCHAR(255),
    description TEXT,
    time TIME,
    place VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS "tparticipant" (
    userID INT REFERENCES tuser(userId) ON DELETE CASCADE,
    eventID INT REFERENCES tevent(eventId) ON DELETE CASCADE
);

```