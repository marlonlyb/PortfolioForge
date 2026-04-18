import { Navigate, createBrowserRouter } from 'react-router-dom';

import { AdminLayout } from './layouts/AdminLayout';
import { StoreLayout } from './layouts/StoreLayout';
import { LandingPage } from '../features/landing/LandingPage';
import { CompleteProfilePage } from '../features/auth/CompleteProfilePage';
import { LoginPage } from '../features/auth/LoginPage';
import { VerifyEmailOtpPage } from '../features/auth/VerifyEmailOtpPage';
import { ProductDetailPage } from '../features/catalog/ProductDetailPage';
import { SearchResultsPage } from '../features/search/SearchResultsPage';
import { AdminProjectListPage } from '../features/admin-projects/AdminProjectListPage';
import { AdminProjectFormPage } from '../features/admin-projects/AdminProjectFormPage';
import { AdminCaseStudyWorkflowPage } from '../features/admin-settings/AdminCaseStudyWorkflowPage';
import { AdminSiteSettingsPage } from '../features/admin-settings/AdminSiteSettingsPage';
import { AdminTechnologyListPage } from '../features/admin-technologies/AdminTechnologyListPage';
import { AdminTechnologyFormPage } from '../features/admin-technologies/AdminTechnologyFormPage';
import { AdminUserListPage } from '../features/admin-users/AdminUserListPage';
import { AdminUserFormPage } from '../features/admin-users/AdminUserFormPage';
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
      { path: 'signup', element: <LoginPage variant="signup" /> },
      { path: 'admin/login', element: <LoginPage variant="admin" /> },
      { path: 'verify-email', element: <VerifyEmailOtpPage /> },
      { path: 'complete-profile', element: <CompleteProfilePage /> },
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
          { path: 'settings/case-studies', element: <AdminCaseStudyWorkflowPage /> },
          { path: 'technologies', element: <AdminTechnologyListPage /> },
          { path: 'technologies/new', element: <AdminTechnologyFormPage /> },
          { path: 'technologies/:id', element: <AdminTechnologyFormPage /> },
          { path: 'users', element: <AdminUserListPage /> },
          { path: 'users/:id', element: <AdminUserFormPage /> },
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
