# rest-api

# user-service

# REST API


GET /users -- list of users -- 200, 404, 500  
POST /users/:id --create user -- 201, 4xx, Header Location: url  
GET /users/:id -- user by id -- 200, 404, 500  
PUT /users/:id -- fully update user -- 204/200, 400, 404, 500  
PATCH /user/:id -- partially update user -- 204/200, 400, 404, 500  
DELETE /users/:id -- delete user by id -- 204/200, 400, 404  