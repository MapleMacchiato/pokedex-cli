module github.com/MapleMacchiato/pokedex-cli

replace github.com/MapleMacchiato/pokedex-cli/internal/pokeapi v0.0.0 => ./internal/pokeapi/

replace github.com/MapleMacchiato/pokedex-cli/internal/pokecache v0.0.0 => ./internal/pokecache/

require (
	github.com/MapleMacchiato/pokedex-cli/internal/pokecache v0.0.0
	github.com/MapleMacchiato/pokedex-cli/internal/pokeapi v0.0.0
)

go 1.22.3
