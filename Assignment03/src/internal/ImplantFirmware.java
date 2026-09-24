package internal;

public interface ImplantFirmware {
    void executeNeuralCommand(NeuralPayload payload) throws RejectionSyndromeException;
}
