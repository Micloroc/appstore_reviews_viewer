import React from 'react';
import { App as AppType } from '../types';

interface AppSelectorProps {
  apps: AppType[];
  selectedAppId: string;
  onAppChange: (appId: string) => void;
  onRemoveApp: (appId: string) => void;
}

const AppSelector: React.FC<AppSelectorProps> = ({ 
  apps, 
  selectedAppId, 
  onAppChange, 
  onRemoveApp 
}) => {
  const selectedApp = apps.find(app => app.id === selectedAppId);

  return (
    <div className="reviews-header">
      <h2>
        Reviews for {selectedApp?.name || `App ${selectedAppId}`}
      </h2>
      <button 
        className="remove-app-btn"
        onClick={() => onRemoveApp(selectedAppId)}
        title="Remove this app"
      >
        ×
      </button>
    </div>
  );
};

export default AppSelector;
