
# ui

A minimal terminal UI library for Go — fuzzy menus, prompts, and multi-select, with no dependencies beyond `golang.org/x/term`.

## Install

```bash
go get github.com/y-d-n-a/ui
```

## Usage

### Fuzzy Menu

Displays a searchable, single-select menu. Arrow keys navigate, typing filters, Enter confirms.

```go
selected, err := ui.Run("Pick a fruit:", []ui.Option{
    {Label: "Apple", Value: "apple"},
    {Label: "Banana", Value: "banana"},
    {Label: "Cherry", Value: "cherry"},
})
if err != nil {
    log.Fatal(err)
}
fmt.Println(selected.Value)
```

### Multi-Select Menu

Space toggles selection, Enter confirms. If nothing is toggled, Enter selects the current item.

```go
selected, err := ui.RunMulti("Pick your toppings:", []ui.Option{
    {Label: "Cheese", Value: "cheese"},
    {Label: "Peppers", Value: "peppers"},
    {Label: "Olives", Value: "olives"},
})
if err != nil {
    log.Fatal(err)
}
for _, opt := range selected {
    fmt.Println(opt.Value)
}
```

### Text Prompt

```go
name, err := ui.Prompt("Your name:")
if err != nil {
    log.Fatal(err)
}
fmt.Println(name)
```

### Secret Prompt

Input is masked with `*` characters.

```go
pass, err := ui.PromptSecret("Password:")
if err != nil {
    log.Fatal(err)
}
```

## API

```go
type Option struct {
    Label string // displayed in the menu
    Value string // returned on selection
}

func Run(prompt string, options []Option) (Option, error)
func RunMulti(prompt string, options []Option) ([]Option, error)
func Prompt(label string) (string, error)
func PromptSecret(label string) (string, error)
```

## Keybindings

| Key        | Action              |
|------------|---------------------|
| `↑` / `↓` | Move cursor         |
| Type       | Filter options      |
| `Space`    | Toggle (multi only) |
| `Enter`    | Confirm             |
| `Esc` / `Ctrl+C` | Cancel        |

## Requirements

- Go 1.19+
- A real terminal (raw mode via `golang.org/x/term`)
