
set GOOS=js

set GOARCH=wasm

go build -o main.wasm -tags wasm

copy \go\misc\wasm\*.js
copy \go\misc\wasm\*.html

rem copy wasm_exec.html index.html