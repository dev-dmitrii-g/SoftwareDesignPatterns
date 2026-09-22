import internal.*;

void main() {
    Client client = new Client();
    PokemonFactory factory = new FireFactory();
    Encounter fireEncounter = client.spawn(factory);
    factory = new ElectricFactory();
    Encounter electricEncounter = client.spawn(factory);
    fireEncounter.performSpecialMove(electricEncounter);
    fireEncounter.getPokemon().stats();
    electricEncounter.getPokemon().stats();
    electricEncounter.performSpecialMove(fireEncounter);
    fireEncounter.getPokemon().stats();
    electricEncounter.getPokemon().stats();
}