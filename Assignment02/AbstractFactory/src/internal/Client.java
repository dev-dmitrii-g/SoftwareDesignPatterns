package internal;

public class Client {
    public Encounter spawn (PokemonFactory factory) {
        Pokemon pokemon = factory.createPokemon();
        SpecialMove specialMove = factory.createSpecialMove();
        Encounter encounter = new Encounter(pokemon, specialMove);

        System.out.println("<\\|/-----------------<\\--\\|/--/>----------------\\|/>\n");
        System.out.println("A wild " + pokemon.getName() + " appeared!\n");
        pokemon.stats();
        System.out.println("Special Move: " + specialMove.getName() + "\n");
        pokemon.speak();
        System.out.println("<\\|/-----------------<\\--\\|/--/>----------------\\|/>\n");
        return encounter;
    }
}
