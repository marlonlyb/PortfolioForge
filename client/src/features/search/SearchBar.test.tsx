import type { ReactNode } from 'react';

import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LocaleProvider } from '../../app/providers/LocaleProvider';
import { SearchBar } from './SearchBar';

function renderSearchBar(ui: ReactNode) {
  return render(
    <MemoryRouter>
      <LocaleProvider>{ui}</LocaleProvider>
    </MemoryRouter>,
  );
}

describe('SearchBar', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  afterEach(() => {
    cleanup();
    vi.restoreAllMocks();
  });

  it('uses the shared placeholder by default and allows a presentation-only override', () => {
    const onSearch = vi.fn();

    const { rerender } = renderSearchBar(<SearchBar initialQuery="CompactLogix" onSearch={onSearch} />);

    expect(screen.getByPlaceholderText('Busca proyectos por tecnología, cliente o concepto…')).toBeInTheDocument();

    rerender(
      <MemoryRouter>
        <LocaleProvider>
          <SearchBar
            initialQuery="CompactLogix"
            onSearch={onSearch}
            placeholderOverride="Busca un proyecto, tecnología o tema..."
          />
        </LocaleProvider>
      </MemoryRouter>,
    );

    const input = screen.getByPlaceholderText('Busca un proyecto, tecnología o tema...');
    expect(input).toHaveValue('CompactLogix');

    fireEvent.submit(input.closest('form') as HTMLFormElement);

    expect(onSearch).toHaveBeenCalledWith('CompactLogix');
  });

  it('keeps clear behavior unchanged when the placeholder override is present', () => {
    renderSearchBar(
      <SearchBar
        initialQuery="Printer 05"
        placeholderOverride="Busca un proyecto, tecnología o tema..."
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Limpiar búsqueda' }));

    expect(screen.getByRole('textbox', { name: 'Buscar' })).toHaveValue('');
  });
});
