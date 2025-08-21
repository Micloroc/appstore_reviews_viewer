import { useState, useEffect } from 'react';
import { App as AppType } from '../types';
import { storageService } from '../services/storage';

const useApps = () => {
  const [apps, setApps] = useState<AppType[]>([]);
  const [selectedAppId, setSelectedAppId] = useState<string>('');

  useEffect(() => {
    const storedApps = storageService.getApps();
    setApps(storedApps);
    
    if (storedApps.length > 0) {
      setSelectedAppId(currentSelectedAppId => 
        currentSelectedAppId || storedApps[0].id
      );
    }
  }, []);

  const refreshApps = () => {
    const updatedApps = storageService.getApps();
    setApps(updatedApps);
    return updatedApps;
  };

  const selectApp = (appId: string) => {
    setSelectedAppId(appId);
  };

  const removeApp = (appId: string) => {
    storageService.removeApp(appId);
    const updatedApps = refreshApps();
    
    if (selectedAppId === appId) {
      setSelectedAppId(updatedApps.length > 0 ? updatedApps[0].id : '');
    }
  };

  return {
    apps,
    selectedAppId,
    refreshApps,
    selectApp,
    removeApp
  };
};

export default useApps;
