if (!WebAssembly.instantiateStreaming) { // polyfill
		WebAssembly.instantiateStreaming = async (resp, importObject) => {
			const source = await (await resp).arrayBuffer();
			return await WebAssembly.instantiate(source, importObject);
		};
}

const go = new Go();
let mod, inst;
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
	mod = result.module;
	inst = result.instance;
	run()
}).catch((err) => {
	console.error(err);
});

async function run() {			
	await go.run(inst);
	inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance			
}

// define a new console
let console=(function(oldCons){
  return {
      log: function(text){
        self.postMessage({text: text})
      },
      info: function (text) {
        self.postMessage({text: text})
      },
      warn: function (text) {
        self.postMessage({text: text})
      },
      error: function (text) {
        self.postMessage({text: text})
      }
  }
}(self.console))

const oldCons = self.console

// redefine the old console
self.console = console

self.addEventListener("message", event => {
	let data = event.data

	let command = data.command

	oldCons.log("received command " + command)

	ExecUciCommandLineWasm(command)        
})
