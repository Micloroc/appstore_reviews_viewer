import React from 'react';
import { render, screen } from '@testing-library/react';
import ReviewCard from './ReviewCard';
import { Review } from '../types';

describe('ReviewCard', () => {
  const mockReview: Review = {
    id: '1',
    appId: '123456789',
    author: 'John Doe',
    content: 'This is a great app! I love using it every day.',
    score: 4,
    submittedAt: '2023-12-01T10:30:00Z',
  };

  test('renders review information', () => {
    render(<ReviewCard review={mockReview} />);
    
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('This is a great app! I love using it every day.')).toBeInTheDocument();
    expect(screen.getByText('(4/5)')).toBeInTheDocument();
  });

  test('renders star rating correctly', () => {
    render(<ReviewCard review={mockReview} />);
    
    const starsElement = screen.getByText('★★★★☆');
    expect(starsElement).toBeInTheDocument();
  });

  test('renders 5-star rating correctly', () => {
    const fiveStarReview = { ...mockReview, score: 5 };
    render(<ReviewCard review={fiveStarReview} />);
    
    const starsElement = screen.getByText('★★★★★');
    expect(starsElement).toBeInTheDocument();
    expect(screen.getByText('(5/5)')).toBeInTheDocument();
  });

  test('renders 1-star rating correctly', () => {
    const oneStarReview = { ...mockReview, score: 1 };
    render(<ReviewCard review={oneStarReview} />);
    
    const starsElement = screen.getByText('★☆☆☆☆');
    expect(starsElement).toBeInTheDocument();
    expect(screen.getByText('(1/5)')).toBeInTheDocument();
  });

  test('formats date correctly', () => {
    render(<ReviewCard review={mockReview} />);
    
    expect(screen.getByText(/Submitted: December 1, 2023/)).toBeInTheDocument();
  });

  test('handles different date formats', () => {
    const reviewWithDifferentDate = {
      ...mockReview,
      submittedAt: '2023-01-15T14:45:30Z',
    };
    
    render(<ReviewCard review={reviewWithDifferentDate} />);
    
    expect(screen.getByText(/Submitted: January 15, 2023/)).toBeInTheDocument();
  });

  test('renders review with empty content', () => {
    const reviewWithEmptyContent = { ...mockReview, content: '' };
    render(<ReviewCard review={reviewWithEmptyContent} />);
    
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    const contentElement = document.querySelector('.review-content');
    expect(contentElement).toBeInTheDocument();
    expect(contentElement).toHaveTextContent('');
  });

  test('renders review with long content', () => {
    const longContent = 'This is a very long review that goes on and on about how amazing this app is. '.repeat(5);
    const reviewWithLongContent = { ...mockReview, content: longContent.trim() };
    
    render(<ReviewCard review={reviewWithLongContent} />);
    
    expect(screen.getByText(longContent.trim())).toBeInTheDocument();
  });

  test('renders review with special characters in author name', () => {
    const reviewWithSpecialAuthor = { ...mockReview, author: 'José María O\'Connor' };
    render(<ReviewCard review={reviewWithSpecialAuthor} />);
    
    expect(screen.getByText('José María O\'Connor')).toBeInTheDocument();
  });
});
