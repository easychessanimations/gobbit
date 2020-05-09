# gobbit

Go language bitboard multi variant chess analysis engine.

# Motivation

- When you want to analyze a variant it is tempting to modify an existing engine, but this carries risks and difficulties. I wanted an engine that I fully understand and is general enough so that it can accomodate many variants.

- Portablity is important for me. Go is a very portable language, building for a given platform takes setting two environment variables before running the build command ( the same is not true for C/C++, building for Linux under Windows is close to impossible ). Go is also pretty low level ( though not as low as C/C++), so it allows creating a fast native executable. Go compiler supports WASM, so porting the engine to the browser is easy.

- The engine was inspired by an existing bitboard Go language engine [Zurichess](https://bitbucket.org/zurichess/zurichess/src/master/) which is listed on [CCRL](https://ccrl.chessdom.com/ccrl/4040/cgi/engine_details.cgi?print=Details&each_game=1&eng=Zurichess%20Neuchatel%2064-bit) and rated 2800 there. Many ideas were taken over, however not slavishly.

# Variants

Supported variants are Standard, [8-Piece](https://www.chessvariants.com/rules/8-piece-chess) and Atomic.

# Online

The WASM build of the engine is available online at

[https://gobbitengine.netlify.app/](https://gobbitengine.netlify.app/)

# Discussion

The engine can be discussed live at

[easychess Discord](https://discord.gg/RKJDzJj)