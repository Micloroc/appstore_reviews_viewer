import { useState, useEffect } from 'react';
import { Review } from '../types';
import { reviewsApi } from '../services/api';

const useReviews = (appId: string) => {
  const [reviews, setReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchReviews = async (targetAppId: string) => {
    if (!targetAppId) {
      setReviews([]);
      return;
    }

    setIsLoading(true);
    setError('');
    
    try {
      const response = await reviewsApi.getByAppId(targetAppId);
      const sortedReviews = response.reviews.sort((a, b) => 
        new Date(b.submittedAt).getTime() - new Date(a.submittedAt).getTime()
      );
      setReviews(sortedReviews);
    } catch (err) {
      setError('Failed to fetch reviews. Please try again.');
      setReviews([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchReviews(appId);
  }, [appId]);

  return { reviews, isLoading, error };
};

export default useReviews;
