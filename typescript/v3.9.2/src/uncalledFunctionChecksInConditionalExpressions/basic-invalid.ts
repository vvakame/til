function hasImportantPermissions(): boolean {
    return true;
}

function deleteAllTheImportantFiles() {
}

// 関数の呼び出しを忘れて条件式に使った場合、エラーになる (TypeScript v3.7 で入った)
// error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
if (hasImportantPermissions) {
    deleteAllTheImportantFiles();
}

// 三項演算子でもこのチェックが働くようになった (今回から)
// error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
hasImportantPermissions ? deleteAllTheImportantFiles() : void 0;

// なお、 Add missing call parentheses のQuick Fixも追加されたので、エラーになっている箇所を見つけたらｼｭｯと直せる
