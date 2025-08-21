import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AddAppForm from './AddAppForm';

jest.mock('../services/api');
jest.mock('../services/storage');

const mockAppsApi = {
  add: jest.fn(),
};

const mockStorageService = {
  addApp: jest.fn(),
};

require('../services/api').appsApi = mockAppsApi;
require('../services/storage').storageService = mockStorageService;

describe('AddAppForm', () => {
  const mockOnAppAdded = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders form elements', () => {
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    expect(screen.getByText('Add New App to Track')).toBeInTheDocument();
    expect(screen.getByLabelText('App ID:')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('e.g., 595068606')).toBeInTheDocument();
    expect(screen.getByText('Add App')).toBeInTheDocument();
    expect(screen.getByText(/Find the App ID in the App Store URL/)).toBeInTheDocument();
  });

  test('requires app ID input', async () => {
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    const appIdInput = screen.getByLabelText('App ID:');
    
    expect(appIdInput).toBeRequired();
    expect(appIdInput).toHaveValue('');
  });

  test('submits form with valid app ID', async () => {
    mockAppsApi.add.mockResolvedValue(undefined);
    
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    const appIdInput = screen.getByLabelText('App ID:');
    const submitButton = screen.getByText('Add App');
    
    await userEvent.type(appIdInput, '123456789');
    await userEvent.click(submitButton);
    
    await waitFor(() => {
      expect(mockAppsApi.add).toHaveBeenCalledWith('123456789');
      expect(mockStorageService.addApp).toHaveBeenCalledWith('123456789');
      expect(mockOnAppAdded).toHaveBeenCalled();
    });
    
    expect(appIdInput).toHaveValue('');
  });

  test('shows loading state during submission', async () => {
    mockAppsApi.add.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)));
    
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    const appIdInput = screen.getByLabelText('App ID:');
    const submitButton = screen.getByText('Add App');
    
    await userEvent.type(appIdInput, '123456789');
    await userEvent.click(submitButton);
    
    expect(screen.getByText('Adding...')).toBeInTheDocument();
    expect(submitButton).toBeDisabled();
    
    await waitFor(() => {
      expect(screen.getByText('Add App')).toBeInTheDocument();
      expect(submitButton).not.toBeDisabled();
    });
  });

  test('shows error message on API failure', async () => {
    mockAppsApi.add.mockRejectedValue(new Error('API Error'));
    
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    const appIdInput = screen.getByLabelText('App ID:');
    const submitButton = screen.getByText('Add App');
    
    await userEvent.type(appIdInput, '123456789');
    await userEvent.click(submitButton);
    
    await waitFor(() => {
      expect(screen.getByText('Failed to add app. Please check the app ID and try again.')).toBeInTheDocument();
    });
    
    expect(mockOnAppAdded).not.toHaveBeenCalled();
    expect(appIdInput).toHaveValue('123456789');
  });

  test('clears error message on new submission', async () => {
    mockAppsApi.add.mockRejectedValueOnce(new Error('API Error'));
    mockAppsApi.add.mockResolvedValueOnce(undefined);
    
    render(<AddAppForm onAppAdded={mockOnAppAdded} />);
    
    const appIdInput = screen.getByLabelText('App ID:');
    const submitButton = screen.getByText('Add App');
    
    await userEvent.type(appIdInput, '123456789');
    await userEvent.click(submitButton);
    
    await waitFor(() => {
      expect(screen.getByText('Failed to add app. Please check the app ID and try again.')).toBeInTheDocument();
    });
    
    await userEvent.clear(appIdInput);
    await userEvent.type(appIdInput, '987654321');
    await userEvent.click(submitButton);
    
    await waitFor(() => {
      expect(screen.queryByText('Failed to add app. Please check the app ID and try again.')).not.toBeInTheDocument();
    });
  });
});
