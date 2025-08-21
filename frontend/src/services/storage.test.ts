import { storageService } from './storage';
import { App } from '../types';

describe('storageService', () => {
  let mockStorage: Record<string, string> = {};

  beforeEach(() => {
    mockStorage = {};
    
    const mockLocalStorage = {
      getItem: jest.fn((key: string) => mockStorage[key] || null),
      setItem: jest.fn((key: string, value: string) => {
        mockStorage[key] = value;
      }),
      removeItem: jest.fn((key: string) => {
        delete mockStorage[key];
      }),
      clear: jest.fn(() => {
        mockStorage = {};
      }),
    };

    Object.defineProperty(window, 'localStorage', {
      value: mockLocalStorage,
      writable: true,
    });
  });

  describe('getApps', () => {
    test('returns empty array when no apps stored', () => {
      const apps = storageService.getApps();
      expect(apps).toEqual([]);
    });

    test('returns stored apps', () => {
      const mockApps: App[] = [
        { id: '123', name: 'Test App' },
        { id: '456', name: 'Another App' },
      ];
      
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(mockApps);
      
      const apps = storageService.getApps();
      expect(apps).toEqual(mockApps);
    });

    test('returns empty array when storage contains invalid JSON', () => {
      mockStorage['appstore_reviews_viewer_apps'] = 'invalid json';
      
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
      const apps = storageService.getApps();
      
      expect(apps).toEqual([]);
      expect(consoleSpy).toHaveBeenCalledWith('Error reading from localStorage:', expect.any(SyntaxError));
      
      consoleSpy.mockRestore();
    });

    test('handles localStorage errors gracefully', () => {
      const mockLocalStorage = window.localStorage as any;
      mockLocalStorage.getItem.mockImplementation(() => {
        throw new Error('Storage error');
      });
      
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
      const apps = storageService.getApps();
      
      expect(apps).toEqual([]);
      expect(consoleSpy).toHaveBeenCalledWith('Error reading from localStorage:', expect.any(Error));
      
      consoleSpy.mockRestore();
    });
  });

  describe('saveApps', () => {
    test('saves apps to localStorage', () => {
      const mockApps: App[] = [
        { id: '123', name: 'Test App' },
        { id: '456', name: 'Another App' },
      ];
      
      storageService.saveApps(mockApps);
      
      expect(mockStorage['appstore_reviews_viewer_apps']).toBe(JSON.stringify(mockApps));
    });

    test('handles localStorage errors gracefully', () => {
      const mockLocalStorage = window.localStorage as any;
      mockLocalStorage.setItem.mockImplementation(() => {
        throw new Error('Storage error');
      });
      
      const consoleSpy = jest.spyOn(console, 'error').mockImplementation();
      
      storageService.saveApps([]);
      
      expect(consoleSpy).toHaveBeenCalledWith('Error writing to localStorage:', expect.any(Error));
      
      consoleSpy.mockRestore();
    });
  });

  describe('addApp', () => {
    test('adds new app to empty storage', () => {
      storageService.addApp('123');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual([{ id: '123' }]);
    });

    test('adds new app to existing apps', () => {
      const existingApps: App[] = [{ id: '123', name: 'Existing App' }];
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(existingApps);
      
      storageService.addApp('456');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual([
        { id: '123', name: 'Existing App' },
        { id: '456' },
      ]);
    });

    test('does not add duplicate app', () => {
      const existingApps: App[] = [{ id: '123', name: 'Existing App' }];
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(existingApps);
      
      storageService.addApp('123');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual(existingApps);
    });
  });

  describe('removeApp', () => {
    test('removes app from storage', () => {
      const existingApps: App[] = [
        { id: '123', name: 'App 1' },
        { id: '456', name: 'App 2' },
        { id: '789', name: 'App 3' },
      ];
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(existingApps);
      
      storageService.removeApp('456');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual([
        { id: '123', name: 'App 1' },
        { id: '789', name: 'App 3' },
      ]);
    });

    test('handles removing non-existent app', () => {
      const existingApps: App[] = [{ id: '123', name: 'App 1' }];
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(existingApps);
      
      storageService.removeApp('999');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual(existingApps);
    });

    test('handles removing from empty storage', () => {
      storageService.removeApp('123');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps'] || '[]');
      expect(stored).toEqual([]);
    });

    test('removes all instances of app with same ID', () => {
      const existingApps: App[] = [
        { id: '123', name: 'App 1' },
        { id: '456', name: 'App 2' },
        { id: '123', name: 'Duplicate App' },
      ];
      mockStorage['appstore_reviews_viewer_apps'] = JSON.stringify(existingApps);
      
      storageService.removeApp('123');
      
      const stored = JSON.parse(mockStorage['appstore_reviews_viewer_apps']);
      expect(stored).toEqual([{ id: '456', name: 'App 2' }]);
    });
  });

  describe('integration tests', () => {
    test('complete workflow: add, get, remove', () => {
      expect(storageService.getApps()).toEqual([]);
      
      storageService.addApp('123');
      expect(storageService.getApps()).toEqual([{ id: '123' }]);
      
      storageService.addApp('456');
      expect(storageService.getApps()).toEqual([
        { id: '123' },
        { id: '456' },
      ]);
      
      storageService.removeApp('123');
      expect(storageService.getApps()).toEqual([{ id: '456' }]);
      
      storageService.removeApp('456');
      expect(storageService.getApps()).toEqual([]);
    });
  });
});
