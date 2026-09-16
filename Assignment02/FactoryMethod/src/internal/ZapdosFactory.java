package internal;

public class ZapdosFactory extends PokemonFactory {
    @Override
    Pokemon createPokemon() {
        return new Zapdos();
    }
}
