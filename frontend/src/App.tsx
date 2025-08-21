import React from 'react';
import AddAppForm from './components/AddAppForm';
import AppHeader from './components/AppHeader';
import AppControls from './components/AppControls';
import AppSelector from './components/AppSelector';
import ReviewsList from './components/ReviewsList';
import NoAppsMessage from './components/NoAppsMessage';
import { useApps, useReviews } from './hooks';
import './App.css';

function App() {
  const { apps, selectedAppId, refreshApps, selectApp, removeApp } = useApps();
  const { reviews, isLoading, error } = useReviews(selectedAppId);

  const handleAppAdded = () => {
    const updatedApps = refreshApps();
    
    if (updatedApps.length > 0) {
      const newApp = updatedApps[updatedApps.length - 1];
      selectApp(newApp.id);
    }
  };

  return (
    <div className="App">
      <AppHeader />

      <main className="App-main">
        <AddAppForm onAppAdded={handleAppAdded} />

        {apps.length > 0 && (
          <AppControls 
            apps={apps} 
            selectedAppId={selectedAppId} 
            onAppChange={selectApp} 
          />
        )}

        {selectedAppId && (
          <div className="reviews-section">
            <AppSelector 
              apps={apps} 
              selectedAppId={selectedAppId} 
              onAppChange={selectApp} 
              onRemoveApp={removeApp} 
            />
            <ReviewsList 
              reviews={reviews} 
              isLoading={isLoading} 
              error={error} 
            />
          </div>
        )}

        {apps.length === 0 && <NoAppsMessage />}
      </main>
    </div>
  );
}

export default App;
