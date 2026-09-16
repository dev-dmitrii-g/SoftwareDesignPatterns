package internal;

public abstract class PokemonFactory {
    abstract Pokemon createPokemon();

    public Pokemon spawn() {
        Pokemon pokemon = createPokemon();
        System.out.println("A wild " + createPokemon().name + " appeared!");
        pokemon.stats();
        pokemon.speak();
        return pokemon;
    }
}
