import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AdminLayout } from './layouts/AdminLayout';
import { StoreLayout } from './layouts/StoreLayout';
import { LandingPage } from '../features/landing/LandingPage';
import { LoginPage } from '../features/auth/LoginPage';
import { ProductDetailPage } from '../features/catalog/ProductDetailPage';
import { SearchResultsPage } from '../features/search/SearchResultsPage';
import { AdminProductListPage } from '../features/admin-products/AdminProductListPage';
import { AdminProductFormPage } from '../features/admin-products/AdminProductFormPage';
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
      { path: 'projects/:id', element: <ProductDetailPage /> },
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
          { path: 'projects', element: <AdminProductListPage /> },
          { path: 'projects/new', element: <AdminProductFormPage /> },
          { path: 'projects/:id', element: <AdminProductFormPage /> },
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
