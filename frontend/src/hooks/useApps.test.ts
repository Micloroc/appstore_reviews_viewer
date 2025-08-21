import { renderHook, act } from '@testing-library/react';
import useApps from './useApps';
import { App } from '../types';

jest.mock('../services/storage');

const mockStorageService = {
  getApps: jest.fn(),
  addApp: jest.fn(),
  removeApp: jest.fn(),
  saveApps: jest.fn(),
};

require('../services/storage').storageService = mockStorageService;

describe('useApps', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('initializes with empty apps when storage is empty', () => {
    mockStorageService.getApps.mockReturnValue([]);

    const { result } = renderHook(() => useApps());

    expect(result.current.apps).toEqual([]);
    expect(result.current.selectedAppId).toBe('');
  });

  test('initializes with apps from storage', () => {
    const mockApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'Another App' },
    ];
    mockStorageService.getApps.mockReturnValue(mockApps);

    const { result } = renderHook(() => useApps());

    expect(result.current.apps).toEqual(mockApps);
    expect(result.current.selectedAppId).toBe('123');
  });

  test('does not auto-select when no apps available', () => {
    mockStorageService.getApps.mockReturnValue([]);

    const { result } = renderHook(() => useApps());

    expect(result.current.selectedAppId).toBe('');
  });

  test('refreshApps updates apps from storage', () => {
    const initialApps: App[] = [{ id: '123', name: 'Test App' }];
    const updatedApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'New App' },
    ];

    mockStorageService.getApps.mockReturnValueOnce(initialApps);
    mockStorageService.getApps.mockReturnValueOnce(updatedApps);

    const { result } = renderHook(() => useApps());

    expect(result.current.apps).toEqual(initialApps);

    act(() => {
      const returnedApps = result.current.refreshApps();
      expect(returnedApps).toEqual(updatedApps);
    });

    expect(result.current.apps).toEqual(updatedApps);
  });

  test('selectApp changes selected app ID', () => {
    const mockApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'Another App' },
    ];
    mockStorageService.getApps.mockReturnValue(mockApps);

    const { result } = renderHook(() => useApps());

    expect(result.current.selectedAppId).toBe('123');

    act(() => {
      result.current.selectApp('456');
    });

    expect(result.current.selectedAppId).toBe('456');
  });

  test('removeApp calls storage service and refreshes apps', () => {
    const initialApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'Another App' },
    ];
    const appsAfterRemoval: App[] = [{ id: '456', name: 'Another App' }];

    mockStorageService.getApps.mockReturnValueOnce(initialApps);
    mockStorageService.getApps.mockReturnValueOnce(appsAfterRemoval);

    const { result } = renderHook(() => useApps());

    expect(result.current.selectedAppId).toBe('123');

    act(() => {
      result.current.removeApp('123');
    });

    expect(mockStorageService.removeApp).toHaveBeenCalledWith('123');
    expect(result.current.apps).toEqual(appsAfterRemoval);
    expect(result.current.selectedAppId).toBe('456');
  });

  test('removeApp clears selection when removing selected app and no apps remain', () => {
    const initialApps: App[] = [{ id: '123', name: 'Test App' }];
    const appsAfterRemoval: App[] = [];

    mockStorageService.getApps.mockReturnValueOnce(initialApps);
    mockStorageService.getApps.mockReturnValueOnce(appsAfterRemoval);

    const { result } = renderHook(() => useApps());

    expect(result.current.selectedAppId).toBe('123');

    act(() => {
      result.current.removeApp('123');
    });

    expect(result.current.selectedAppId).toBe('');
  });

  test('removeApp does not change selection when removing non-selected app', () => {
    const initialApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'Another App' },
    ];
    const appsAfterRemoval: App[] = [{ id: '123', name: 'Test App' }];

    mockStorageService.getApps.mockReturnValueOnce(initialApps);
    mockStorageService.getApps.mockReturnValueOnce(appsAfterRemoval);

    const { result } = renderHook(() => useApps());

    expect(result.current.selectedAppId).toBe('123');

    act(() => {
      result.current.removeApp('456');
    });

    expect(result.current.selectedAppId).toBe('123');
  });

  test('hook does not auto-select when apps are already selected', () => {
    const mockApps: App[] = [
      { id: '123', name: 'Test App' },
      { id: '456', name: 'Another App' },
    ];

    mockStorageService.getApps.mockReturnValue(mockApps);

    const { result } = renderHook(() => useApps());

    act(() => {
      result.current.selectApp('456');
    });

    expect(result.current.selectedAppId).toBe('456');

    const { result: result2 } = renderHook(() => useApps());
    expect(result2.current.selectedAppId).toBe('123');
  });
});
