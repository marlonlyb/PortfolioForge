import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { LandingPage } from './LandingPage';

vi.mock('../catalog/CatalogPage', () => ({
  CatalogPage: () => <div data-testid="catalog-page" />,
}));

describe('LandingPage', () => {
  beforeEach(() => {
    window.localStorage.clear();
    document.documentElement.lang = 'es';
    Object.defineProperty(HTMLElement.prototype, 'scrollIntoView', {
      configurable: true,
      value: vi.fn(),
    });
    vi.spyOn(HTMLElement.prototype, 'scrollIntoView').mockImplementation(() => {});
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it('renders the polished guided-search copy and maps prompt labels to deterministic queries', () => {
    render(
      <MemoryRouter>
        <LocaleProvider>
          <LandingPage />
        </LocaleProvider>
      </MemoryRouter>,
    );

    expect(screen.getByText('BÚSQUEDA GUIADA')).toBeInTheDocument();
    expect(screen.queryByRole('heading', { level: 1 })).not.toBeInTheDocument();
    expect(screen.getByText('Busca proyectos, casos y experiencias reales.')).toBeInTheDocument();

    const prompt = screen.getByRole('button', {
      name: 'Muéstrame la migración PLC de Printer 05',
    });

    fireEvent.click(prompt);

    expect(screen.getByRole('textbox', { name: 'Buscar' })).toHaveValue('Printer 05');
    expect(HTMLElement.prototype.scrollIntoView).toHaveBeenCalledTimes(1);
  });
});
