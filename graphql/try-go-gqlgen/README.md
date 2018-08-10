# README

```
$ vgo generate ./...
$ vgo run ./server/server.go
```

```
query QueryTodos {
  todos {
    id
    text
    done
    user {
      id
      name
    }
  }
}

mutation CreateTodo {
  createTodo(input: {text: "do performance", userId: "User:123"}) {
    id
    text
    done
    user {
      id
      name
    }
  }
}
```
