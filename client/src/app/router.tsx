import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AdminLayout } from './layouts/AdminLayout';
import { StoreLayout } from './layouts/StoreLayout';
import { LandingPage } from '../features/landing/LandingPage';
import { LoginPage } from '../features/auth/LoginPage';
import { ProductDetailPage } from '../features/catalog/ProductDetailPage';
import { SearchResultsPage } from '../features/search/SearchResultsPage';
import { AdminProjectListPage } from '../features/admin-projects/AdminProjectListPage';
import { AdminProjectFormPage } from '../features/admin-projects/AdminProjectFormPage';
import { AdminSiteSettingsPage } from '../features/admin-settings/AdminSiteSettingsPage';
import { AdminTechnologyListPage } from '../features/admin-technologies/AdminTechnologyListPage';
import { AdminTechnologyFormPage } from '../features/admin-technologies/AdminTechnologyFormPage';
import { RequireAdmin } from '../shared/routing/RequireAdmin';
import { NotFoundPage } from '../shared/ui/NotFoundPage';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <StoreLayout />,
    children: [
      { index: true, element: <LandingPage /> },
      { path: 'projects', element: <Navigate replace to="/" /> },
      { path: 'projects/:slug', element: <ProductDetailPage /> },
      { path: 'search', element: <SearchResultsPage /> },
      { path: 'login', element: <LoginPage /> },
    ],
  },
  {
    element: <RequireAdmin />,
    children: [
      {
        path: '/admin',
        element: <AdminLayout />,
        children: [
          { index: true, element: <Navigate replace to="/admin/projects" /> },
          { path: 'projects', element: <AdminProjectListPage /> },
          { path: 'projects/new', element: <AdminProjectFormPage /> },
          { path: 'projects/:id', element: <AdminProjectFormPage /> },
          { path: 'settings', element: <AdminSiteSettingsPage /> },
          { path: 'technologies', element: <AdminTechnologyListPage /> },
          { path: 'technologies/new', element: <AdminTechnologyFormPage /> },
          { path: 'technologies/:id', element: <AdminTechnologyFormPage /> },
        ],
      },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
  {
    path: '/home',
    element: <Navigate to="/" replace />,
  },
]);
