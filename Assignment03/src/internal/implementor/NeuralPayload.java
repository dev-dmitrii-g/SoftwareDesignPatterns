package internal.implementor;

public class NeuralPayload {
    public enum TargetRegion {
        OPTIC_NERVE,
        MOTOR_CORTEX,
        LIMBIC_SYSTEM
    }

    private final TargetRegion targetRegion;
    private final int intensity;
    private final String commandText;

    public NeuralPayload(TargetRegion targetRegion, int intensity, String commandText) {
        this.targetRegion = targetRegion;
        this.intensity = intensity;
        this.commandText = commandText;
    }

    public TargetRegion getTargetRegion() {
        return this.targetRegion;
    }

    public int getIntensity() {
        return this.intensity;
    }

    public String getCommandText() {
        return this.commandText;
    }
}
