#!/bin/sh
GOOS=js GOARCH=wasm go build -o tape.wasm ..
base64 -i tape.wasm | tr -d '\n' > tape.wasmb64
printf "const wasm_strbuffer = atob(\"" > tape.js
cat tape.wasmb64 >> tape.js
printf "\");\\n" >> tape.js
cat tape_template.js >> tape.js
cp tape.html index.html
