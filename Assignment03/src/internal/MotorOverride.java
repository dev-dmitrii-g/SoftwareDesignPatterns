package internal;

public class MotorOverride extends NeuralInterface {
    public MotorOverride(ImplantFirmware firmware) {
        super(firmware);
    }

    @Override
    public void triggerImplant() {
        System.out.println("\\...Forcing-Motor-Override...//");
        NeuralPayload payload = new NeuralPayload(NeuralPayload.TargetRegion.MOTOR_CORTEX, 10, "AUTH_ENC_RUNNING_ENHANCEMENT");
        firmware.executeNeuralCommand(payload);
    }
}
