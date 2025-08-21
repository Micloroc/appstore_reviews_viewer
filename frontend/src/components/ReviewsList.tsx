import React from 'react';
import { Review } from '../types';
import ReviewCard from './ReviewCard';

interface ReviewsListProps {
  reviews: Review[];
  isLoading: boolean;
  error: string;
}

const ReviewsList: React.FC<ReviewsListProps> = ({ 
  reviews, 
  isLoading, 
  error 
}) => {
  if (error) {
    return <div className="error-message">{error}</div>;
  }

  if (isLoading) {
    return <div className="loading">Loading reviews...</div>;
  }

  if (reviews.length === 0) {
    return <div className="no-reviews">No reviews found for this app.</div>;
  }

  return (
    <div className="reviews-list">
      <div className="reviews-count">
        {reviews.length} review{reviews.length !== 1 ? 's' : ''}
      </div>
      {reviews.map((review) => (
        <ReviewCard key={review.id} review={review} />
      ))}
    </div>
  );
};

export default ReviewsList;
