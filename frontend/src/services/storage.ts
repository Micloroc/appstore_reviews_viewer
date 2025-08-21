import { App } from '../types';

const STORAGE_KEY = 'appstore_reviews_viewer_apps';

export const storageService = {
  getApps: (): App[] => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      return stored ? JSON.parse(stored) : [];
    } catch (error) {
      console.error('Error reading from localStorage:', error);
      return [];
    }
  },

  saveApps: (apps: App[]): void => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(apps));
    } catch (error) {
      console.error('Error writing to localStorage:', error);
    }
  },

  addApp: (appId: string): void => {
    const apps = storageService.getApps();
    const existingApp = apps.find(app => app.id === appId);
    
    if (!existingApp) {
      apps.push({ id: appId });
      storageService.saveApps(apps);
    }
  },

  removeApp: (appId: string): void => {
    const apps = storageService.getApps();
    const filteredApps = apps.filter(app => app.id !== appId);
    storageService.saveApps(filteredApps);
  },
};
