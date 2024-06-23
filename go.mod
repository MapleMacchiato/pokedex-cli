module github.com/MapleMacchiato/pokedex-cli

replace github.com/MapleMacchiato/pokedex-cli/internal/pokeapi v0.0.0 => ./internal/pokeapi/

require (
	github.com/MapleMacchiato/pokedex-cli/internal/pokeapi v0.0.0
)

go 1.22.3
