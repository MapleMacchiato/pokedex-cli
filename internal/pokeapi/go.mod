module github.com/MapleMacchiato/pokedex-cli/internal/pokeapi

replace github.com/MapleMacchiato/pokedex-cli/internal/pokecache v0.0.0 => ../pokecache/

require (
	github.com/MapleMacchiato/pokedex-cli/internal/pokecache v0.0.0
)

go 1.22.3
