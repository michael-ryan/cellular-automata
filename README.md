# Cellular Automata Simulator

## 📝 Documentation 

For using the tool, see [here](https://pkg.go.dev/github.com/michael-ryan/cellularautomata).

For constructing your own automata, see [here](https://pkg.go.dev/github.com/michael-ryan/cellularautomata/models).

## 🚀 Usage 
The simplest usage is as follows. In your `main` package, set up a `Config` object and supply it to `Launch`.
```Go
package main

import (
	"github.com/michael-ryan/cellularautomata"
	"github.com/michael-ryan/cellularautomata/examples"
)

func main() {
	c := cellularautomata.Config{
		Fps:          15,
		CellsX:       128,
		CellsY:       72,
		WindowX:      1280,
		WindowY:      720,
		Automaton:    examples.NewForest(),
		InitialState: 0,
		SkipEditor:   true,
	}

	cellularautomata.Launch(c)
}
```

This will open a GUI window and run a simulation.

There is an optional edit mode which the program will start in if `SkipEditor` is `false`. In this mode, you can click on cells to cycle their initial state, then press `S` on your keyboard to start the simulation. 

Note: Launch must be called from the main goroutine, due to a limitation in OpenGL.

Feel free to play with the config parameters. You can swap out `Automaton` for other sample models (defined [here](models/)), or you can easily construct your own automata. For examples, see [here](examples/).

## 🐛 Known Issues & Planned Improvements

- Analysis tools to record cell state counts and how they change over time.

## 🔧 Troubleshooting

### *There's some error about x11, package 'gl' missing etc etc*

See requirements for the graphics library I'm using [here](https://github.com/gopxl/pixel?tab=readme-ov-file#requirements). You'll probably need to install `gcc`, `libgl1-mesa-dev` and `xorg-dev`.

### *My issue isn't listed here!*

Please file an issue on this repository, and I'll take a look.

## 🤝 Contributing

Feel free to file pull requests.
