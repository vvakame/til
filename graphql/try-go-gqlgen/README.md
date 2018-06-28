# README

```
$ vgo generate ./...
$ vgo run main.go
```

```
mutation createTodo {
  createTodo(text: "test") {
    user {
      id
    }
    text
    done
  }
}

query getTodos {
  hoge: todos {
    id
    text
    piyo: done
    user {
      id
      name
    }
  }
}

query getNode {
  user: node(id: "U5577006791947779410") {
    id
    ... on User {
      name
    }
  }
  todo: node(id: "T6129484611666145821") {
    id
    ... on Todo {
      text
      done
      user {
        id
        name
      }
    }
  }
}

query getNodes {
  nodes(ids: ["U5577006791947779410", "T6129484611666145821"]) {
    id
    ... on User {
      name
    }
    ... on Todo {
      text
      done
      user {
        id
        name
      }
    }
  }
}
```