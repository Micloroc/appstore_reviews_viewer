import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import App from './App';

jest.mock('./services/api');
jest.mock('./services/storage');

const mockStorageService = {
  getApps: jest.fn(),
  addApp: jest.fn(),
  removeApp: jest.fn(),
  saveApps: jest.fn(),
};

const mockAppsApi = {
  add: jest.fn(),
};

const mockReviewsApi = {
  getByAppId: jest.fn(),
};

require('./services/storage').storageService = mockStorageService;
require('./services/api').appsApi = mockAppsApi;
require('./services/api').reviewsApi = mockReviewsApi;

describe('App', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockStorageService.getApps.mockReturnValue([]);
  });

  test('renders app header', () => {
    render(<App />);
    expect(screen.getByText('App Store Reviews Viewer')).toBeInTheDocument();
  });

  test('shows no apps message when no apps are added', () => {
    render(<App />);
    expect(screen.getByText('No apps added yet. Add your first app above to get started!')).toBeInTheDocument();
  });

  test('renders add app form', () => {
    render(<App />);
    expect(screen.getByText('Add New App to Track')).toBeInTheDocument();
    expect(screen.getByLabelText('App ID:')).toBeInTheDocument();
    expect(screen.getByText('Add App')).toBeInTheDocument();
  });

  test('shows app controls when apps are available', () => {
    mockStorageService.getApps.mockReturnValue([
      { id: '12345', name: 'Test App' }
    ]);

    render(<App />);
    expect(screen.getByText('Select App:')).toBeInTheDocument();
    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('12345');
  });

  test('shows reviews section when app is selected', async () => {
    mockStorageService.getApps.mockReturnValue([
      { id: '12345', name: 'Test App' }
    ]);

    mockReviewsApi.getByAppId.mockResolvedValue({
      reviews: [
        {
          id: '1',
          appId: '12345',
          author: 'Test User',
          content: 'Great app!',
          score: 5,
          submittedAt: '2023-01-01T00:00:00Z'
        }
      ]
    });

    render(<App />);
    
    await waitFor(() => {
      expect(screen.getByText('Reviews for Test App')).toBeInTheDocument();
    });
  });
});
