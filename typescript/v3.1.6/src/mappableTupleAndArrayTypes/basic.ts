type Stringify<T> = { [K in keyof T]: string };

function stringifyProps<T>(v: T): Stringify<T> {
  const result = {} as Stringify<T>;
  for (const prop in v) {
    result[prop] = String(v[prop]);
  }
  return result;
}

{
  const obj: { no: string } = stringifyProps({ no: 151 });
  console.log(obj);
}

function stringifyAll<T extends unknown[]>(...args: T): Stringify<T> {
  return args.map(v => String(v)) as any;
}

{
  const array = stringifyAll([1, true]);
  // TypeScript 3.1以前だと forEach も length も string と解釈される
  // 一般的には配列の要素部分だけ変換されれば十分だよなぁ…？
  array.forEach(v => console.log(v));
  const len: number = array.length;
  console.log(array, len);
}
{
  const tuple: [1, true] = [1, true];
  const array = stringifyAll(tuple);
  // TypeScript 3.1以前だと forEach も length も string と解釈される
  array.forEach(v => console.log(v));
  const len: number = array.length;
  console.log(tuple, array, len);
}
