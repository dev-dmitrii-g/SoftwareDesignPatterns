package internal;

import java.util.concurrent.ThreadLocalRandom;

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
    public void attack(Pokemon pokemon) {
        if (this.fainted) {
            return;
        }
        System.out.println("<\\-MOVE-/> " + this.name + " tackled " + pokemon.name + "!\n");
        pokemon.takeDamage(this.attack / 3 + ThreadLocalRandom.current().nextInt(0, 10));
    };

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

    public boolean isFainted() {
        return this.fainted;
    }

    public void stats() {
        System.out.println(
                "NAME: " + this.name + "\n"
                + "TYPE: " + this.type + "\n"
                + "HP: " + this.currentHp + "/" + this.maxHp + "\n"
                + "ATK: " + this.attack + "\n"
        );
    }

    public String getName() {
        return this.name;
    }
}
