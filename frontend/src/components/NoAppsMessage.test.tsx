import React from 'react';
import { render, screen } from '@testing-library/react';
import NoAppsMessage from './NoAppsMessage';

describe('NoAppsMessage', () => {
  test('renders no apps message', () => {
    render(<NoAppsMessage />);
    
    expect(screen.getByText('No apps added yet. Add your first app above to get started!')).toBeInTheDocument();
  });

  test('renders message in a paragraph element', () => {
    render(<NoAppsMessage />);
    
    const message = screen.getByText('No apps added yet. Add your first app above to get started!');
    expect(message.tagName).toBe('P');
  });

  test('has correct container structure', () => {
    const { container } = render(<NoAppsMessage />);
    
    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('no-apps-message');
  });
});
