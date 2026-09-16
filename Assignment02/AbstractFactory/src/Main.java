import internal.*;

void main() {
    Client client = new Client();
    PokemonFactory factory = new FireFactory();
    Encounter fireEncounter = client.spawn(factory);
    factory = new ElectricFactory();
    Encounter electricEncounter = client.spawn(factory);
    fireEncounter.getSpecialMove().execute(fireEncounter.getPokemon(), electricEncounter.getPokemon());
    electricEncounter.getSpecialMove().execute(electricEncounter.getPokemon(), fireEncounter.getPokemon());
    fireEncounter.getPokemon().stats();
    electricEncounter.getPokemon().stats();
    fireEncounter.getPokemon().attack(electricEncounter.getPokemon());
    electricEncounter.getPokemon().attack(fireEncounter.getPokemon());
    fireEncounter.getPokemon().stats();
    electricEncounter.getPokemon().stats();
    fireEncounter.getSpecialMove().execute(fireEncounter.getPokemon(), electricEncounter.getPokemon());
    electricEncounter.getSpecialMove().execute(electricEncounter.getPokemon(), fireEncounter.getPokemon());
}