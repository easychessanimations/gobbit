# gobbit

Go language bitboard multi variant chess analysis engine.

# Motivation

- When you want to analyze a variant it is tempting to modify an existing engine, but this carries risks and difficulties. I wanted an engine that I **fully understand** and is **general enough** so that it can accomodate **many variants**.

- **Portability** is important for me. Go is a very portable language, building for a given platform only takes setting two environment variables before running the build command ( the same is not true for C/C++, building for Linux under Windows is close to impossible ). Go is also a pretty low level language ( though not as low level as C/C++ ), so it allows creating a fast native executable. Go compiler has a WASM port, which allows running the engine in the browser.

- The engine was inspired by an existing Go language bitboard engine [Zurichess](https://bitbucket.org/zurichess/zurichess/src/master/) which is listed on [CCRL](https://ccrl.chessdom.com/ccrl/4040/cgi/engine_details.cgi?print=Details&each_game=1&eng=Zurichess%20Neuchatel%2064-bit) and is rated 2800 there. Many ideas were taken over, however not slavishly.

# Variants

Supported variants are Standard, [8-Piece](https://www.chessvariants.com/rules/8-piece-chess) and Atomic.

# Protocol

The engine operates on a useful fraction of the [UCI protocol](http://wbec-ridderkerk.nl/html/UCIProtocol.html). Only analysis features are supported, the engine has no time management and it cannot be used to play a live game.

# Online

The WASM build of the engine is available online at

[https://gobbitengine.netlify.app/](https://gobbitengine.netlify.app/)

# Binaries

Precompiled binaries are available in the `dist` folder.

## Windows binary

`dist\gobbit.exe`

## Linux binary

`dist\gobbit`

# Building

Clone the repository:

```
git clone https://github.com/easychessanimations/gobbit.git
```

In `go.mod` change `C:/gomodules/modules/gobbit` to the absolute path of the folder where you cloned the repository.

## Building on Windows

In `s\b.bat` change `C:/gomodules/modules/gobbit` to the absolute path of the folder where you cloned the repository.

Then run the build script:

```
s\b
```

This will create the Windows and Linux executable in the `dist` folder and the web site, including the WASM executable and the glue code, in the `site` folder.

## Building on other systems

Study the build instructions for Windows, and convert the Windows build script to a script that your system understands.

# Discussion

The engine can be discussed live at

[easychess Discord](https://discord.gg/RKJDzJj)
