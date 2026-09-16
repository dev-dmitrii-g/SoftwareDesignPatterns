package internal;

public class ThunderShock extends SpecialMove {
    private final String name = "THUNDER SHOCK";
    private final int power = 40;

    @Override
    public void execute(Pokemon user, Pokemon target) {
        if (user.fainted || target.fainted) {
            System.out.println((user.fainted ? user.name : target.name) + " is already fainted!\n");
            return;
        }

        System.out.println("<\\-MOVE-/> " + user.name + " zapped " + target.name + " with THUNDER SHOCK!\n");
        target.takeDamage(this.power);
    }

    @Override
    public String getName() {
        return this.name;
    }
}