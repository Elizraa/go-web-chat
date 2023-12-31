# Globes
Reference https://github.com/DavidSchott/chitchat

A live version of the current master branch is deployed to [54.169.251.97](http://54.169.251.97).
But curently is limited by aws free tier package, only the endoint API still working (for now).

## Documentation ##
[postman](https://documenter.getpostman.com/view/15592249/2s9YJaXiUx)

## Install ##
All that you need is [Golang](https://golang.org/). Once you run the application, it will expose a [target port](./config.json) on the host.

```
20:5:13 app         | Globes 0.4 started at 127.0.0.1:443
```

## Configuration ##
Edit [config.json](./config.json) to configure HTTP server settings.
Add [.env](./env.loval) file to configure mongodb data.




## ChatRoom API ##
Send HTTP requests to `/chats`:
  * `GET /chats/<room_id>`: retrieve a chat room by ID or title
  * `POST /chats`: create a new chat room
  * `PUT /chats/<room_id>`: update existing chat room by ID or title
  * `DEL /chats/<room_id>`: delete a chat room

E.g. A basic chat room POST request could look as follows:
```json
{
   "title":"My chat room",
   "description":"There are many like it, but this one is mine",
   "visibility":"public"
}
```
If successful, the server will respond with HTTP code 201 and the newly created chat room resource:
```json
{
    "title": "My chat room",
    "description": "There are many like it, but this one is mine",
    "visibility": "public",
    "createdAt": "2020-12-03T21:23:54.4213184-08:00",
    "updatedAt": "2020-12-03T21:23:54.4213184-08:00",
    "id": 2,
    "users": []
}
```

If unsuccessful, the server will return [an error](./data/apierror.go) in its response body, e.g. if a room with the same title already exists HTTP code 400 will be sent along with the following body content: 
```json
{
    "status": false,
    "error": {
        "code": 102,
        "error": "Room error: Duplicate room",
        "field": "title"
    }
}
```
The required fields to create a public chat room are `"title"` and `"visibility"`. "Private" and "hidden" chat rooms also require a password.

  ## Authorization
  Authorization is implemented using [JSON Web Tokens](https://jwt.io/introduction/) in an attempt to achieve statelessness and scalability. In order to authorize access or modification of a password-protected chat room, users will need to pass in a token in the `Authorization` HTTP header using the Bearer scheme.
  
  The token can be requested by sending a JSON body containing the password to the `/chats/<room_id>/token` endpoint. For example:
  ```json
  {
   "secret":"my_secret",
   "name":"david"
}
  ```
  If successful, the server will respond with the token in the response body.
```json
{
    "status": true,
    "name": "david",
    "room_id": 2,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImRhdmlkIiwicm9vbV9pZCI6MiwiZXhwIjoxNjA3NTg1MjQ5fQ.b6XnNqrFnFmuUMhTBKfyR3PAyCQkxbUaPupBXgknl8w"
}
```  
Otherwise an error will be returned if the secret is missing, does not match the password hash, or is invalid.
```json
{
    "status": false,
    "error": {
        "code": 304,
        "error": "Unauthorized operation",
        "field": "secret"
    }
}
```
For added convenience, it is also possible to request a new token before the current one expires by sending an authorized request to `/chats/<room_id>/token/renew`.

## Chatting ##
Chatting is implemented using WebSockets through the `chats/<room_id>/ws` endpoint. The steps are:
  1. Send a HTTP POST to `/chats/<room_id>/token` to obtain [authorization](#authorization) (only required for private and hidden rooms)
  2. Open WebSocket to `/chats/<room_id>/ws` with the auth token set in the `Sec-WebSocket-Protocol` header
  3.  Broadcast messages in the following format:
```
    {
        "event_type": "join/leave/send",
        "name": "my-username",
        "color": "color of chat message",
        "msg": "message to send"
    }
```
It is recommended to send a `join` event broadcast when joining chat rooms so everyone else is notified. The server will send a "leave" event upon disconnect

  4. Send a HTTP GET to `/chats/<room_id>/token/renew` to renew the authorization token before expiration.
