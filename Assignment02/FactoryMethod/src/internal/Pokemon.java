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

    public abstract void speak();
    public abstract void attack(Pokemon pokemon);

    protected void takeDamage(int damage) {
        if (this.fainted) {
            System.out.println(this.name + " is already fainted!\n");
            return;
        }
        this.currentHp -= damage;
        if (this.currentHp <= 0) {
            this.currentHp = 0;
            faint();
        }
    }

    private void faint() {
        System.out.println(this.name + " has fainted!\n");
        this.fainted = true;
    }

    public void stats() {
        System.out.println(
                "NAME: " + this.name + "\n"
                + "TYPE: " + this.type + "\n"
                + "HP: " + this.currentHp + "/" + this.maxHp + "\n"
                + "ATK: " + this.attack + "\n"
        );
    }
}
