import { reviewsApi, appsApi } from './api';
import { Review } from '../types';

global.fetch = jest.fn();
const mockFetch = fetch as jest.MockedFunction<typeof fetch>;

describe('API Services', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('reviewsApi', () => {
    describe('getByAppId', () => {
      test('fetches reviews successfully', async () => {
        const mockReviews: Review[] = [
          {
            id: '1',
            appId: '123',
            author: 'Test User',
            content: 'Great app!',
            score: 5,
            submittedAt: '2023-12-01T10:00:00Z',
          },
        ];

        const mockResponse = { reviews: mockReviews };

        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          text: () => Promise.resolve(JSON.stringify(mockResponse)),
        } as Response);

        const result = await reviewsApi.getByAppId('123');

        expect(mockFetch).toHaveBeenCalledWith(
          'http://localhost:8080/api/v1/app/123/reviews/recent',
          {
            headers: {
              'Content-Type': 'application/json',
            },
          }
        );
        expect(result).toEqual(mockResponse);
      });

      test('throws error on HTTP error response', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: false,
          status: 404,
        } as Response);

        await expect(reviewsApi.getByAppId('123')).rejects.toThrow('HTTP error! status: 404');
      });

      test('handles empty response', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 204,
          headers: new Headers(),
        } as Response);

        const result = await reviewsApi.getByAppId('123');

        expect(result).toEqual({});
      });

      test('handles response with zero content length', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-length': '0' }),
        } as Response);

        const result = await reviewsApi.getByAppId('123');

        expect(result).toEqual({});
      });

      test('handles invalid JSON response gracefully', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          text: () => Promise.resolve('invalid json'),
        } as Response);

        const consoleSpy = jest.spyOn(console, 'warn').mockImplementation();

        const result = await reviewsApi.getByAppId('123');

        expect(result).toEqual({});
        expect(consoleSpy).toHaveBeenCalledWith('Failed to parse JSON response:', expect.any(SyntaxError));

        consoleSpy.mockRestore();
      });

      test('handles empty text response', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          text: () => Promise.resolve(''),
        } as Response);

        const result = await reviewsApi.getByAppId('123');

        expect(result).toEqual({});
      });

      test('uses custom API base URL from environment', async () => {
        const originalModule = require('./api');
        jest.resetModules();
        
        process.env.REACT_APP_API_URL = 'https://custom-api.com';
        
        const { reviewsApi: customReviewsApi } = require('./api');

        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers({ 'content-type': 'application/json' }),
          text: () => Promise.resolve('{}'),
        } as Response);

        await customReviewsApi.getByAppId('123');

        expect(mockFetch).toHaveBeenCalledWith(
          'https://custom-api.com/api/v1/app/123/reviews/recent',
          expect.any(Object)
        );

        delete process.env.REACT_APP_API_URL;
        jest.resetModules();
      });
    });
  });

  describe('appsApi', () => {
    describe('add', () => {
      test('adds app successfully', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 201,
          headers: new Headers(),
        } as Response);

        await appsApi.add('123456789');

        expect(mockFetch).toHaveBeenCalledWith(
          'http://localhost:8080/api/v1/app',
          {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ appId: '123456789' }),
          }
        );
      });

      test('throws error on HTTP error response', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: false,
          status: 400,
        } as Response);

        await expect(appsApi.add('123456789')).rejects.toThrow('HTTP error! status: 400');
      });

      test('handles successful response with no content', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 204,
          headers: new Headers(),
        } as Response);

        const result = await appsApi.add('123456789');

        expect(result).toEqual({});
      });

      test('sends correct request body format', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers(),
        } as Response);

        await appsApi.add('987654321');

        expect(mockFetch).toHaveBeenCalledWith(
          expect.any(String),
          expect.objectContaining({
            body: '{"appId":"987654321"}',
          })
        );
      });

      test('uses POST method', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers(),
        } as Response);

        await appsApi.add('123456789');

        expect(mockFetch).toHaveBeenCalledWith(
          expect.any(String),
          expect.objectContaining({
            method: 'POST',
          })
        );
      });

      test('includes correct headers', async () => {
        mockFetch.mockResolvedValueOnce({
          ok: true,
          status: 200,
          headers: new Headers(),
        } as Response);

        await appsApi.add('123456789');

        expect(mockFetch).toHaveBeenCalledWith(
          expect.any(String),
          expect.objectContaining({
            headers: {
              'Content-Type': 'application/json',
            },
          })
        );
      });
    });
  });

  describe('error handling', () => {
    test('handles network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(reviewsApi.getByAppId('123')).rejects.toThrow('Network error');
    });

    test('handles fetch abort errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('AbortError'));

      await expect(appsApi.add('123')).rejects.toThrow('AbortError');
    });
  });
});
