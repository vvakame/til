// --allowUmdGlobalAccess 指定無しの場合次のエラーになる
// error TS2686: 'foobar' refers to a UMD global, but the current file is a module. Consider adding an import instead.
console.log(foobar.name);

export { }
