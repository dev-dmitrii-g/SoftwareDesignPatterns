package internal.implementor;

import internal.exception.RejectionSyndromeException;

public interface ImplantFirmware {
    void executeNeuralCommand(NeuralPayload payload) throws RejectionSyndromeException;
}
