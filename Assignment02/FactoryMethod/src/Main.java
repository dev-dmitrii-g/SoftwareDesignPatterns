import internal.MagbyFactory;
import internal.Pokemon;
import internal.PokemonFactory;
import internal.SlowbroFactory;

void main() {
    PokemonFactory factory = new MagbyFactory();
    Pokemon magby = factory.spawn();

    factory = new SlowbroFactory();
    Pokemon slowbro = factory.spawn();

    magby.attack(slowbro);
    slowbro.attack(magby);

    magby.attack(slowbro);
    slowbro.attack(magby);
}