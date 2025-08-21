import React from 'react';
import { Review } from '../types';
import './ReviewCard.css';

interface ReviewCardProps {
  review: Review;
}

const ReviewCard: React.FC<ReviewCardProps> = ({ review }) => {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const renderStars = (score: number) => {
    return '★'.repeat(score) + '☆'.repeat(5 - score);
  };

  return (
    <div className="review-card">
      <div className="review-header">
        <div className="review-author">{review.author}</div>
        <div className="review-score">
          <span className="stars">{renderStars(review.score)}</span>
          <span className="score-number">({review.score}/5)</span>
        </div>
      </div>
      <div className="review-content">{review.content}</div>
      <div className="review-date">
        Submitted: {formatDate(review.submittedAt)}
      </div>
    </div>
  );
};

export default ReviewCard;
