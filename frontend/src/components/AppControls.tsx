import React from 'react';
import { App as AppType } from '../types';

interface AppControlsProps {
  apps: AppType[];
  selectedAppId: string;
  onAppChange: (appId: string) => void;
}

const AppControls: React.FC<AppControlsProps> = ({ 
  apps, 
  selectedAppId, 
  onAppChange 
}) => (
  <div className="controls">
    <div className="control-group">
      <label htmlFor="app-select">Select App:</label>
      <select
        id="app-select"
        value={selectedAppId}
        onChange={(e) => onAppChange(e.target.value)}
      >
        <option value="">Choose an app...</option>
        {apps.map((app) => (
          <option key={app.id} value={app.id}>
            {app.name || `App ${app.id}`}
          </option>
        ))}
      </select>
    </div>
  </div>
);

export default AppControls;
