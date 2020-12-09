# HTTP API
The server runs on port 4000

Responses should always be in json format
They should always contain a field that tells if the operation was succesful.

### Example of an error message:
```json
{
    "Success": false,
    "Message": "Sample error message"
}
```
# Routes:
### Get all subjects:
GET `/subjects`
### Get specific subject:
GET `/subject/:id`

### Get a specific course:
GET `/course/:id`
### Get all courses:
GET `/courses`

### Get groups in a specific season:
GET `/groups/:seasonID`

### Get an users schedule:
POST `/schedule`

### Websocket connection for realtime updates from server:
Websocket connection to: `/ws`

### Login and authentication:
POST `/auth/login`

POST `/auth/adduser`