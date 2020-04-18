package basic

// generated by gen.js, don't edit

type Figure int

const (
   NoFigure         = Figure(0)
   Pawn             = Figure(1)
   Knight           = Figure(2)
   Bishop           = Figure(3)
   Rook             = Figure(4)
   Queen            = Figure(5)
   King             = Figure(6)
   Lancer           = Figure(7)
   LancerN          = Figure(8)
   LancerNE         = Figure(9)
   LancerE          = Figure(10)
   LancerSE         = Figure(11)
   LancerS          = Figure(12)
   LancerSW         = Figure(13)
   LancerW          = Figure(14)
   LancerNW         = Figure(15)
   Sentry           = Figure(16)
   Jailer           = Figure(17)
)

const FigureMinValue = Pawn
const FigureMaxValue = Jailer
const FigureArraySize = FigureMaxValue - FigureMinValue + 1

const LancerMinValue = LancerN
const LancerMaxValue = LancerNW

const LANCER_DIRECTION_MASK = 0b111

// SymbolOf tells the symbol of a Figure
var SymbolOf = [18]string{"." , "p" , "n" , "b" , "r" , "q" , "k" , "l" , "ln" , "lne" , "le" , "lse" , "ls" , "lsw" , "lw" , "lnw" , "s" , "j"}

type Piece int

const (
   NoPiece               = Piece(0)
   DummyPiece            = Piece(1)
   BlackPawn             = Piece(2)
   WhitePawn             = Piece(3)
   BlackKnight           = Piece(4)
   WhiteKnight           = Piece(5)
   BlackBishop           = Piece(6)
   WhiteBishop           = Piece(7)
   BlackRook             = Piece(8)
   WhiteRook             = Piece(9)
   BlackQueen            = Piece(10)
   WhiteQueen            = Piece(11)
   BlackKing             = Piece(12)
   WhiteKing             = Piece(13)
   BlackLancer           = Piece(14)
   WhiteLancer           = Piece(15)
   BlackLancerN          = Piece(16)
   WhiteLancerN          = Piece(17)
   BlackLancerNE         = Piece(18)
   WhiteLancerNE         = Piece(19)
   BlackLancerE          = Piece(20)
   WhiteLancerE          = Piece(21)
   BlackLancerSE         = Piece(22)
   WhiteLancerSE         = Piece(23)
   BlackLancerS          = Piece(24)
   WhiteLancerS          = Piece(25)
   BlackLancerSW         = Piece(26)
   WhiteLancerSW         = Piece(27)
   BlackLancerW          = Piece(28)
   WhiteLancerW          = Piece(29)
   BlackLancerNW         = Piece(30)
   WhiteLancerNW         = Piece(31)
   BlackSentry           = Piece(32)
   WhiteSentry           = Piece(33)
   BlackJailer           = Piece(34)
   WhiteJailer           = Piece(35)
)

const PieceMinValue = BlackPawn
const PieceMaxValue = WhiteJailer
const PieceArraySize = PieceMaxValue - PieceMinValue + 1

var FigureOf [36]Figure

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
var ColorOf [36]Color

// ColorFigure constructs a Piece from Color and Figure
var ColorFigure[2][18]Piece

// SymbolToPiece tells Piece for a FEN symbol
var SymbolToPiece = map[string]Piece{
   "p"    : BlackPawn,
   "P"    : WhitePawn,
   "n"    : BlackKnight,
   "N"    : WhiteKnight,
   "b"    : BlackBishop,
   "B"    : WhiteBishop,
   "r"    : BlackRook,
   "R"    : WhiteRook,
   "q"    : BlackQueen,
   "Q"    : WhiteQueen,
   "k"    : BlackKing,
   "K"    : WhiteKing,
   "l"    : BlackLancer,
   "L"    : WhiteLancer,
   "ln"   : BlackLancerN,
   "Ln"   : WhiteLancerN,
   "lne"  : BlackLancerNE,
   "Lne"  : WhiteLancerNE,
   "le"   : BlackLancerE,
   "Le"   : WhiteLancerE,
   "lse"  : BlackLancerSE,
   "Lse"  : WhiteLancerSE,
   "ls"   : BlackLancerS,
   "Ls"   : WhiteLancerS,
   "lsw"  : BlackLancerSW,
   "Lsw"  : WhiteLancerSW,
   "lw"   : BlackLancerW,
   "Lw"   : WhiteLancerW,
   "lnw"  : BlackLancerNW,
   "Lnw"  : WhiteLancerNW,
   "s"    : BlackSentry,
   "S"    : WhiteSentry,
   "j"    : BlackJailer,
   "J"    : WhiteJailer,
}

