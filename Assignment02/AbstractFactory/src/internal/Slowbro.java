package internal;

class Slowbro extends Pokemon {
    Slowbro() {
        super("Slowbro", Type.WATER, 95, 75);
    }

    @Override
    public void speak() {
        System.out.println(this.name + " says: When a SLOWPOKE went hunting in the sea, its tail was bitten by a SHELLDER. That made it evolve into SLOWBRO!\n");
    }
}
