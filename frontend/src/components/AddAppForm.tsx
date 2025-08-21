import React, { useState } from 'react';
import { appsApi } from '../services/api';
import { storageService } from '../services/storage';
import './AddAppForm.css';

interface AddAppFormProps {
  onAppAdded: () => void;
}

const AddAppForm: React.FC<AddAppFormProps> = ({ onAppAdded }) => {
  const [appId, setAppId] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      await appsApi.add(appId);
      
      storageService.addApp(appId);
      
      setAppId('');
      onAppAdded();
    } catch (err) {
      console.error(err);
      setError('Failed to add app. Please check the app ID and try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="add-app-form">
      <h3>Add New App to Track</h3>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="appId">App ID:</label>
          <input
            type="text"
            id="appId"
            value={appId}
            onChange={(e) => setAppId(e.target.value)}
            placeholder="e.g., 595068606"
            required
          />
          <small>
            Find the App ID in the App Store URL: https://apps.apple.com/us/app/appname/id[APP_ID]
          </small>
        </div>

        {error && <div className="error-message">{error}</div>}

        <button type="submit" disabled={isLoading}>
          {isLoading ? 'Adding...' : 'Add App'}
        </button>
      </form>
    </div>
  );
};

export default AddAppForm;
