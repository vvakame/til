// JSXのテキストで、JSXの仕様的に使ってはいけない > とか } の利用がコンパイルエラーになるようになった
// error TS1382: Unexpected token. Did you mean `{'>'}` or `&gt;`?
// error TS1381: Unexpected token. Did you mean `{'}'}` or `&rbrace;`?
let node = (
    <div>
        テキストに > とか } があったらダメ！
        これからは {">"} とか {"}"} を使うか、
        &gt; や &rbrace; を使おう！
    </div>
);

// Quick Fix も2パターン追加されている
// Wrap invalid character in an expression container
// Convert invalid character to its html entity code
// …が、 v3.9.2 + VSCode 1.46.0-insider の組み合わせだと選んでも直してくれないぽい
