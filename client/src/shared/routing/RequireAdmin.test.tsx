import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router-dom';
import { describe, expect, it } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SessionProvider } from '../../app/providers/SessionProvider';
import { RequireAdmin } from './RequireAdmin';

function LoginDestination() {
  const location = useLocation();
  return <p>{`login redirect:${(location.state as { from?: string } | null)?.from ?? 'missing'}`}</p>;
}

describe('RequireAdmin', () => {
  it('redirects unauthenticated admin requests to /login and preserves state.from', async () => {
    render(
      <MemoryRouter initialEntries={['/admin/projects?tab=active#hero']}>
        <SessionProvider>
          <LocaleProvider>
            <Routes>
              <Route element={<RequireAdmin />}>
                <Route path="/admin/projects" element={<p>admin page</p>} />
              </Route>
              <Route path="/login" element={<LoginDestination />} />
            </Routes>
          </LocaleProvider>
        </SessionProvider>
      </MemoryRouter>,
    );

    expect(await screen.findByText('login redirect:/admin/projects?tab=active#hero')).toBeInTheDocument();
  });
});
