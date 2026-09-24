package internal;

public class SensoryAugmentation extends NeuralInterface {
    public SensoryAugmentation(ImplantFirmware firmware) {
        super(firmware);
    }
    @Override
    public void triggerImplant() {
        System.out.println("\\...Activating-Sensory-Augmentation...//");
        NeuralPayload payload = new NeuralPayload(NeuralPayload.TargetRegion.OPTIC_NERVE, 3, "AUTH_ENC_NIGHT_VISION");
        firmware.executeNeuralCommand(payload);
    }
}
