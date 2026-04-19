import { describe, expect, it } from 'vitest';

import type { SessionUser } from '../../app/providers/SessionProvider';
import type { Project } from '../../shared/types/project';
import {
  ASSISTANT_DISCOVERY_VARIANT,
  resolveAssistantAccess,
} from './assistantAccess';

function buildProject(overrides: Partial<Project> = {}): Project {
  return {
    id: 'project-1',
    name: 'PortfolioForge',
    slug: 'portfolioforge',
    description: 'Detailed project description.',
    category: 'platform',
    status: 'published',
    featured: false,
    active: true,
    assistant_available: true,
    images: [],
    media: [],
    created_at: 1710000000,
    updated_at: 1710000000,
    technologies: [],
    ...overrides,
  };
}

function buildSessionUser(overrides: Partial<SessionUser> = {}): SessionUser {
  return {
    id: 'user-1',
    email: 'ada@example.com',
    is_admin: false,
    auth_provider: 'google',
    email_verified: true,
    full_name: 'Ada Lovelace',
    company: 'Analytical Engines',
    profile_completed: true,
    assistant_eligible: true,
    can_use_project_assistant: true,
    created_at: '2026-04-15T00:00:00Z',
    ...overrides,
  };
}

describe('resolveAssistantAccess', () => {
  it('returns no CTA while the session is still loading', () => {
    expect(resolveAssistantAccess(buildProject(), null, true, '/projects/portfolioforge')).toEqual({
      showLauncher: false,
      inlineVariant: null,
    });
  });

  it('returns no CTA when the assistant is unavailable', () => {
    expect(resolveAssistantAccess(buildProject({ assistant_available: false }), null, false, '/projects/portfolioforge')).toEqual({
      showLauncher: false,
      inlineVariant: null,
    });
  });

  it('routes anonymous visitors to login with return state', () => {
    expect(resolveAssistantAccess(buildProject(), null, false, '/projects/portfolioforge')).toEqual({
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.ANONYMOUS,
      ctaTo: '/login',
      ctaState: { from: '/projects/portfolioforge' },
    });
  });

  it('routes local unverified users to verify-email with email state', () => {
    expect(
      resolveAssistantAccess(
        buildProject(),
        buildSessionUser({
          auth_provider: 'local',
          email_verified: false,
          assistant_eligible: false,
          can_use_project_assistant: false,
        }),
        false,
        '/projects/portfolioforge',
      ),
    ).toEqual({
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.VERIFY_EMAIL,
      ctaTo: '/verify-email',
      ctaState: { from: '/projects/portfolioforge', email: 'ada@example.com' },
    });
  });

  it('routes verified users with incomplete profile to complete-profile', () => {
    expect(
      resolveAssistantAccess(
        buildProject(),
        buildSessionUser({
          profile_completed: false,
          assistant_eligible: false,
          can_use_project_assistant: false,
        }),
        false,
        '/projects/portfolioforge',
      ),
    ).toEqual({
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.COMPLETE_PROFILE,
      ctaTo: '/complete-profile',
      ctaState: { from: '/projects/portfolioforge' },
    });
  });

  it('returns a restricted informational card for residual non-eligible users', () => {
    expect(
      resolveAssistantAccess(
        buildProject(),
        buildSessionUser({
          auth_provider: 'google',
          email_verified: false,
          assistant_eligible: false,
          can_use_project_assistant: false,
        }),
        false,
        '/projects/portfolioforge',
      ),
    ).toEqual({
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.RESTRICTED,
    });
  });

  it('returns launcher-only access for eligible users', () => {
    expect(resolveAssistantAccess(buildProject(), buildSessionUser(), false, '/projects/portfolioforge')).toEqual({
      showLauncher: true,
      inlineVariant: null,
    });
  });
});
