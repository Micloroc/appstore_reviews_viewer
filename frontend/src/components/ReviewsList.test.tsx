import React from 'react';
import { render, screen } from '@testing-library/react';
import ReviewsList from './ReviewsList';
import { Review } from '../types';

describe('ReviewsList', () => {
  const mockReviews: Review[] = [
    {
      id: '1',
      appId: '123456789',
      author: 'John Doe',
      content: 'Great app!',
      score: 5,
      submittedAt: '2023-12-01T10:30:00Z',
    },
    {
      id: '2',
      appId: '123456789',
      author: 'Jane Smith',
      content: 'Could be better.',
      score: 3,
      submittedAt: '2023-12-02T15:45:00Z',
    },
  ];

  test('renders loading state', () => {
    render(<ReviewsList reviews={[]} isLoading={true} error="" />);
    
    expect(screen.getByText('Loading reviews...')).toBeInTheDocument();
  });

  test('renders error state', () => {
    const errorMessage = 'Failed to fetch reviews';
    render(<ReviewsList reviews={[]} isLoading={false} error={errorMessage} />);
    
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  test('renders no reviews message when list is empty', () => {
    render(<ReviewsList reviews={[]} isLoading={false} error="" />);
    
    expect(screen.getByText('No reviews found for this app.')).toBeInTheDocument();
  });

  test('renders reviews when available', () => {
    render(<ReviewsList reviews={mockReviews} isLoading={false} error="" />);
    
    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    expect(screen.getByText('Great app!')).toBeInTheDocument();
    expect(screen.getByText('Could be better.')).toBeInTheDocument();
  });

  test('renders correct review count for multiple reviews', () => {
    render(<ReviewsList reviews={mockReviews} isLoading={false} error="" />);
    
    expect(screen.getByText('2 reviews')).toBeInTheDocument();
  });

  test('renders correct review count for single review', () => {
    const singleReview = [mockReviews[0]];
    render(<ReviewsList reviews={singleReview} isLoading={false} error="" />);
    
    expect(screen.getByText('1 review')).toBeInTheDocument();
  });

  test('does not render reviews count when loading', () => {
    render(<ReviewsList reviews={mockReviews} isLoading={true} error="" />);
    
    expect(screen.queryByText('2 reviews')).not.toBeInTheDocument();
    expect(screen.getByText('Loading reviews...')).toBeInTheDocument();
  });

  test('does not render reviews count when error occurs', () => {
    render(<ReviewsList reviews={mockReviews} isLoading={false} error="Some error" />);
    
    expect(screen.queryByText('2 reviews')).not.toBeInTheDocument();
    expect(screen.getByText('Some error')).toBeInTheDocument();
  });

  test('prioritizes error over loading state', () => {
    render(<ReviewsList reviews={[]} isLoading={true} error="Network error" />);
    
    expect(screen.getByText('Network error')).toBeInTheDocument();
    expect(screen.queryByText('Loading reviews...')).not.toBeInTheDocument();
  });

  test('renders large number of reviews', () => {
    const manyReviews = Array.from({ length: 50 }, (_, index) => ({
      id: `${index + 1}`,
      appId: '123456789',
      author: `User ${index + 1}`,
      content: `Review content ${index + 1}`,
      score: (index % 5) + 1,
      submittedAt: new Date(2023, 0, index + 1).toISOString(),
    }));

    render(<ReviewsList reviews={manyReviews} isLoading={false} error="" />);
    
    expect(screen.getByText('50 reviews')).toBeInTheDocument();
    expect(screen.getByText('User 1')).toBeInTheDocument();
    expect(screen.getByText('User 50')).toBeInTheDocument();
  });
});
