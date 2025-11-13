// React example with Security Platform Frontend SDK
import React, { useEffect } from 'react';
import { initSecuritySDK } from '@security-platform/frontend-sdk';

function App() {
  useEffect(() => {
    // Initialize SDK on app mount
    const sdk = initSecuritySDK({
      controlPlaneUrl: process.env.REACT_APP_CONTROL_PLANE_URL || 'https://console.securityplatform.com',
      authKey: process.env.REACT_APP_AUTH_KEY || '',
      serviceName: 'react-app',
      environment: process.env.NODE_ENV,
      redactionRules: ['password', 'card', 'ssn'],
      enableSessionReplay: true,
      enableTelemetry: true,
    });

    // Cleanup on unmount
    return () => {
      sdk.destroy();
    };
  }, []);

  const handleButtonClick = () => {
    // SDK automatically captures events, but you can also manually capture
    // Custom events are captured via the SDK instance
  };

  return (
    <div className="App">
      <button onClick={handleButtonClick}>Click Me</button>
    </div>
  );
}

export default App;
