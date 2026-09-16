package internal;

public class SlowbroFactory extends PokemonFactory {
    @Override
    Pokemon createPokemon() {
        return new Slowbro();
    }
}
