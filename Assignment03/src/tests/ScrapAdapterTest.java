package tests;

import internal.BlackMarketAdapter;
import internal.BlackMarketWiring;
import internal.NeuralPayload;
import internal.RejectionSyndromeException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyFloat;
import static org.mockito.Mockito.*;

public class ScrapAdapterTest {
    @Test
    void executeNeuralCommand_translatesPayloadToRawData() {
        BlackMarketWiring mockHardware = mock(BlackMarketWiring.class);
        when(mockHardware.shockSystem(anyFloat(), any(byte[].class))).thenReturn(0);

        BlackMarketAdapter adapter = new BlackMarketAdapter(mockHardware);
        NeuralPayload payload = new NeuralPayload(NeuralPayload.TargetRegion.OPTIC_NERVE, 4, "TEST");
        adapter.executeNeuralCommand(payload);

        verify(mockHardware, times(1)).shockSystem(eq(9.48f), any(byte[].class));
    }

    @Test
    void executeNeuralCommand_translatesLegacyErrorCodeToException() {
        BlackMarketWiring mockHardware = mock(BlackMarketWiring.class);
        when(mockHardware.shockSystem(anyFloat(), any(byte[].class))).thenReturn(137);

        BlackMarketAdapter adapter = new BlackMarketAdapter(mockHardware);
        NeuralPayload payload = new NeuralPayload(NeuralPayload.TargetRegion.MOTOR_CORTEX, 10, "SPEED");

        RejectionSyndromeException thrown = assertThrows(RejectionSyndromeException.class, () -> adapter.executeNeuralCommand(payload));
        assertTrue(thrown.getMessage().contains("> Neural Rejection Error!"));
    }
}
