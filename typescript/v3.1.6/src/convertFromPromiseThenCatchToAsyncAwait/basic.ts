function timeout(timeout?: number) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve();
    }, timeout);
  });
}

function exec() {
  return timeout(100)
    .then(() => {
      if (Math.random() > 0.5) {
        throw new Error("random failed");
      }
      console.log("hi!");
    })
    .catch(err => {
      console.error(err);
    });
}
