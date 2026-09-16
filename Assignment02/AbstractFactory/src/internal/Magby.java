package internal;

class Magby extends Pokemon {
    Magby() {
        super("Magby", Type.FIRE, 45, 75);
    }

    @Override
    public void speak() {
        System.out.println(this.name + " says: When flames drip from its nose, that means it has a cold. Have it lie down for a nice rest in some magma!\n");
    }
}
