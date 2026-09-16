package internal;

public class MagbyFactory extends PokemonFactory {
    @Override
    Pokemon createPokemon() {
        return new Magby();
    }
}
