package internal;

public class WaterPulse extends SpecialMove {
    private final String name = "WATER PULSE";
    private final int power = 60;

    @Override
    public void execute(Pokemon user, Pokemon target) {
        if (user.fainted || target.fainted) {
            System.out.println((user.fainted ? user.name : target.name) + " is already fainted!\n");
            return;
        }

        System.out.println("<\\-MOVE-/> " + user.name + " struck " + target.name + " with WATER PULSE!\n");
        target.takeDamage(this.power);
    }

    @Override
    public String getName() {
        return this.name;
    }
}