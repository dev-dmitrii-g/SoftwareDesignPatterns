package internal;

public abstract class PokemonFactory {
    abstract Pokemon createPokemon();

    public Pokemon spawn() {
        Pokemon pokemon = createPokemon();
        System.out.println("<\\|/-----------------<\\--\\|/--/>----------------\\|/>\n");
        System.out.println("A wild " + pokemon.name + " appeared!\n");
        pokemon.stats();
        pokemon.speak();
        System.out.println("<\\|/-----------------<\\--\\|/--/>----------------\\|/>\n");
        return pokemon;
    }
}
