package internal;

public class Ember extends SpecialMove {
    private final String name = "EMBER";
    private final int power = 40;

    @Override
    public void execute(Pokemon user, Pokemon target) {
        if (user.fainted || target.fainted) {
            System.out.println((user.fainted ? user.name : target.name) + " is already fainted!\n");
            return;
        }

        System.out.println("<\\-MOVE-/> " + user.name + " attacks " + target.name + " with EMBER!\n");
        target.takeDamage(this.power);
    }

    @Override
    public String getName() {
        return this.name;
    }
}
