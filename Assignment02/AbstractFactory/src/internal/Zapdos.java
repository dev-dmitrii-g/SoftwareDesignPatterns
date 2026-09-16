package internal;

class Zapdos extends Pokemon {
    Zapdos() {
        super("Zapdos", Type.ELECTRIC, 90, 90);
    }

    @Override
    public void speak() {
        System.out.println(this.name + " says: ZAPDOS is a legendary bird POKéMON that has the ability to control electricity. It usually lives in thunderclouds. It gains power if it is stricken by lightning bolts!\n");
    }
}
