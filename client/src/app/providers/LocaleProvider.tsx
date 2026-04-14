import { createContext, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';

import { PUBLIC_LOCALE, type PublicLocale } from '../../shared/i18n/config';
import { getMessages } from '../../shared/i18n/messages';

const STORAGE_KEY = 'portfolioforge.locale';

interface LocaleContextValue {
  locale: PublicLocale;
  setLocale: (locale: PublicLocale) => void;
  t: ReturnType<typeof getMessages>;
}

const LocaleContext = createContext<LocaleContextValue | null>(null);

export function LocaleProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<PublicLocale>(PUBLIC_LOCALE.ES);

  useEffect(() => {
    const stored = window.localStorage.getItem(STORAGE_KEY) as PublicLocale | null;
    if (stored && Object.values(PUBLIC_LOCALE).includes(stored)) {
      setLocaleState(stored);
    }
  }, []);

  function setLocale(nextLocale: PublicLocale) {
    setLocaleState(nextLocale);
    window.localStorage.setItem(STORAGE_KEY, nextLocale);
    document.documentElement.lang = nextLocale;
  }

  const value = useMemo(
    () => ({ locale, setLocale, t: getMessages(locale) }),
    [locale],
  );

  return <LocaleContext.Provider value={value}>{children}</LocaleContext.Provider>;
}

export function useLocale() {
  const context = useContext(LocaleContext);
  if (!context) {
    throw new Error('useLocale must be used within LocaleProvider');
  }
  return context;
}
