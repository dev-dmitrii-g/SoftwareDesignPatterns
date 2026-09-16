package internal;

public class ElectricFactory extends PokemonFactory {
    @Override
    public Pokemon createPokemon() {
        return new Zapdos();
    }

    @Override
    public SpecialMove createSpecialMove() {
        return new ThunderShock();
    }
}
