
# RESTful API for Task Management

A simple CRUD API using Gin with PostgreSQL and JWT authentication.


## Installation

1 . Clone the repository:

```bash
  git clone https://github.com/casual-user-asm/task-manager.git
  cd task-manager
```
2 . Set up the environment variables in .env file(take the values from .env.example file)

3 . Build and run the Docker containers:

```bash
  docker-compose up -d
```
## API Endpoints

#### Register a new user.

```
  POST /user/register
  
  Example fields for JSON:
  
  {
    "username": "test",
    "email": "test",
    "password": "test",
  }
```


#### Log in with registered user credentials and receive a JWT token.

```
  POST /user/login

  Example fields for JSON:
  
  {
    "email": "test",
    "password": "test",
  }
```

#### Log out and invalidate the JWT token.
```
  PUT /user/logout
```


#### Invalidate the JWT token and Delete User.

```
  DELETE /user/delete
```
#### Get a single task by ID.
```
  GET /tasks/:id
```


#### Get all tasks.

```
  GET /tasks/
```
#### Create a new task.
```
  POST /tasks/create

  Example fields for JSON:
  
  {
    "title": "test",
    "description": "test",
  }
```


#### Update an existing task by ID.

```
  PUT /tasks/update/:id

  Example fields for JSON:
  
  {
    "title": "test",
    "description": "test",
  }
```
#### Delete a task by ID.
```
  DELETE /tasks/delete/:id
```


