package internal.asbtraction;

import internal.implementor.ImplantFirmware;

public abstract class NeuralInterface {
    ImplantFirmware firmware;

    NeuralInterface(ImplantFirmware firmware) {
        this.firmware = firmware;
    }

    public abstract void triggerImplant();
}
