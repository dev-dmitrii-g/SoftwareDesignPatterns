package internal.implementor;

import internal.exception.RejectionSyndromeException;

public class ArasakaKernel implements ImplantFirmware {
    @Override
    public void executeNeuralCommand(NeuralPayload payload) throws RejectionSyndromeException {
        if (!payload.getCommandText().startsWith("AUTH_ENC_")) {
            throw new RejectionSyndromeException("|-<A>-| ArasakaKernel Safety Override: Unauthorized command. |-<A>-|");
        }
        System.out.println("ArasakaKernel executing " + payload.getCommandText());
    }
}
