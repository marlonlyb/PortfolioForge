import type { SessionUser } from '../../app/providers/SessionProvider';
import type { Project } from '../../shared/types/project';

export const ASSISTANT_DISCOVERY_VARIANT = {
  ANONYMOUS: 'anonymous',
  VERIFY_EMAIL: 'verify-email',
  COMPLETE_PROFILE: 'complete-profile',
  RESTRICTED: 'restricted',
} as const;

export type AssistantDiscoveryVariant =
  (typeof ASSISTANT_DISCOVERY_VARIANT)[keyof typeof ASSISTANT_DISCOVERY_VARIANT];

interface AssistantAccessCtaState {
  from: string;
  email?: string;
}

interface AssistantAccessResolutionBase {
  showLauncher: boolean;
  inlineVariant: AssistantDiscoveryVariant | null;
}

interface AssistantAccessResolutionWithCta extends AssistantAccessResolutionBase {
  ctaTo: '/login' | '/verify-email' | '/complete-profile';
  ctaState: AssistantAccessCtaState;
}

interface AssistantAccessResolutionWithoutCta extends AssistantAccessResolutionBase {
  ctaTo?: undefined;
  ctaState?: undefined;
}

export type AssistantAccessResolution =
  | AssistantAccessResolutionWithCta
  | AssistantAccessResolutionWithoutCta;

type AssistantAccessProject = Pick<Project, 'assistant_available'> | null | undefined;

export function resolveAssistantAccess(
  project: AssistantAccessProject,
  user: SessionUser | null,
  sessionLoading: boolean,
  fromPath: string,
): AssistantAccessResolution {
  if (sessionLoading || !project?.assistant_available) {
    return {
      showLauncher: false,
      inlineVariant: null,
    };
  }

  if (!user) {
    return {
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.ANONYMOUS,
      ctaTo: '/login',
      ctaState: { from: fromPath },
    };
  }

  if (user.can_use_project_assistant) {
    return {
      showLauncher: true,
      inlineVariant: null,
    };
  }

  if (user.auth_provider === 'local' && !user.email_verified) {
    return {
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.VERIFY_EMAIL,
      ctaTo: '/verify-email',
      ctaState: {
        from: fromPath,
        email: user.email,
      },
    };
  }

  if (!user.profile_completed) {
    return {
      showLauncher: false,
      inlineVariant: ASSISTANT_DISCOVERY_VARIANT.COMPLETE_PROFILE,
      ctaTo: '/complete-profile',
      ctaState: { from: fromPath },
    };
  }

  return {
    showLauncher: false,
    inlineVariant: ASSISTANT_DISCOVERY_VARIANT.RESTRICTED,
  };
}
