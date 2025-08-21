import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AppSelector from './AppSelector';
import { App } from '../types';

describe('AppSelector', () => {
  const mockApps: App[] = [
    { id: '123', name: 'Test App' },
    { id: '456', name: 'Another App' },
    { id: '789' },
  ];

  const mockOnAppChange = jest.fn();
  const mockOnRemoveApp = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders reviews header with app name', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="123"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    expect(screen.getByText('Reviews for Test App')).toBeInTheDocument();
  });

  test('renders fallback name when app has no name', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="789"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    expect(screen.getByText('Reviews for App 789')).toBeInTheDocument();
  });

  test('renders fallback name when app is not found', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="999"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    expect(screen.getByText('Reviews for App 999')).toBeInTheDocument();
  });

  test('renders remove button', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="123"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    const removeButton = screen.getByTitle('Remove this app');
    expect(removeButton).toBeInTheDocument();
    expect(removeButton).toHaveTextContent('×');
  });

  test('calls onRemoveApp when remove button is clicked', async () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="123"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    const removeButton = screen.getByTitle('Remove this app');
    await userEvent.click(removeButton);

    expect(mockOnRemoveApp).toHaveBeenCalledWith('123');
  });

  test('has proper heading structure', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="123"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    const heading = screen.getByRole('heading', { level: 2 });
    expect(heading).toBeInTheDocument();
    expect(heading).toHaveTextContent('Reviews for Test App');
  });

  test('remove button has correct accessibility attributes', () => {
    render(
      <AppSelector
        apps={mockApps}
        selectedAppId="123"
        onAppChange={mockOnAppChange}
        onRemoveApp={mockOnRemoveApp}
      />
    );

    const removeButton = screen.getByTitle('Remove this app');
    expect(removeButton).toHaveAttribute('title', 'Remove this app');
  });
});
