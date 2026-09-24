package internal;

public abstract class NeuralInterface {
    ImplantFirmware firmware;

    NeuralInterface(ImplantFirmware firmware) {
        this.firmware = firmware;
    }

    abstract void triggerImplant();
}
