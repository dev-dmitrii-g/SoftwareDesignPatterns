package internal;

public class WaterFactory extends PokemonFactory {
    @Override
    public Pokemon createPokemon() {
        return new Slowbro();
    }

    @Override
    public SpecialMove createSpecialMove() {
        return new WaterPulse();
    }
}
