*shortcuts* for easy usage of gobbit from the command line

**command aliases**

the engine defines command aliases for certain frequently used uci commands

<https://github.com/easychessanimations/gobbit/blob/c01d8bcd5f9896fdbd3b157ca17d36e972de3444/maincommon.go#L6>

```
var UCI_COMMAND_ALIASES = map[string]string{
    "vs" : "setoption name UCI_Variant value Standard",
    "ve" : "setoption name UCI_Variant value Eightpiece",
    "va" : "setoption name UCI_Variant value Atomic",
}
```

notably to set the variant to Eightpiece, type

```
ve
```

**engineconfig.txt**

you can save a text file called `engineconfig.txt` in the directory where `gobbit` executabale resides

uci commands and gobbit commands / command aliases listed in this file will be executed upon engine startup ( empty lines are not allowed )

so if you are regularly analyzing a position, or playing a game, you can set this up in `engineconfig.txt` like

```
ve
position startpos moves g1f3 d7d5
```

which would set the variant to Eightpiece and play the moves 1. Nf3 d5 on every engine startup

**matein4.txt**

if you save

<https://raw.githubusercontent.com/easychessanimations/gobbit/master/matein4.txt>

as `matein4.txt` in the engine directory

then typing `u` will give you a random mate in 4 puzzle; to start solving it, type `g`