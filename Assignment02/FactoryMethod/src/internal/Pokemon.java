package internal;

public abstract class Pokemon {
    protected final String name;
    protected final Type type;
    protected final int maxHp;
    protected int currentHp;
    protected final int attack;
    protected boolean fainted = false;

    Pokemon(String name, Type type, int maxHp, int attack) {
        this.name = name;
        this.type = type;
        this.currentHp = this.maxHp = maxHp;
        this.attack = attack;
    }

    abstract void speak();
    abstract void attack(Pokemon pokemon);

    void takeDamage(int damage) {
        if (this.fainted) {
            System.out.println(this.name + " has fainted!");
            return;
        }
        this.currentHp -= damage;
        if (this.currentHp <= 0) {
            this.currentHp = 0;
            this.fainted = true;
        }
    }

    void stats() {
        System.out.println(
                "NAME: " + this.name + "\n"
                + "TYPE: " + this.type + "\n"
                + "HP: " + this.maxHp + "/" + this.currentHp + "\n"
                + "ATK: " + this.attack
        );
    }
}
