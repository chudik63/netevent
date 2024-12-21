# Netevent
App for networking on events 
# Installation
```
git clone git@gitlab.crja72.ru:gospec/go9/netevent.git
cd netevent
```
# Main Commands
Build and start all containers:
```
docker-compose up
```

# Usage
## API Documentation

### 1. SignUp: POST /api/v1/sign-up
Create a new user account.

#### Request Body:
```json
{
  "id": "USER_ID",
  "name": "USER_NAME",
  "email": "USER_EMAIL",
  "password": "USER_PASSWORD",
  "role": "USER_ROLE",
  "interests": ["INTEREST_1", "INTEREST_2"]
}
```

#### Response:
```json
{
  "message": "success or error message"
}
```

#### Example cURL Request:
```bash
curl -X POST http://localhost:80/api/v1/sign-up \
     -d '{"id": 1, "name": "John Doe", "email": "john.doe@example.com", "password": "password123", "role": "user", "interests": ["sports", "music"]}' \
     -H "Content-Type: application/json"
```

---

### 2. SignIn: POST /api/v1/sign-in
User sign-in to get authentication tokens.

#### Request Body:
```json
{
  "name": "USER_NAME",
  "password": "USER_PASSWORD"
}
```

#### Response:
```json
{
  "access_token": "ACCESS_TOKEN",
  "access_token_ttl": "ACCESS_TOKEN_TTL",
  "refresh_token": "REFRESH_TOKEN",
  "refresh_token_ttl": "REFRESH_TOKEN_TTL"
}
```

#### Example cURL Request:
```bash
curl -X POST http://localhost:80/api/v1/sign-in \
     -d '{"name": "John Doe", "password": "password123"}' \
     -H "Content-Type: application/json"
```

---

### 4. Create Event: POST /api/v1/event
Create a new event.

#### Request Body:
```json
{
  "event": {
    "event_id": "EVENT_ID",
    "creator_id": "CREATOR_ID",
    "title": "EVENT_TITLE",
    "description": "EVENT_DESCRIPTION",
    "time": "EVENT_TIME",
    "place": "EVENT_PLACE",
    "interests": ["INTEREST_1", "INTEREST_2"]
  }
}
```

#### Response:
```json
{
  "event_id": "EVENT_ID"
}
```

#### Example cURL Request:
```bash
curl -X POST http://localhost:80/api/v1/event \
     -d '{"event": {"creator_id": 1, "title": "My Event", "description": "Event Description", "time": "2024-12-21T10:00:00", "place": "Place", "interests": ["coding", "tech"]}}' \
     -H "Content-Type: application/json"
```

---

### 5. Read Event: GET /api/v1/event/{event_id}
Read event details by event ID.

#### Response:
```json
{
  "event": {
    "event_id": "EVENT_ID",
    "creator_id": "CREATOR_ID",
    "title": "EVENT_TITLE",
    "description": "EVENT_DESCRIPTION",
    "time": "EVENT_TIME",
    "place": "EVENT_PLACE",
    "interests": ["INTEREST_1", "INTEREST_2"]
  }
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event/1
```

---

### 6. Update Event: PUT /api/v1/event/{event_id}
Update event details by event ID.

#### Request Body:
```json
{
  "event": {
    "event_id": "EVENT_ID",
    "creator_id": "CREATOR_ID",
    "title": "EVENT_TITLE",
    "description": "EVENT_DESCRIPTION",
    "time": "EVENT_TIME",
    "place": "EVENT_PLACE",
    "interests": ["INTEREST_1", "INTEREST_2"]
  }
}
```

#### Response:
```json
{}
```

#### Example cURL Request:
```bash
curl -X PUT http://localhost:80/api/v1/event/1 \
     -d '{"event": {"title": "Updated Event", "description": "Updated Description", "time": "2024-12-22T12:00:00", "place": "New Place", "interests": ["new_interest"]}}' \
     -H "Content-Type: application/json"
```

---

### 7. Delete Event: DELETE /api/v1/event/{event_id}
Delete an event by event ID.

#### Response:
```json
{}
```

#### Example cURL Request:
```bash
curl -X DELETE http://localhost:80/api/v1/event/1
```

---

### 8. List Events: GET /api/v1/event
List all events.

#### Response:
```json
{
  "events": [
    {
      "event_id": "EVENT_ID",
      "creator_id": "CREATOR_ID",
      "title": "EVENT_TITLE",
      "description": "EVENT_DESCRIPTION",
      "time": "EVENT_TIME",
      "place": "EVENT_PLACE",
      "interests": ["INTEREST_1", "INTEREST_2"]
    }
  ]
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event
```

---

### 9. List Events By Creator: GET /api/v1/event/creator/{creator_id}
List all events created by a specific creator.

#### Response:
```json
{
  "events": [
    {
      "event_id": "EVENT_ID",
      "creator_id": "CREATOR_ID",
      "title": "EVENT_TITLE",
      "description": "EVENT_DESCRIPTION",
      "time": "EVENT_TIME",
      "place": "EVENT_PLACE",
      "interests": ["INTEREST_1", "INTEREST_2"]
    }
  ]
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event/creator/1
```

---

### 10. List Events By Interests: GET /api/v1/event/interests/{user_id}
List all events based on user's interests.

#### Response:
```json
{
  "events": [
    {
      "event_id": "EVENT_ID",
      "creator_id": "CREATOR_ID",
      "title": "EVENT_TITLE",
      "description": "EVENT_DESCRIPTION",
      "time": "EVENT_TIME",
      "place": "EVENT_PLACE",
      "interests": ["INTEREST_1", "INTEREST_2"]
    }
  ]
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event/interests/1
```

---

### 11. Register User to Event: POST /api/v1/event/register
Register a user for an event.

#### Request Body:
```json
{
  "user_id": "USER_ID",
  "event_id": "EVENT_ID"
}
```

#### Response:
```json
{}
```

#### Example cURL Request:
```bash
curl -X POST http://localhost:80/api/v1/event/register \
     -d '{"user_id": 1, "event_id": 1}' \
     -H "Content-Type: application/json"
```

---

### 12. Set Chat Status: PUT /api/v1/event/chat
Set the chat status for a user in an event.

#### Request Body:
```json
{
  "user_id": "USER_ID",
  "event_id": "EVENT_ID",
  "is_ready": "true_or_false"
}
```

#### Response:
```json
{}
```

#### Example cURL Request:
```bash
curl -X PUT http://localhost:80/api/v1/event/chat \
     -d '{"user_id": 1, "event_id": 1, "is_ready": true}' \
     -H "Content-Type: application/json"
```

---

### 13. List Users to Chat: GET /api/v1/event/chat/{event_id}/{user_id}
List users available to chat in an event.

#### Response:
```json
{
  "participants": [
    {
      "user_id": "USER_ID",
      "name": "USER_NAME",
      "interests": ["INTEREST_1", "INTEREST_2"]
    }
  ]
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event/chat/1/1
```

---

### 14. List Registered Events: GET /api/v1/event/registrated/{user_id}
List all events a user has registered for.

#### Response:
```json
{
  "events": [
    {
      "event_id": "EVENT_ID",
      "creator_id": "CREATOR_ID",
      "title": "EVENT_TITLE",
      "description": "EVENT_DESCRIPTION",
      "time": "EVENT_TIME",
      "place": "EVENT_PLACE",
      "interests": ["INTEREST_1", "INTEREST_2"]
    }
  ]
}
```

#### Example cURL Request:
```bash
curl -X GET http://localhost:80/api/v1/event/registrated/1
```

---


# Technologies Used
- Golang
- GRPC
- PostgreSQL
- Reddis
- Kafka
- Docker
