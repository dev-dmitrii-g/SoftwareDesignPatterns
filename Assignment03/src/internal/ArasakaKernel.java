package internal;

public class ArasakaKernel implements ImplantFirmware {
    @Override
    public void executeNeuralCommand(NeuralPayload payload) throws RejectionSyndromeException {
        if (!payload.getCommandText().startsWith("AUTH_ENC_")) {
            throw new RejectionSyndromeException("ArasakaKernel Safety Override: Unauthorized command.");
        }
        System.out.println("ArasakaKernel executing " + payload.getCommandText());
    }
}
