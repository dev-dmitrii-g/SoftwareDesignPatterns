package internal;

public class Encounter {
    private final Pokemon pokemon;
    private final SpecialMove specialMove;

    public Encounter(Pokemon pokemon, SpecialMove specialMove) {
        this.pokemon = pokemon;
        this.specialMove = specialMove;
    }

    public void performSpecialMove(Encounter targetEncounter) {
        this.specialMove.execute(this.pokemon, targetEncounter.getPokemon());
    }

    public Pokemon getPokemon() {
        return this.pokemon;
    }

    public SpecialMove getSpecialMove() {
        return this.specialMove;
    }
}