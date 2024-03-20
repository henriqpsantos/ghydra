# Ghydra

A menu-based runner for the commands you don't want to type and can't remember the key-bindings

Ghydra is a terminal-based program to run other commands through user-defined menus. Configuration is based on toml. Ideally, it should hit a sweet spot between the speed of a single key-binding and the ease-of-use of a GUI-like interface.

When Ghydra executes, an interface with a set of keybinds and their descriptions shows up. The user is then prompted to select an option or terminate the program with `C-c`.

## Configuration

Ghydra uses a toml based configuration with the following style:

```toml
# You can also configure the colors from the top-level menu
keyColor    = "#A0D822" # This is the color that the key within [] is shown
sepColor    = "#888888" # Color for the brackets
descColor   = "#A0AAEE" # Color for the description text
headerColor = "#A0AAEE"

[[menu]] # This is a top-level menu
desc = "The menu's description"
key = "k"

[[menu.action]]
desc = "This action is inside the 'k' menu"
key = "r"
command = "echo hello"

[[menu.action]]
desc = "This action is also inside the 'k' menu"
key = "w"
command = "echo world" 

[[menu.menu]] # You can define menus within menus
desc = "This is a menu within a menu"
key = "m"

[[menu]] # Another top-level menu
desc = "This menu is on the root level"
key = "l"

...

```

As a tree, the menu will look something like this:

```
root
 ├── k    > The menu's description 
 │   ├─ r > This action is inside the 'k' menu
 │   ├─ w > This action is also inside the 'k' menu
 │   └─ m > This is a menu within a menu
 │      ├ ...
 │     ...
 │
 └── l    > This menu is on the root level
```

## Shell support

Ghydra will output the interface to `stderr` and output the result of whichever `command` is selected to `stdout`. To run the command output by Ghydra in fish:

```fish
eval (ghydra)
````

To setup a keybind to execute Ghydra and run the selected program:

```fish
	bind \cx 'eval (ghydra
)'
```

## Acknowledgements

Though the author has never used tydra or hydra, this is inspired by the tydra/hydra system.
