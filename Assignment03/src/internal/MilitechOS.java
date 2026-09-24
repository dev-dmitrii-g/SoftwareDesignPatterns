package internal;

public class MilitechOS implements ImplantFirmware {
    @Override
    public void executeNeuralCommand(NeuralPayload payload) {
        if (payload.getIntensity() > 8) {
            throw new RejectionSyndromeException("MilitechOS Safety Override: Intensity too high.");
        }
        System.out.println("MilitechOS routing " + payload.getCommandText() + " to " + payload.getTargetRegion());
    }
}
