function ex1(): Promise<boolean> {
  return fetch("https://microsoft.com").then(result => result.ok);
}

function ex2(): Promise<void> {
  return fetch("https://microsoft.com").then(
    result => console.log(result),
    rej => console.log("error", rej)
  );
}

function ex3(): Promise<void> {
  return fetch("https://microsoft.com")
    .then(result => console.log(result))
    .catch(err => console.log(err));
}

export {};
