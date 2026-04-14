import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider } from 'react-router-dom';

import { LocaleProvider } from './app/providers/LocaleProvider';
import { SessionProvider } from './app/providers/SessionProvider';
import { router } from './app/router';
import './styles.css';

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element #root was not found');
}

createRoot(rootElement).render(
  <StrictMode>
    <SessionProvider>
      <LocaleProvider>
        <RouterProvider router={router} />
      </LocaleProvider>
    </SessionProvider>
  </StrictMode>,
);
