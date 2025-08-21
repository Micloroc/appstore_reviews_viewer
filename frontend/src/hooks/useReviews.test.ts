import { renderHook, waitFor } from '@testing-library/react';
import useReviews from './useReviews';
import { Review } from '../types';

jest.mock('../services/api');

const mockReviewsApi = {
  getByAppId: jest.fn(),
};

require('../services/api').reviewsApi = mockReviewsApi;

describe('useReviews', () => {
  const mockReviews: Review[] = [
    {
      id: '1',
      appId: '123',
      author: 'John Doe',
      content: 'Great app!',
      score: 5,
      submittedAt: '2023-12-01T10:00:00Z',
    },
    {
      id: '2',
      appId: '123',
      author: 'Jane Smith',
      content: 'Good but could be better.',
      score: 4,
      submittedAt: '2023-12-02T15:30:00Z',
    },
  ];

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('initializes with empty state', () => {
    const { result } = renderHook(() => useReviews(''));

    expect(result.current.reviews).toEqual([]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBe('');
  });

  test('does not fetch reviews when appId is empty', () => {
    renderHook(() => useReviews(''));

    expect(mockReviewsApi.getByAppId).not.toHaveBeenCalled();
  });

  test('fetches reviews when appId is provided', async () => {
    mockReviewsApi.getByAppId.mockResolvedValue({ reviews: mockReviews });

    const { result } = renderHook(() => useReviews('123'));

    expect(result.current.isLoading).toBe(true);
    expect(mockReviewsApi.getByAppId).toHaveBeenCalledWith('123');

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.reviews).toEqual(mockReviews);
    expect(result.current.error).toBe('');
  });

  test('sorts reviews by submission date (newest first)', async () => {
    const unsortedReviews = [
      {
        id: '1',
        appId: '123',
        author: 'User 1',
        content: 'Review 1',
        score: 5,
        submittedAt: '2023-12-01T10:00:00Z',
      },
      {
        id: '2',
        appId: '123',
        author: 'User 2',
        content: 'Review 2',
        score: 4,
        submittedAt: '2023-12-03T10:00:00Z',
      },
      {
        id: '3',
        appId: '123',
        author: 'User 3',
        content: 'Review 3',
        score: 3,
        submittedAt: '2023-12-02T10:00:00Z',
      },
    ];

    mockReviewsApi.getByAppId.mockResolvedValue({ reviews: unsortedReviews });

    const { result } = renderHook(() => useReviews('123'));

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.reviews[0].id).toBe('2');
    expect(result.current.reviews[1].id).toBe('3');
    expect(result.current.reviews[2].id).toBe('1');
  });

  test('handles API errors', async () => {
    mockReviewsApi.getByAppId.mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useReviews('123'));

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.reviews).toEqual([]);
    expect(result.current.error).toBe('Failed to fetch reviews. Please try again.');
  });

  test('clears reviews when appId changes to empty', async () => {
    mockReviewsApi.getByAppId.mockResolvedValue({ reviews: mockReviews });

    const { result, rerender } = renderHook(
      ({ appId }) => useReviews(appId),
      { initialProps: { appId: '123' } }
    );

    await waitFor(() => {
      expect(result.current.reviews).toEqual(mockReviews);
    });

    rerender({ appId: '' });

    expect(result.current.reviews).toEqual([]);
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBe('');
  });

  test('refetches reviews when appId changes', async () => {
    const app123Reviews = [mockReviews[0]];
    const app456Reviews = [mockReviews[1]];

    mockReviewsApi.getByAppId.mockResolvedValueOnce({ reviews: app123Reviews });
    mockReviewsApi.getByAppId.mockResolvedValueOnce({ reviews: app456Reviews });

    const { result, rerender } = renderHook(
      ({ appId }) => useReviews(appId),
      { initialProps: { appId: '123' } }
    );

    await waitFor(() => {
      expect(result.current.reviews).toEqual(app123Reviews);
    });

    rerender({ appId: '456' });

    expect(result.current.isLoading).toBe(true);

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.reviews).toEqual(app456Reviews);
    expect(mockReviewsApi.getByAppId).toHaveBeenCalledTimes(2);
    expect(mockReviewsApi.getByAppId).toHaveBeenNthCalledWith(1, '123');
    expect(mockReviewsApi.getByAppId).toHaveBeenNthCalledWith(2, '456');
  });

  test('clears error on successful refetch', async () => {
    mockReviewsApi.getByAppId.mockRejectedValueOnce(new Error('Network error'));
    mockReviewsApi.getByAppId.mockResolvedValueOnce({ reviews: mockReviews });

    const { result, rerender } = renderHook(
      ({ appId }) => useReviews(appId),
      { initialProps: { appId: '123' } }
    );

    await waitFor(() => {
      expect(result.current.error).toBe('Failed to fetch reviews. Please try again.');
    });

    rerender({ appId: '456' });

    await waitFor(() => {
      expect(result.current.reviews).toEqual(mockReviews);
      expect(result.current.error).toBe('');
    });
  });

  test('maintains loading state during fetch', async () => {
    let resolvePromise: (value: any) => void;
    const promise = new Promise((resolve) => {
      resolvePromise = resolve;
    });

    mockReviewsApi.getByAppId.mockReturnValue(promise);

    const { result } = renderHook(() => useReviews('123'));

    expect(result.current.isLoading).toBe(true);
    expect(result.current.reviews).toEqual([]);
    expect(result.current.error).toBe('');

    resolvePromise!({ reviews: mockReviews });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });
  });
});
