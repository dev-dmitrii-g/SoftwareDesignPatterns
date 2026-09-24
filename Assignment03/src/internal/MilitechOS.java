package internal;

public class MilitechOS implements ImplantFirmware {
    @Override
    public void executeNeuralCommand(NeuralPayload payload) {
        if (payload.getIntensity() > 8) {
            throw new RejectionSyndromeException("\\\\->->->-\\-MILITECHOS_SAFETY_OVERRIDE-//->->->-\\-INTENSITY_TOO_HIGH-//->->->-////");
        }
        System.out.println("MilitechOS routing " + payload.getCommandText() + " to " + payload.getTargetRegion());
    }
}
