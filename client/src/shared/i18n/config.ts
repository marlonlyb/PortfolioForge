export const PUBLIC_LOCALE = {
  ES: 'es',
  CA: 'ca',
  EN: 'en',
  DE: 'de',
} as const;

export type PublicLocale = (typeof PUBLIC_LOCALE)[keyof typeof PUBLIC_LOCALE];

export const PUBLIC_LOCALE_LABELS: Record<PublicLocale, string> = {
  es: 'ES',
  ca: 'CA',
  en: 'EN',
  de: 'DE',
};

export const TRANSLATION_MODE = {
  AUTO: 'auto',
  MANUAL: 'manual',
} as const;

export type TranslationMode = (typeof TRANSLATION_MODE)[keyof typeof TRANSLATION_MODE];

export const PUBLIC_CONTENT_FIELDS = [
  'name',
  'description',
  'category',
  'business_goal',
  'problem_statement',
  'solution_summary',
  'architecture',
  'ai_usage',
  'integrations',
  'technical_decisions',
  'challenges',
  'results',
  'metrics',
  'timeline',
] as const;

export type PublicContentFieldKey = (typeof PUBLIC_CONTENT_FIELDS)[number];
