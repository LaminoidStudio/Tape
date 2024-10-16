wasm_codearray = new Uint8Array(wasm_strbuffer.length);
for (var i in wasm_strbuffer) wasm_codearray[i] = wasm_strbuffer.charCodeAt(i);

const native_console = window.console;
window.console = {
	log: function(str) {
		var node = document.createElement("div");
		node.appendChild(document.createTextNode(str));
		document.getElementById("log").appendChild(node);
		native_console.log(str);
	},

	clear: function(str) {
		document.getElementById("log").innerHTML = '';
		native_console.clear();
	}
}

const go = new Go();
let mod, inst;
WebAssembly.instantiate(wasm_codearray, go.importObject).then((result) => {
	mod = result.module;
	inst = result.instance;
	document.getElementById("runButton").disabled = false;
	document.getElementById("explainButton").disabled = false;
	document.getElementById("helpButton").disabled = false;
}).catch((err) => {
	console.error(err);
});

async function run() {
	// Remember that argv[0] is the program name.
	console.clear();
	go.argv = ['tape', '-run', '-code', document.getElementById('codeTextArea').value, '-memory', document.getElementById('memoryNumber').value, '-seed', document.getElementById('seedNumber').value, '-timeout', document.getElementById('timeoutText').value];
	if (document.getElementById('originalCheckbox').checked)
		go.argv.push('-original');
	if (document.getElementById('signedCheckbox').checked)
		go.argv.push('-signed');
	await go.run(inst);
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}

async function explain() {
	// Remember that argv[0] is the program name.
	console.clear();
	go.argv = ['tape', '-code', document.getElementById('codeTextArea').value];
	if (document.getElementById('originalCheckbox').checked)
		go.argv.push('-original');
	await go.run(inst);
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}

function copy() {
	navigator.clipboard.writeText(document.getElementById('codeTextArea').value);
}

async function help() {
	// Remember that argv[0] is the program name.
	console.clear();
	go.argv = ['tape', '-help'];
	await go.run(inst);
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
}

function reset() {
	console.clear();
	document.getElementById('codeTextArea').value = '';
}
