package internal;

public class FireFactory extends PokemonFactory {
    @Override
    public Pokemon createPokemon() {
        return new Magby();
    }

    @Override
    public SpecialMove createSpecialMove() {
        return new Ember();
    }
}
