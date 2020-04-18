const fs = require('fs')

const NUM_RANKS = 8
const NUM_FILES = 8

const rankLetters = ["1", "2", "3", "4", "5", "6", "7", "8"]
const fileLetters = ["a", "b", "c", "d", "e", "f", "g", "h"]

function writeFile(name, package, content){
    let pcontent = `package ${package}

// generated by gen.js, don't edit
${content}
`

    fs.writeFileSync(`${package}/${name}.go`, pcontent)
}

let sqs = []

for(let rank=0; rank<NUM_RANKS; rank++) for(let file=0; file<NUM_FILES; file++) sqs.push(
    `   Square${fileLetters[file].toUpperCase()}${rankLetters[rank]} = Square(${rank*NUM_RANKS+file})`
)

let ranks = []

for(let rank=0; rank<NUM_RANKS; rank++) ranks.push(
    `   Rank${rankLetters[rank]} = Rank(${rank})`
)

let files = []

for(let file=0; file<NUM_FILES; file++) files.push(
    `   File${fileLetters[file].toUpperCase()} = File(${file})`
)

const square_go = `
const NUM_RANKS = ${NUM_RANKS}
const LAST_RANK = NUM_RANKS - 1
const NUM_FILES = ${NUM_RANKS}
const LAST_FILE = NUM_FILES - 1

const BOARD_AREA = NUM_RANKS * NUM_FILES

const RANK_STORAGE_SIZE_IN_BITS = 3
const FILE_STORAGE_SIZE_IN_BITS = 3

const SQUARE_STORAGE_SIZE_IN_BITS = RANK_STORAGE_SIZE_IN_BITS + FILE_STORAGE_SIZE_IN_BITS

const RANK_SHIFT_IN_BITS = FILE_STORAGE_SIZE_IN_BITS
const FILE_SHIFT_IN_BITS = 0

const RANK_MASK = ((1 << RANK_STORAGE_SIZE_IN_BITS) - 1) << RANK_SHIFT_IN_BITS
const FILE_MASK = (1 << FILE_STORAGE_SIZE_IN_BITS) - 1

type Rank int
type File int
type Square uint

const (
${sqs.join("\n")}
)

const SquareMinValue = SquareA1
const SquareMaxValue = SquareH8

const (
${ranks.join("\n")}
)

const (
${files.join("\n")}
)

var RankLetterOf = [NUM_RANKS]string{${rankLetters.map(rl => `"${rl}"`).join(' , ')}}
var FileLetterOf = [NUM_FILES]string{${fileLetters.map(rl => `"${rl}"`).join(' , ')}}

var UCIOf [BOARD_AREA]string
var UCIToSquare map[string]Square
`

writeFile("square", "basic", square_go)

const figures = [
    ["NoFigure", "."],    
    ["Pawn", "p"],
    ["Knight", "n"],
    ["Bishop", "b"],
    ["Rook", "r"],
    ["Queen", "q"],
    ["King", "k"],
    ["Lancer", "l"],
    ["LancerN", "ln"],
    ["LancerNE", "lne"],
    ["LancerE", "le"],
    ["LancerSE", "lse"],
    ["LancerS", "ls"],
    ["LancerSW", "lsw"],
    ["LancerW", "lw"],
    ["LancerNW", "lnw"],
    ["Sentry", "s"],
    ["Jailer", "j"]
]

const piece_go = `
type Figure int

const (
${figures.map((p,i) => "   " + p[0].padEnd(16, " ") + " = Figure(" + i + ")").join("\n")}
)

const FigureMinValue = ${figures[1][0]}
const FigureMaxValue = ${figures[figures.length-1][0]}
const FigureArraySize = FigureMaxValue - FigureMinValue + 1

const LancerMinValue = LancerN
const LancerMaxValue = LancerNW

const LANCER_DIRECTION_MASK = 0b111

// SymbolOf tells the symbol of a Figure
var SymbolOf = [${figures.length}]string{
${figures.map((fig,i) => ('   "' + fig[1] + '"').padEnd(10, " ") + ', // ' + (""+i).padEnd(4, " ") + fig[0]).join("\n")}
}

type Piece int

const (
   NoPiece               = Piece(0)
   DummyPiece            = Piece(1)
${figures.slice(1).map((p,i) =>
    "   " + ("Black"+p[0]).padEnd(21, " ") + " = Piece(" + (2*i + 2) + ")\n" +
    "   " + ("White"+p[0]).padEnd(21, " ") + " = Piece(" + (2*i + 3) + ")"
).join("\n")}
)

const PieceMinValue = Black${figures[1][0]}
const PieceMaxValue = White${figures[figures.length-1][0]}
const PieceArraySize = PieceMaxValue - PieceMinValue + 1

var FigureOf [${figures.length*2}]Figure

type Color int

const (   
   Black   =  Color(0) 
   White   =  Color(1)
   NoColor =  Color(2)
)

const COLOR_MASK = White

const ColorMinValue = 0
const ColorMaxValue = NoColor
const ColorArraySize = ColorMaxValue - ColorMinValue + 1

// ColorOf tells the color of a Piece
var ColorOf [${figures.length*2}]Color

// ColorFigure constructs a Piece from Color and Figure
var ColorFigure[2][${figures.length}]Piece

// SymbolToPiece tells Piece for a FEN symbol
var SymbolToPiece = map[string]Piece{
${figures.slice(1).map((fig,i) => {
    return  `   "${fig[1]}"`.padEnd(10, " ") + `: Black${fig[0]},`.padEnd(20, " ") + "// " + (i*2+2) + "\n" + 
            `   "${fig[1].substring(0,1).toUpperCase()+fig[1].substring(1)}"`.padEnd(10, " ") + `: White${fig[0]},`.padEnd(20, " ") + "// " + (i*2+3)
}).join("\n")}
}
`

writeFile("piece", "basic", piece_go)