import React from 'react';
import { render, screen } from '@testing-library/react';
import AppHeader from './AppHeader';

describe('AppHeader', () => {
  test('renders app title', () => {
    render(<AppHeader />);
    
    expect(screen.getByText('App Store Reviews Viewer')).toBeInTheDocument();
  });

  test('renders as header element', () => {
    render(<AppHeader />);
    
    const header = screen.getByRole('banner');
    expect(header).toBeInTheDocument();
  });

  test('has correct heading level', () => {
    render(<AppHeader />);
    
    const heading = screen.getByRole('heading', { level: 1 });
    expect(heading).toBeInTheDocument();
    expect(heading).toHaveTextContent('App Store Reviews Viewer');
  });
});
