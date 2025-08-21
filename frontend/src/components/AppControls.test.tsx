import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AppControls from './AppControls';
import { App } from '../types';

describe('AppControls', () => {
  const mockApps: App[] = [
    { id: '123', name: 'First App' },
    { id: '456', name: 'Second App' },
    { id: '789' },
  ];

  const mockOnAppChange = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders select label and dropdown', () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    expect(screen.getByLabelText('Select App:')).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument();
  });

  test('renders default option', () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    expect(screen.getByText('Choose an app...')).toBeInTheDocument();
  });

  test('renders all apps as options', () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    expect(screen.getByText('First App')).toBeInTheDocument();
    expect(screen.getByText('Second App')).toBeInTheDocument();
    expect(screen.getByText('App 789')).toBeInTheDocument();
  });

  test('shows fallback name for apps without name', () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    expect(screen.getByText('App 789')).toBeInTheDocument();
  });

  test('displays selected app', () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId="456"
        onAppChange={mockOnAppChange}
      />
    );

    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('456');
  });

  test('calls onAppChange when selection changes', async () => {
    render(
      <AppControls
        apps={mockApps}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    const select = screen.getByRole('combobox');
    await userEvent.selectOptions(select, '123');

    expect(mockOnAppChange).toHaveBeenCalledWith('123');
  });

  test('handles empty apps array', () => {
    render(
      <AppControls
        apps={[]}
        selectedAppId=""
        onAppChange={mockOnAppChange}
      />
    );

    expect(screen.getByText('Choose an app...')).toBeInTheDocument();
    expect(screen.queryByText('First App')).not.toBeInTheDocument();
  });

  test('maintains selection when apps list updates', () => {
    const { rerender } = render(
      <AppControls
        apps={mockApps}
        selectedAppId="456"
        onAppChange={mockOnAppChange}
      />
    );

    const updatedApps = [...mockApps, { id: '999', name: 'New App' }];
    
    rerender(
      <AppControls
        apps={updatedApps}
        selectedAppId="456"
        onAppChange={mockOnAppChange}
      />
    );

    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('456');
  });
});
