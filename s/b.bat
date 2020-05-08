echo off

set GOOS=js

set GOARCH=wasm

go build -o main.wasm -tags wasm

copy \go\misc\wasm\*.js
type wasm_exec.js loadwasm.js > wasm_loader.js
rem copy \go\misc\wasm\*.html

rem copy wasm_exec.html index.html

copy index.html site
copy main.wasm site
copy wasm_exec.js site
copy wasm_loader.js site
