import type { PublicLocale } from './config';

export interface LandingQuickPrompt {
  label: string;
  query: string;
}

interface Messages {
  headerTitle: string;
  headerSummary: string;
  navHome: string;
  navLogin: string;
  navSearch: string;
  navAdmin: string;
  navLogout: string;
  headerCaption: string;
  searchPlaceholder: string;
  searchButton: string;
  searchClear: string;
  searchSuggestionsLabel: string;
  landingSearchEyebrow: string;
  landingSearchTitle: string;
  landingSearchLead: string;
  landingSearchPlaceholder: string;
  landingSearchContextHint: string;
  landingQuickPrompts: LandingQuickPrompt[];
  landingEyebrow: string;
  landingTitle: string;
  landingLead: string;
  landingPrimaryCta: string;
  landingSecondaryCta: string;
  landingDesignIntent: string;
  landingPrinciples: string[];
  landingHighlights: Array<{ value: string; label: string }>;
  landingPortfolioSystem: string;
  landingShowcaseTitle: string;
  landingShowcaseCopy: string;
  landingQuoteEyebrow: string;
  landingQuote: string;
  catalogEyebrow: string;
  catalogTitle: string;
  catalogIntro: string;
  catalogSearchLabel: string;
  catalogSearchPlaceholder: string;
  catalogCategoryPlaceholder: string;
  catalogClearFilters: string;
  catalogViewCaseStudy: string;
  catalogOpenProject: string;
  catalogNoImage: string;
  catalogNoSummary: string;
  catalogLoading: string;
  catalogError: string;
  catalogNoMatches: string;
  catalogNoFilterMatches: string;
  detailBack: string;
  detailLoading: string;
  detailNotFound: string;
  detailProjectOverview: string;
  detailCategory: string;
  detailClient: string;
  detailUpdated: string;
  detailTechnologies: string;
  detailIndependent: string;
  detailRecentlyCurated: string;
  detailNotSpecified: string;
  detailStrategyLayer: string;
  detailExecutionLayer: string;
  detailTechnicalLayer: string;
  detailBusinessGoal: string;
  detailProblem: string;
  detailSolution: string;
  detailDeliveryScope: string;
  detailResponsibilityScope: string;
  detailArchitecture: string;
  detailAIUsage: string;
  detailIntegrations: string;
  detailTechnicalDecisions: string;
  detailChallenges: string;
  detailResults: string;
  detailMetrics: string;
  detailTimeline: string;
  detailClientContext: string;
  detailVisualUnavailable: string;
  authPublicEyebrow: string;
  authPublicLoginTitle: string;
  authPublicLoginDescription: string;
  authPublicSignupTitle: string;
  authPublicSignupDescription: string;
  authPublicLocalRestriction: string;
  authAdminEyebrow: string;
  authAdminTitle: string;
  authAdminDescription: string;
  authBackToPortfolio: string;
  authBackToPublicLogin: string;
  detailAssistantLocalRestriction: string;
  searchResultsEyebrow: string;
  searchResultsTitle: string;
  searchResultsIntro: string;
  searchResultsCountSingular: string;
  searchResultsCountPlural: string;
  searchResultsError: string;
  searchResultsSearching: string;
  searchResultsNoResults: string;
  searchResultsViewCatalog: string;
  searchResultsLoadMore: string;
  searchResultsLoadingMore: string;
  searchResultsMinCharacters: string;
  searchFiltersTitle: string;
  searchFiltersCategory: string;
  searchFiltersClient: string;
  searchFiltersTechnologies: string;
  searchFiltersClear: string;
  searchResultProjectVisual: string;
  searchResultMoreTechnologies: string;
  searchResultMatchDetailsAria: string;
  searchContextTitle: string;
  searchContextEvidenceTitle: string;
  searchContextExplanationPrefix: string;
  searchContextExplanationConnector: string;
  searchContextAnd: string;
  searchContextRelevantProjectFields: string;
  searchEvidenceFieldTitle: string;
  searchEvidenceFieldSummary: string;
  searchEvidenceFieldDescription: string;
  searchEvidenceFieldClient: string;
  searchEvidenceFieldCategory: string;
  searchEvidenceFieldTechnology: string;
  searchEvidenceFieldTechnologies: string;
  searchEvidenceFieldSolution: string;
  searchEvidenceFieldArchitecture: string;
  searchEvidenceFieldBusinessGoal: string;
  searchEvidenceFieldAIUsage: string;
  searchEvidenceFieldTechnicalDecisions: string;
  searchEvidenceFieldResults: string;
  searchMatchTypeFTS: string;
  searchMatchTypeFuzzy: string;
  searchMatchTypeSemantic: string;
  searchMatchTypeStructured: string;
  searchResultOpenCaseStudy: string;
  detailCaseStudyEyebrow: string;
  detailErrorEyebrow: string;
  detailProjectHighlightsAria: string;
  detailAdminMarkdownSource: string;
  detailTechnologiesUsedAria: string;
  detailHeroGalleryAria: string;
  detailGalleryFeatured: string;
  detailGalleryOpenImage: string;
  detailGalleryControlsAria: string;
  detailGalleryPreviousImage: string;
  detailGalleryNextImage: string;
  detailGalleryFallbackCaption: string;
  detailGalleryViewFull: string;
  detailAssistantAccessRequirementsAria: string;
  detailAssistantEyebrow: string;
  detailAssistantLoginPrompt: string;
  detailAssistantLoginCta: string;
  detailAssistantVerifyPrompt: string;
  detailAssistantVerifyCta: string;
  detailAssistantCompleteProfilePrompt: string;
  detailAssistantCompleteProfileCta: string;
  detailAssistantGoogleRestriction: string;
  detailLightboxAria: string;
  detailLightboxClose: string;
  detailAssistantToggleOpen: string;
  detailAssistantToggleClose: string;
  detailAssistantPanelAria: string;
  detailAssistantConversationResume: string;
  detailAssistantConversationIntro: string;
  detailAssistantEmpty: string;
  detailAssistantRoleAssistant: string;
  detailAssistantRoleYou: string;
  detailAssistantThinking: string;
  detailAssistantPlaceholder: string;
  detailAssistantClear: string;
  detailAssistantSend: string;
  detailAssistantUnavailable: string;
  authFieldEmail: string;
  authFieldPassword: string;
  authFieldConfirmPassword: string;
  authValidationEmailRequired: string;
  authValidationPasswordInvalid: string;
  authValidationConfirmRequired: string;
  authValidationConfirmMismatch: string;
  authInvalidCredentials: string;
  authForbiddenAdmin: string;
  authPasswordSetupRequired: string;
  authUnexpectedError: string;
  authGoogleCredentialMissing: string;
  authGoogleSignInFailed: string;
  authGoogleContinue: string;
  authGoogleNotConfigured: string;
  authGoogleLoading: string;
  authGoogleCompleting: string;
  authSignupSuccessEyebrow: string;
  authSignupSuccessTitle: string;
  authSignupSuccessHelper: string;
  authSignupSuccessVerifyCta: string;
  authSignupSuccessBackLogin: string;
  authSubmitCreateAccount: string;
  authSubmitCreatingAccount: string;
  authSubmitSignIn: string;
  authSubmitSigningIn: string;
  authAltAlreadyHaveAccount: string;
  authAltNeedAccount: string;
  authVerifyEyebrow: string;
  authVerifyTitle: string;
  authVerifyDefaultNotice: string;
  authVerifyLoginNotice: string;
  authVerifyCodeLabel: string;
  authVerifySubmit: string;
  authVerifySubmitting: string;
  authVerifyResend: string;
  authVerifyResending: string;
  authVerifyResendIn: string;
  authVerifyBackLogin: string;
  authVerifyInvalidCode: string;
  authVerifyExpiredCode: string;
  authVerifyUnableToVerify: string;
  authVerifyUnableToResend: string;
  authCompleteProfileEyebrow: string;
  authCompleteProfileTitle: string;
  authCompleteProfileDescription: string;
  authCompleteProfileFullName: string;
  authCompleteProfileCompany: string;
  authCompleteProfileSave: string;
  authCompleteProfileSaving: string;
  authCompleteProfileBack: string;
  authCompleteProfileUnableToSave: string;
}

type AuthMessageKey =
  | 'authPublicEyebrow'
  | 'authPublicLoginTitle'
  | 'authPublicLoginDescription'
  | 'authPublicSignupTitle'
  | 'authPublicSignupDescription'
  | 'authPublicLocalRestriction'
  | 'authAdminEyebrow'
  | 'authAdminTitle'
  | 'authAdminDescription'
  | 'authBackToPortfolio'
  | 'authBackToPublicLogin'
  | 'detailAssistantLocalRestriction'
  | 'authFieldEmail'
  | 'authFieldPassword'
  | 'authFieldConfirmPassword'
  | 'authValidationEmailRequired'
  | 'authValidationPasswordInvalid'
  | 'authValidationConfirmRequired'
  | 'authValidationConfirmMismatch'
  | 'authInvalidCredentials'
  | 'authForbiddenAdmin'
  | 'authPasswordSetupRequired'
  | 'authUnexpectedError'
  | 'authGoogleCredentialMissing'
  | 'authGoogleSignInFailed'
  | 'authGoogleContinue'
  | 'authGoogleNotConfigured'
  | 'authGoogleLoading'
  | 'authGoogleCompleting'
  | 'authSignupSuccessEyebrow'
  | 'authSignupSuccessTitle'
  | 'authSignupSuccessHelper'
  | 'authSignupSuccessVerifyCta'
  | 'authSignupSuccessBackLogin'
  | 'authSubmitCreateAccount'
  | 'authSubmitCreatingAccount'
  | 'authSubmitSignIn'
  | 'authSubmitSigningIn'
  | 'authAltAlreadyHaveAccount'
  | 'authAltNeedAccount'
  | 'authVerifyEyebrow'
  | 'authVerifyTitle'
  | 'authVerifyDefaultNotice'
  | 'authVerifyLoginNotice'
  | 'authVerifyCodeLabel'
  | 'authVerifySubmit'
  | 'authVerifySubmitting'
  | 'authVerifyResend'
  | 'authVerifyResending'
  | 'authVerifyResendIn'
  | 'authVerifyBackLogin'
  | 'authVerifyInvalidCode'
  | 'authVerifyExpiredCode'
  | 'authVerifyUnableToVerify'
  | 'authVerifyUnableToResend'
  | 'authCompleteProfileEyebrow'
  | 'authCompleteProfileTitle'
  | 'authCompleteProfileDescription'
  | 'authCompleteProfileFullName'
  | 'authCompleteProfileCompany'
  | 'authCompleteProfileSave'
  | 'authCompleteProfileSaving'
  | 'authCompleteProfileBack'
  | 'authCompleteProfileUnableToSave';

type CoreMessages = Omit<Messages, AuthMessageKey>;

const authMessages: Record<PublicLocale, Pick<Messages, AuthMessageKey>> = {
  es: {
    authPublicEyebrow: 'Acceso público',
    authPublicLoginTitle: 'Accede a PortfolioForge',
    authPublicLoginDescription:
      'Usa Google o entra con tu email y contraseña locales.',
    authPublicSignupTitle: 'Crea tu cuenta',
    authPublicSignupDescription:
      'Regístrate con Google o crea una cuenta local con email y contraseña. La verificación del email es obligatoria antes de desbloquear el asistente.',
    authPublicLocalRestriction:
      'Las cuentas locales aún necesitan un email verificado y un perfil completado antes de desbloquear el asistente.',
    authAdminEyebrow: 'Acceso admin',
    authAdminTitle: 'Acceso admin',
    authAdminDescription:
      'Esta ruta permanece oculta en la UI pública. Usa aquí solo credenciales locales de administrador.',
    authBackToPortfolio: 'Volver al portfolio',
    authBackToPublicLogin: 'Volver al login público',
    detailAssistantLocalRestriction:
      'Los usuarios de email sin contraseña aún necesitan un email verificado y un perfil completado antes de desbloquear el asistente del proyecto.',
    authFieldEmail: 'Email',
    authFieldPassword: 'Contraseña',
    authFieldConfirmPassword: 'Confirmar contraseña',
    authValidationEmailRequired: 'El email es obligatorio.',
    authValidationPasswordInvalid: 'La contraseña debe tener al menos 8 caracteres.',
    authValidationConfirmRequired: 'Debes confirmar la contraseña.',
    authValidationConfirmMismatch: 'Las contraseñas deben coincidir.',
    authInvalidCredentials: 'Email o contraseña incorrectos.',
    authForbiddenAdmin: 'Esta cuenta no tiene acceso de administrador.',
    authPasswordSetupRequired: 'Esta cuenta todavía necesita configurar o restablecer una contraseña antes de iniciar sesión.',
    authUnexpectedError: 'Ocurrió un error inesperado. Inténtalo de nuevo.',
    authGoogleCredentialMissing: 'Google no devolvió una credencial válida.',
    authGoogleSignInFailed: 'No se pudo completar el acceso con Google.',
    authGoogleContinue: 'Continuar con Google',
    authGoogleNotConfigured: 'El acceso con Google no está configurado en este entorno.',
    authGoogleLoading: 'Cargando acceso con Google…',
    authGoogleCompleting: 'Completando acceso con Google…',
    authSignupSuccessEyebrow: 'Alta pública',
    authSignupSuccessTitle: 'Revisa tu email',
    authSignupSuccessHelper: 'Verifica primero el código y luego inicia sesión con tu nueva contraseña.',
    authSignupSuccessVerifyCta: 'Verificar email',
    authSignupSuccessBackLogin: 'Volver al login',
    authSubmitCreateAccount: 'Crear cuenta',
    authSubmitCreatingAccount: 'Creando cuenta…',
    authSubmitSignIn: 'Iniciar sesión',
    authSubmitSigningIn: 'Iniciando sesión…',
    authAltAlreadyHaveAccount: '¿Ya tienes cuenta? Inicia sesión',
    authAltNeedAccount: '¿Necesitas una cuenta? Regístrate',
    authVerifyEyebrow: 'Verificación de email',
    authVerifyTitle: 'Verifica tu email',
    authVerifyDefaultNotice: 'Introduce el código de verificación de 6 dígitos que enviamos a tu email.',
    authVerifyLoginNotice: 'Email verificado. Inicia sesión con tu contraseña para continuar.',
    authVerifyCodeLabel: 'Código de verificación',
    authVerifySubmit: 'Verificar email',
    authVerifySubmitting: 'Verificando…',
    authVerifyResend: 'Reenviar código',
    authVerifyResending: 'Enviando…',
    authVerifyResendIn: 'Reenviar código en',
    authVerifyBackLogin: 'Volver al login',
    authVerifyInvalidCode: 'El código de verificación no es válido. Inténtalo de nuevo.',
    authVerifyExpiredCode: 'Este código de verificación ha expirado. Solicita uno nuevo para continuar.',
    authVerifyUnableToVerify: 'No se puede verificar el código ahora mismo.',
    authVerifyUnableToResend: 'No se puede reenviar el código de verificación ahora mismo.',
    authCompleteProfileEyebrow: 'Completa tu perfil',
    authCompleteProfileTitle: 'Desbloquea el asistente del proyecto',
    authCompleteProfileDescription: 'Añade tu nombre completo y empresa para continuar con el chat específico del proyecto.',
    authCompleteProfileFullName: 'Nombre completo',
    authCompleteProfileCompany: 'Empresa',
    authCompleteProfileSave: 'Guardar perfil',
    authCompleteProfileSaving: 'Guardando…',
    authCompleteProfileBack: 'Volver al proyecto',
    authCompleteProfileUnableToSave: 'No se puede guardar tu perfil ahora mismo.',
  },
  ca: {
    authPublicEyebrow: 'Accés públic',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
    authFieldEmail: 'Email',
    authFieldPassword: 'Password',
    authFieldConfirmPassword: 'Confirm password',
    authValidationEmailRequired: 'Email is required.',
    authValidationPasswordInvalid: 'Password must contain at least 8 characters.',
    authValidationConfirmRequired: 'Confirm password is required.',
    authValidationConfirmMismatch: 'Passwords must match.',
    authInvalidCredentials: 'Invalid email or password.',
    authForbiddenAdmin: 'This account does not have admin access.',
    authPasswordSetupRequired: 'This account still needs a password setup or reset before it can sign in.',
    authUnexpectedError: 'An unexpected error occurred. Please try again.',
    authGoogleCredentialMissing: 'Google sign-in did not return a valid credential.',
    authGoogleSignInFailed: 'Unable to complete Google sign-in.',
    authGoogleContinue: 'Continue with Google',
    authGoogleNotConfigured: 'Google sign-in is not configured in this environment.',
    authGoogleLoading: 'Loading Google sign-in…',
    authGoogleCompleting: 'Completing Google sign-in…',
    authSignupSuccessEyebrow: 'Public sign up',
    authSignupSuccessTitle: 'Check your email',
    authSignupSuccessHelper: 'Verify the code first, then log in with your new password.',
    authSignupSuccessVerifyCta: 'Verify email',
    authSignupSuccessBackLogin: 'Back to login',
    authSubmitCreateAccount: 'Create account',
    authSubmitCreatingAccount: 'Creating account…',
    authSubmitSignIn: 'Sign in',
    authSubmitSigningIn: 'Signing in…',
    authAltAlreadyHaveAccount: 'Already have an account? Log in',
    authAltNeedAccount: 'Need an account? Sign up',
    authVerifyEyebrow: 'Email verification',
    authVerifyTitle: 'Verify your email',
    authVerifyDefaultNotice: 'Enter the 6-digit verification code we sent to your email.',
    authVerifyLoginNotice: 'Email verified. Log in with your password to continue.',
    authVerifyCodeLabel: 'Verification code',
    authVerifySubmit: 'Verify email',
    authVerifySubmitting: 'Verifying…',
    authVerifyResend: 'Resend code',
    authVerifyResending: 'Sending…',
    authVerifyResendIn: 'Resend code in',
    authVerifyBackLogin: 'Back to login',
    authVerifyInvalidCode: 'The verification code is invalid. Please try again.',
    authVerifyExpiredCode: 'This verification code expired. Request a new code to continue.',
    authVerifyUnableToVerify: 'Unable to verify the code right now.',
    authVerifyUnableToResend: 'Unable to resend the verification code right now.',
    authCompleteProfileEyebrow: 'Complete your profile',
    authCompleteProfileTitle: 'Unlock the project assistant',
    authCompleteProfileDescription: 'Add your full name and company to continue with project-specific chat.',
    authCompleteProfileFullName: 'Full name',
    authCompleteProfileCompany: 'Company',
    authCompleteProfileSave: 'Save profile',
    authCompleteProfileSaving: 'Saving…',
    authCompleteProfileBack: 'Back to project',
    authCompleteProfileUnableToSave: 'Unable to save your profile right now.',
  },
  en: {
    authPublicEyebrow: 'Public access',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
    authFieldEmail: 'Email',
    authFieldPassword: 'Password',
    authFieldConfirmPassword: 'Confirm password',
    authValidationEmailRequired: 'Email is required.',
    authValidationPasswordInvalid: 'Password must contain at least 8 characters.',
    authValidationConfirmRequired: 'Confirm password is required.',
    authValidationConfirmMismatch: 'Passwords must match.',
    authInvalidCredentials: 'Invalid email or password.',
    authForbiddenAdmin: 'This account does not have admin access.',
    authPasswordSetupRequired: 'This account still needs a password setup or reset before it can sign in.',
    authUnexpectedError: 'An unexpected error occurred. Please try again.',
    authGoogleCredentialMissing: 'Google sign-in did not return a valid credential.',
    authGoogleSignInFailed: 'Unable to complete Google sign-in.',
    authGoogleContinue: 'Continue with Google',
    authGoogleNotConfigured: 'Google sign-in is not configured in this environment.',
    authGoogleLoading: 'Loading Google sign-in…',
    authGoogleCompleting: 'Completing Google sign-in…',
    authSignupSuccessEyebrow: 'Public sign up',
    authSignupSuccessTitle: 'Check your email',
    authSignupSuccessHelper: 'Verify the code first, then log in with your new password.',
    authSignupSuccessVerifyCta: 'Verify email',
    authSignupSuccessBackLogin: 'Back to login',
    authSubmitCreateAccount: 'Create account',
    authSubmitCreatingAccount: 'Creating account…',
    authSubmitSignIn: 'Sign in',
    authSubmitSigningIn: 'Signing in…',
    authAltAlreadyHaveAccount: 'Already have an account? Log in',
    authAltNeedAccount: 'Need an account? Sign up',
    authVerifyEyebrow: 'Email verification',
    authVerifyTitle: 'Verify your email',
    authVerifyDefaultNotice: 'Enter the 6-digit verification code we sent to your email.',
    authVerifyLoginNotice: 'Email verified. Log in with your password to continue.',
    authVerifyCodeLabel: 'Verification code',
    authVerifySubmit: 'Verify email',
    authVerifySubmitting: 'Verifying…',
    authVerifyResend: 'Resend code',
    authVerifyResending: 'Sending…',
    authVerifyResendIn: 'Resend code in',
    authVerifyBackLogin: 'Back to login',
    authVerifyInvalidCode: 'The verification code is invalid. Please try again.',
    authVerifyExpiredCode: 'This verification code expired. Request a new code to continue.',
    authVerifyUnableToVerify: 'Unable to verify the code right now.',
    authVerifyUnableToResend: 'Unable to resend the verification code right now.',
    authCompleteProfileEyebrow: 'Complete your profile',
    authCompleteProfileTitle: 'Unlock the project assistant',
    authCompleteProfileDescription: 'Add your full name and company to continue with project-specific chat.',
    authCompleteProfileFullName: 'Full name',
    authCompleteProfileCompany: 'Company',
    authCompleteProfileSave: 'Save profile',
    authCompleteProfileSaving: 'Saving…',
    authCompleteProfileBack: 'Back to project',
    authCompleteProfileUnableToSave: 'Unable to save your profile right now.',
  },
  de: {
    authPublicEyebrow: 'Öffentlicher Zugang',
    authPublicLoginTitle: 'Login to PortfolioForge',
    authPublicLoginDescription:
      'Use Google or sign in with your local email and password.',
    authPublicSignupTitle: 'Create your account',
    authPublicSignupDescription:
      'Sign up with Google or create a local account with email and password. Email verification is required before the assistant unlocks.',
    authPublicLocalRestriction:
      'Local accounts still need a verified email and a completed profile before the assistant unlocks.',
    authAdminEyebrow: 'Admin access',
    authAdminTitle: 'Admin access',
    authAdminDescription:
      'This route stays hidden from the public UI. Use local admin credentials here only.',
    authBackToPortfolio: 'Back to portfolio',
    authBackToPublicLogin: 'Back to public login',
    detailAssistantLocalRestriction:
      'Passwordless email users still need a verified email and a completed profile before the project assistant unlocks.',
    authFieldEmail: 'Email',
    authFieldPassword: 'Password',
    authFieldConfirmPassword: 'Confirm password',
    authValidationEmailRequired: 'Email is required.',
    authValidationPasswordInvalid: 'Password must contain at least 8 characters.',
    authValidationConfirmRequired: 'Confirm password is required.',
    authValidationConfirmMismatch: 'Passwords must match.',
    authInvalidCredentials: 'Invalid email or password.',
    authForbiddenAdmin: 'This account does not have admin access.',
    authPasswordSetupRequired: 'This account still needs a password setup or reset before it can sign in.',
    authUnexpectedError: 'An unexpected error occurred. Please try again.',
    authGoogleCredentialMissing: 'Google sign-in did not return a valid credential.',
    authGoogleSignInFailed: 'Unable to complete Google sign-in.',
    authGoogleContinue: 'Continue with Google',
    authGoogleNotConfigured: 'Google sign-in is not configured in this environment.',
    authGoogleLoading: 'Loading Google sign-in…',
    authGoogleCompleting: 'Completing Google sign-in…',
    authSignupSuccessEyebrow: 'Public sign up',
    authSignupSuccessTitle: 'Check your email',
    authSignupSuccessHelper: 'Verify the code first, then log in with your new password.',
    authSignupSuccessVerifyCta: 'Verify email',
    authSignupSuccessBackLogin: 'Back to login',
    authSubmitCreateAccount: 'Create account',
    authSubmitCreatingAccount: 'Creating account…',
    authSubmitSignIn: 'Sign in',
    authSubmitSigningIn: 'Signing in…',
    authAltAlreadyHaveAccount: 'Already have an account? Log in',
    authAltNeedAccount: 'Need an account? Sign up',
    authVerifyEyebrow: 'Email verification',
    authVerifyTitle: 'Verify your email',
    authVerifyDefaultNotice: 'Enter the 6-digit verification code we sent to your email.',
    authVerifyLoginNotice: 'Email verified. Log in with your password to continue.',
    authVerifyCodeLabel: 'Verification code',
    authVerifySubmit: 'Verify email',
    authVerifySubmitting: 'Verifying…',
    authVerifyResend: 'Resend code',
    authVerifyResending: 'Sending…',
    authVerifyResendIn: 'Resend code in',
    authVerifyBackLogin: 'Back to login',
    authVerifyInvalidCode: 'The verification code is invalid. Please try again.',
    authVerifyExpiredCode: 'This verification code expired. Request a new code to continue.',
    authVerifyUnableToVerify: 'Unable to verify the code right now.',
    authVerifyUnableToResend: 'Unable to resend the verification code right now.',
    authCompleteProfileEyebrow: 'Complete your profile',
    authCompleteProfileTitle: 'Unlock the project assistant',
    authCompleteProfileDescription: 'Add your full name and company to continue with project-specific chat.',
    authCompleteProfileFullName: 'Full name',
    authCompleteProfileCompany: 'Company',
    authCompleteProfileSave: 'Save profile',
    authCompleteProfileSaving: 'Saving…',
    authCompleteProfileBack: 'Back to project',
    authCompleteProfileUnableToSave: 'Unable to save your profile right now.',
  },
};

const baseMessages: Record<PublicLocale, CoreMessages> = {
  es: {
    headerTitle: 'Portfolio de proyectos',
    headerSummary: 'Estrategia, ejecución y criterio técnico.',
    navHome: 'Inicio',
    navLogin: 'Login',
    navSearch: 'Buscar',
    navAdmin: 'Admin',
    navLogout: 'Salir',
    headerCaption: 'Ing. Marlon Ly Bellido',
    searchPlaceholder: 'Busca proyectos por tecnología, cliente o concepto…',
    searchButton: 'Buscar',
    searchClear: 'Limpiar búsqueda',
    searchSuggestionsLabel: 'Sugerencias de búsqueda',
    landingSearchEyebrow: 'BÚSQUEDA GUIADA',
    landingSearchTitle: '',
    landingSearchLead: 'Busca proyectos, casos y experiencias reales.',
    landingSearchPlaceholder: 'Busca un proyecto, tecnología o tema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Muéstrame la migración PLC de Printer 05', query: 'Printer 05' },
      { label: 'Quiero casos con Allen-Bradley y CompactLogix', query: 'CompactLogix' },
      { label: 'Enséñame automatización industrial con Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Busca motion control con SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Trabajo digital seleccionado',
    landingTitle: 'Una portada más editorial para presentar producto, arquitectura y ejecución.',
    landingLead:
      'PortfolioForge reúne proyectos públicos con mejor jerarquía, bloques más definidos y una lectura más clara del valor detrás de cada entrega.',
    landingPrimaryCta: 'Explorar proyectos',
    landingSecondaryCta: 'Ver catálogo',
    landingDesignIntent: 'Intención visual',
    landingPrinciples: [
      'Casos de estudio con narrativa clara y foco en decisiones.',
      'Composición modular pensada para desktop, tablet y móvil.',
      'Paleta dark mantenida con contraste editorial y ritmo visual.',
    ],
    landingHighlights: [
      { value: '01', label: 'Portfolio público con historias de proyecto estructuradas' },
      { value: '02', label: 'Catálogo modular y responsive para trabajo seleccionado' },
      { value: '03', label: 'Búsqueda y detalle alineados con el mismo sistema visual' },
    ],
    landingPortfolioSystem: 'Sistema de portfolio',
    landingShowcaseTitle: 'Una composición modular para leer el portfolio como publicación, no como listado.',
    landingShowcaseCopy:
      'La landing ahora separa mensaje, exploración y catálogo en bloques más amplios para que el contenido respire mejor en desktop sin perder claridad en móvil.',
    landingQuoteEyebrow: 'Ritmo editorial',
    landingQuote: '“Bloques sólidos, mejor uso del ancho y una estructura visual que prioriza contexto, narrativa y exploración.”',
    catalogEyebrow: 'Índice de proyectos',
    catalogTitle: 'Casos de estudio seleccionados',
    catalogIntro: '',
    catalogSearchLabel: 'Buscar proyectos',
    catalogSearchPlaceholder: 'Buscar por nombre del proyecto…',
    catalogCategoryPlaceholder: 'Categoría',
    catalogClearFilters: 'Limpiar filtros',
    catalogViewCaseStudy: 'Ver caso de estudio',
    catalogOpenProject: 'Abrir proyecto',
    catalogNoImage: 'Sin imagen',
    catalogNoSummary: 'Sin resumen disponible.',
    catalogLoading: 'Cargando proyectos…',
    catalogError: 'Error',
    catalogNoMatches: 'No hay proyectos que coincidan con tu búsqueda.',
    catalogNoFilterMatches: 'No hay proyectos para los filtros actuales.',
    detailBack: '← Volver a proyectos',
    detailLoading: 'Cargando proyecto…',
    detailNotFound: 'Proyecto no encontrado.',
    detailProjectOverview: 'Resumen del proyecto',
    detailCategory: 'Categoría',
    detailClient: 'Contexto',
    detailUpdated: 'Actualizado',
    detailTechnologies: 'Tecnologías',
    detailIndependent: 'Independiente / Interno',
    detailRecentlyCurated: 'Curado recientemente',
    detailNotSpecified: 'No especificado',
    detailStrategyLayer: 'Estrategia',
    detailExecutionLayer: 'Ejecución',
    detailTechnicalLayer: 'Técnica',
    detailBusinessGoal: 'Contexto de negocio / Objetivo',
    detailProblem: 'Problema',
    detailSolution: 'Solución',
    detailDeliveryScope: 'Delivery Scope',
    detailResponsibilityScope: 'Responsibility Scope',
    detailArchitecture: 'Arquitectura',
    detailAIUsage: 'Uso de IA',
    detailIntegrations: 'Integraciones',
    detailTechnicalDecisions: 'Decisiones técnicas',
    detailChallenges: 'Desafíos',
    detailResults: 'Resultados',
    detailMetrics: 'Métricas',
    detailTimeline: 'Timeline',
    detailClientContext: 'Contexto cliente',
    detailVisualUnavailable: 'Visual del caso de estudio no disponible',
    searchResultsEyebrow: 'Búsqueda pública',
    searchResultsTitle: 'Busca proyectos por tecnología, cliente o concepto.',
    searchResultsIntro: 'La búsqueda pública ahora comparte el mismo lenguaje editorial del catálogo y del detalle para que la exploración se sienta consistente.',
    searchResultsCountSingular: 'proyecto encontrado',
    searchResultsCountPlural: 'proyectos encontrados',
    searchResultsError: 'Error al buscar proyectos.',
    searchResultsSearching: 'Buscando',
    searchResultsNoResults: 'No se encontraron proyectos para',
    searchResultsViewCatalog: 'Ver catálogo completo',
    searchResultsLoadMore: 'Cargar más',
    searchResultsLoadingMore: 'Cargando…',
    searchResultsMinCharacters: 'Escribe al menos 2 caracteres para buscar',
    searchFiltersTitle: 'Filtros',
    searchFiltersCategory: 'Categoría',
    searchFiltersClient: 'Cliente',
    searchFiltersTechnologies: 'Tecnologías',
    searchFiltersClear: 'Limpiar filtros',
    searchResultProjectVisual: 'Visual del proyecto',
    searchResultMoreTechnologies: 'más',
    searchResultMatchDetailsAria: 'Detalles de coincidencia',
    searchContextTitle: 'Por qué coincide',
    searchContextEvidenceTitle: 'Evidencia utilizada',
    searchContextExplanationPrefix: 'Coincide con la búsqueda',
    searchContextExplanationConnector: 'en',
    searchContextAnd: 'y',
    searchContextRelevantProjectFields: 'campos relevantes del proyecto',
    searchEvidenceFieldTitle: 'Título del proyecto',
    searchEvidenceFieldSummary: 'Resumen',
    searchEvidenceFieldDescription: 'Descripción',
    searchEvidenceFieldClient: 'Cliente',
    searchEvidenceFieldCategory: 'Categoría',
    searchEvidenceFieldTechnology: 'Tecnología',
    searchEvidenceFieldTechnologies: 'Tecnologías',
    searchEvidenceFieldSolution: 'Solución implementada',
    searchEvidenceFieldArchitecture: 'Arquitectura',
    searchEvidenceFieldBusinessGoal: 'Objetivo de negocio',
    searchEvidenceFieldAIUsage: 'Uso de IA',
    searchEvidenceFieldTechnicalDecisions: 'Decisiones técnicas',
    searchEvidenceFieldResults: 'Resultados',
    searchMatchTypeFTS: 'Texto coincidente',
    searchMatchTypeFuzzy: 'Coincidencia aproximada',
    searchMatchTypeSemantic: 'Coincidencia semántica',
    searchMatchTypeStructured: 'Coincidencia estructurada',
    searchResultOpenCaseStudy: 'Abrir caso de estudio',
    detailCaseStudyEyebrow: 'Caso de estudio',
    detailErrorEyebrow: 'Error',
    detailProjectHighlightsAria: 'Aspectos destacados del proyecto',
    detailAdminMarkdownSource: 'Fuente markdown de admin',
    detailTechnologiesUsedAria: 'Tecnologías utilizadas',
    detailHeroGalleryAria: 'Galería principal del proyecto',
    detailGalleryFeatured: 'Destacada',
    detailGalleryOpenImage: 'Abrir imagen',
    detailGalleryControlsAria: 'Controles de galería',
    detailGalleryPreviousImage: 'Imagen anterior',
    detailGalleryNextImage: 'Imagen siguiente',
    detailGalleryFallbackCaption: 'Visual del proyecto',
    detailGalleryViewFull: 'Ver completa',
    detailAssistantAccessRequirementsAria: 'Requisitos de acceso al asistente',
    detailAssistantEyebrow: 'Asistente del proyecto',
    detailAssistantLoginPrompt: 'Inicia sesión para desbloquear el chat específico del proyecto.',
    detailAssistantLoginCta: 'Iniciar sesión',
    detailAssistantVerifyPrompt: 'Verifica tu email para mantener tu cuenta local habilitada para el asistente.',
    detailAssistantVerifyCta: 'Verificar email',
    detailAssistantCompleteProfilePrompt: 'Completa tu perfil con tu nombre completo y empresa para habilitar el asistente.',
    detailAssistantCompleteProfileCta: 'Completar perfil',
    detailAssistantGoogleRestriction: 'El acceso con Google requiere un email verificado antes de habilitar el asistente.',
    detailLightboxAria: 'Vista previa de imagen',
    detailLightboxClose: 'Cerrar vista previa de imagen',
    detailAssistantToggleOpen: 'Preguntar al asistente del proyecto',
    detailAssistantToggleClose: 'Cerrar asistente',
    detailAssistantPanelAria: 'Asistente del proyecto',
    detailAssistantConversationResume: 'Continúa la conversación con el contexto de esta sesión del navegador solamente.',
    detailAssistantConversationIntro: 'Haz preguntas detalladas apoyadas en la documentación del proyecto.',
    detailAssistantEmpty: 'Prueba a preguntar por arquitectura, resultados, integraciones o tradeoffs.',
    detailAssistantRoleAssistant: 'Asistente',
    detailAssistantRoleYou: 'Tú',
    detailAssistantThinking: 'Pensando…',
    detailAssistantPlaceholder: 'Haz una pregunta detallada sobre este proyecto',
    detailAssistantClear: 'Limpiar chat',
    detailAssistantSend: 'Enviar',
    detailAssistantUnavailable: 'Asistente no disponible.',
  },
  ca: {
    headerTitle: 'Portfoli de projectes',
    headerSummary: 'Estratègia, execució i criteri tècnic.',
    navHome: 'Inici',
    navLogin: 'Login',
    navSearch: 'Cercar',
    navAdmin: 'Admin',
    navLogout: 'Sortir',
    headerCaption: 'Eng. Marlon Ly Bellido',
    searchPlaceholder: 'Cerca projectes per tecnologia, client o concepte…',
    searchButton: 'Cercar',
    searchClear: 'Netejar cerca',
    searchSuggestionsLabel: 'Suggeriments de cerca',
    landingSearchEyebrow: 'CERCA GUIADA',
    landingSearchTitle: '',
    landingSearchLead: 'Cerca projectes, casos i experiències reals.',
    landingSearchPlaceholder: 'Cerca un projecte, tecnologia o tema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Mostra’m la migració PLC de Printer 05', query: 'Printer 05' },
      { label: 'Vull casos amb Allen-Bradley i CompactLogix', query: 'CompactLogix' },
      { label: 'Ensenya’m automatització industrial amb Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Busca motion control amb SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Treball digital seleccionat',
    landingTitle: 'Una portada més editorial per presentar producte, arquitectura i execució.',
    landingLead:
      'PortfolioForge reuneix projectes públics amb millor jerarquia, blocs més definits i una lectura més clara del valor de cada entrega.',
    landingPrimaryCta: 'Explorar projectes',
    landingSecondaryCta: 'Veure catàleg',
    landingDesignIntent: 'Intenció visual',
    landingPrinciples: [
      'Casos d’estudi amb narrativa clara i focus en decisions.',
      'Composició modular pensada per a escriptori, tauleta i mòbil.',
      'Paleta dark mantinguda amb contrast editorial i ritme visual.',
    ],
    landingHighlights: [
      { value: '01', label: 'Portfoli públic amb històries de projecte estructurades' },
      { value: '02', label: 'Catàleg modular i responsive per al treball seleccionat' },
      { value: '03', label: 'Cerca i detall alineats amb el mateix sistema visual' },
    ],
    landingPortfolioSystem: 'Sistema de portfoli',
    landingShowcaseTitle: 'Una composició modular per llegir el portfoli com una publicació, no com un llistat.',
    landingShowcaseCopy:
      'La landing ara separa missatge, exploració i catàleg en blocs més amplis perquè el contingut respiri millor en desktop sense perdre claredat en mòbil.',
    landingQuoteEyebrow: 'Ritme editorial',
    landingQuote: '“Blocs sòlids, millor ús de l’amplada i una estructura visual que prioritza context, narrativa i exploració.”',
    catalogEyebrow: 'Índex de projectes',
    catalogTitle: 'Casos d’estudi seleccionats',
    catalogIntro: '',
    catalogSearchLabel: 'Cercar projectes',
    catalogSearchPlaceholder: 'Cercar pel nom del projecte…',
    catalogCategoryPlaceholder: 'Categoria',
    catalogClearFilters: 'Netejar filtres',
    catalogViewCaseStudy: 'Veure cas d’estudi',
    catalogOpenProject: 'Obrir projecte',
    catalogNoImage: 'Sense imatge',
    catalogNoSummary: 'Sense resum disponible.',
    catalogLoading: 'Carregant projectes…',
    catalogError: 'Error',
    catalogNoMatches: 'No hi ha projectes que coincideixin amb la teva cerca.',
    catalogNoFilterMatches: 'No hi ha projectes per als filtres actuals.',
    detailBack: '← Tornar als projectes',
    detailLoading: 'Carregant projecte…',
    detailNotFound: 'Projecte no trobat.',
    detailProjectOverview: 'Resum del projecte',
    detailCategory: 'Categoria',
    detailClient: 'Context',
    detailUpdated: 'Actualitzat',
    detailTechnologies: 'Tecnologies',
    detailIndependent: 'Independent / Intern',
    detailRecentlyCurated: 'Curat recentment',
    detailNotSpecified: 'No especificat',
    detailStrategyLayer: 'Estratègia',
    detailExecutionLayer: 'Execució',
    detailTechnicalLayer: 'Tècnica',
    detailBusinessGoal: 'Context de negoci / Objectiu',
    detailProblem: 'Problema',
    detailSolution: 'Solució',
    detailDeliveryScope: 'Abast del delivery',
    detailResponsibilityScope: 'Abast de responsabilitat',
    detailArchitecture: 'Arquitectura',
    detailAIUsage: 'Ús d’IA',
    detailIntegrations: 'Integracions',
    detailTechnicalDecisions: 'Decisions tècniques',
    detailChallenges: 'Reptes',
    detailResults: 'Resultats',
    detailMetrics: 'Mètriques',
    detailTimeline: 'Timeline',
    detailClientContext: 'Context client',
    detailVisualUnavailable: 'Visual del cas d’estudi no disponible',
    searchResultsEyebrow: 'Cerca pública',
    searchResultsTitle: 'Cerca projectes per tecnologia, client o concepte.',
    searchResultsIntro: 'La cerca pública ara comparteix el mateix llenguatge editorial del catàleg i del detall perquè l’exploració se senti consistent.',
    searchResultsCountSingular: 'projecte trobat',
    searchResultsCountPlural: 'projectes trobats',
    searchResultsError: 'Error en cercar projectes.',
    searchResultsSearching: 'Cercant',
    searchResultsNoResults: 'No s’han trobat projectes per a',
    searchResultsViewCatalog: 'Veure catàleg complet',
    searchResultsLoadMore: 'Carregar més',
    searchResultsLoadingMore: 'Carregant…',
    searchResultsMinCharacters: 'Escriu almenys 2 caràcters per cercar',
    searchFiltersTitle: 'Filtres',
    searchFiltersCategory: 'Categoria',
    searchFiltersClient: 'Client',
    searchFiltersTechnologies: 'Tecnologies',
    searchFiltersClear: 'Netejar filtres',
    searchResultProjectVisual: 'Visual del projecte',
    searchResultMoreTechnologies: 'més',
    searchResultMatchDetailsAria: 'Detalls de coincidència',
    searchContextTitle: 'Per què coincideix',
    searchContextEvidenceTitle: 'Evidència utilitzada',
    searchContextExplanationPrefix: 'Coincideix amb la cerca',
    searchContextExplanationConnector: 'a',
    searchContextAnd: 'i',
    searchContextRelevantProjectFields: 'camps rellevants del projecte',
    searchEvidenceFieldTitle: 'Títol del projecte',
    searchEvidenceFieldSummary: 'Resum',
    searchEvidenceFieldDescription: 'Descripció',
    searchEvidenceFieldClient: 'Client',
    searchEvidenceFieldCategory: 'Categoria',
    searchEvidenceFieldTechnology: 'Tecnologia',
    searchEvidenceFieldTechnologies: 'Tecnologies',
    searchEvidenceFieldSolution: 'Solució implementada',
    searchEvidenceFieldArchitecture: 'Arquitectura',
    searchEvidenceFieldBusinessGoal: 'Objectiu de negoci',
    searchEvidenceFieldAIUsage: 'Ús d’IA',
    searchEvidenceFieldTechnicalDecisions: 'Decisions tècniques',
    searchEvidenceFieldResults: 'Resultats',
    searchMatchTypeFTS: 'Text coincident',
    searchMatchTypeFuzzy: 'Coincidència aproximada',
    searchMatchTypeSemantic: 'Coincidència semàntica',
    searchMatchTypeStructured: 'Coincidència estructurada',
    searchResultOpenCaseStudy: 'Obrir cas d’estudi',
    detailCaseStudyEyebrow: 'Cas d’estudi',
    detailErrorEyebrow: 'Error',
    detailProjectHighlightsAria: 'Aspectes destacats del projecte',
    detailAdminMarkdownSource: 'Font markdown d’admin',
    detailTechnologiesUsedAria: 'Tecnologies utilitzades',
    detailHeroGalleryAria: 'Galeria principal del projecte',
    detailGalleryFeatured: 'Destacada',
    detailGalleryOpenImage: 'Obrir imatge',
    detailGalleryControlsAria: 'Controls de galeria',
    detailGalleryPreviousImage: 'Imatge anterior',
    detailGalleryNextImage: 'Imatge següent',
    detailGalleryFallbackCaption: 'Visual del projecte',
    detailGalleryViewFull: 'Veure completa',
    detailAssistantAccessRequirementsAria: 'Requisits d’accés a l’assistent',
    detailAssistantEyebrow: 'Assistent del projecte',
    detailAssistantLoginPrompt: 'Inicia sessió per desbloquejar el xat específic del projecte.',
    detailAssistantLoginCta: 'Iniciar sessió',
    detailAssistantVerifyPrompt: 'Verifica el teu correu per mantenir el compte local habilitat per a l’assistent.',
    detailAssistantVerifyCta: 'Verificar correu',
    detailAssistantCompleteProfilePrompt: 'Completa el teu perfil amb el teu nom complet i empresa per habilitar l’assistent.',
    detailAssistantCompleteProfileCta: 'Completar perfil',
    detailAssistantGoogleRestriction: 'L’accés amb Google requereix un correu verificat abans d’habilitar l’assistent.',
    detailLightboxAria: 'Vista prèvia de la imatge',
    detailLightboxClose: 'Tancar la vista prèvia de la imatge',
    detailAssistantToggleOpen: 'Preguntar a l’assistent del projecte',
    detailAssistantToggleClose: 'Tancar assistent',
    detailAssistantPanelAria: 'Assistent del projecte',
    detailAssistantConversationResume: 'Continua la conversa amb el context d’aquesta sessió del navegador només.',
    detailAssistantConversationIntro: 'Fes preguntes detallades recolzades en la documentació del projecte.',
    detailAssistantEmpty: 'Prova de preguntar per arquitectura, resultats, integracions o tradeoffs.',
    detailAssistantRoleAssistant: 'Assistent',
    detailAssistantRoleYou: 'Tu',
    detailAssistantThinking: 'Pensant…',
    detailAssistantPlaceholder: 'Fes una pregunta detallada sobre aquest projecte',
    detailAssistantClear: 'Netejar xat',
    detailAssistantSend: 'Enviar',
    detailAssistantUnavailable: 'Assistent no disponible.',
  },
  en: {
    headerTitle: 'Project portfolio',
    headerSummary: 'Strategy, execution, and technical judgment.',
    navHome: 'Home',
    navLogin: 'Login',
    navSearch: 'Search',
    navAdmin: 'Admin',
    navLogout: 'Logout',
    headerCaption: 'Marlon Ly Bellido · Engineer',
    searchPlaceholder: 'Search projects by technology, client, or concept…',
    searchButton: 'Search',
    searchClear: 'Clear search',
    searchSuggestionsLabel: 'Search suggestions',
    landingSearchEyebrow: 'GUIDED SEARCH',
    landingSearchTitle: '',
    landingSearchLead: 'Search projects, cases, and real-world work.',
    landingSearchPlaceholder: 'Search a project, technology, or topic...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Show me the Printer 05 PLC migration', query: 'Printer 05' },
      { label: 'I want Allen-Bradley and CompactLogix work', query: 'CompactLogix' },
      { label: 'Find industrial automation projects with Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Search motion-control work with SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Selected digital work',
    landingTitle: 'A more editorial cover to present product, architecture, and execution.',
    landingLead:
      'PortfolioForge brings public projects together with stronger hierarchy, clearer blocks, and a more legible reading of the value behind each delivery.',
    landingPrimaryCta: 'Explore projects',
    landingSecondaryCta: 'Browse catalog',
    landingDesignIntent: 'Design intent',
    landingPrinciples: [
      'Case studies with clear narrative and decision-focused structure.',
      'Modular composition designed for desktop, tablet, and mobile.',
      'Dark palette preserved with editorial contrast and visual rhythm.',
    ],
    landingHighlights: [
      { value: '01', label: 'Public portfolio with structured project stories' },
      { value: '02', label: 'Responsive modular catalog for selected work' },
      { value: '03', label: 'Search and detail aligned under the same visual system' },
    ],
    landingPortfolioSystem: 'Portfolio system',
    landingShowcaseTitle: 'A modular composition to read the portfolio like a publication, not a list.',
    landingShowcaseCopy:
      'The landing page now separates message, exploration, and catalog into broader blocks so content can breathe on desktop without losing clarity on mobile.',
    landingQuoteEyebrow: 'Editorial rhythm',
    landingQuote: '“Solid blocks, better use of width, and a visual structure that prioritizes context, narrative, and exploration.”',
    catalogEyebrow: 'Project index',
    catalogTitle: 'Selected case studies',
    catalogIntro: '',
    catalogSearchLabel: 'Search projects',
    catalogSearchPlaceholder: 'Search by project name…',
    catalogCategoryPlaceholder: 'Category',
    catalogClearFilters: 'Clear filters',
    catalogViewCaseStudy: 'View case study',
    catalogOpenProject: 'Open project',
    catalogNoImage: 'No image',
    catalogNoSummary: 'No summary available.',
    catalogLoading: 'Loading projects…',
    catalogError: 'Error',
    catalogNoMatches: 'No projects match your search.',
    catalogNoFilterMatches: 'No projects match the current filters.',
    detailBack: '← Back to projects',
    detailLoading: 'Loading project…',
    detailNotFound: 'Project not found.',
    detailProjectOverview: 'Project overview',
    detailCategory: 'Category',
    detailClient: 'Context',
    detailUpdated: 'Updated',
    detailTechnologies: 'Technologies',
    detailIndependent: 'Independent / Internal',
    detailRecentlyCurated: 'Recently curated',
    detailNotSpecified: 'Not specified',
    detailStrategyLayer: 'Strategy',
    detailExecutionLayer: 'Execution',
    detailTechnicalLayer: 'Technical',
    detailBusinessGoal: 'Business context / goal',
    detailProblem: 'Problem',
    detailSolution: 'Solution',
    detailDeliveryScope: 'Delivery scope',
    detailResponsibilityScope: 'Responsibility scope',
    detailArchitecture: 'Architecture',
    detailAIUsage: 'AI usage',
    detailIntegrations: 'Integrations',
    detailTechnicalDecisions: 'Technical decisions',
    detailChallenges: 'Challenges',
    detailResults: 'Results',
    detailMetrics: 'Metrics',
    detailTimeline: 'Timeline',
    detailClientContext: 'Client context',
    detailVisualUnavailable: 'Case study visual unavailable',
    searchResultsEyebrow: 'Public search',
    searchResultsTitle: 'Search projects by technology, client, or concept.',
    searchResultsIntro: 'Public search now shares the same editorial language as the catalog and detail views so exploration feels consistent.',
    searchResultsCountSingular: 'project found',
    searchResultsCountPlural: 'projects found',
    searchResultsError: 'Error while searching projects.',
    searchResultsSearching: 'Searching',
    searchResultsNoResults: 'No projects were found for',
    searchResultsViewCatalog: 'View full catalog',
    searchResultsLoadMore: 'Load more',
    searchResultsLoadingMore: 'Loading…',
    searchResultsMinCharacters: 'Type at least 2 characters to search',
    searchFiltersTitle: 'Filters',
    searchFiltersCategory: 'Category',
    searchFiltersClient: 'Client',
    searchFiltersTechnologies: 'Technologies',
    searchFiltersClear: 'Clear filters',
    searchResultProjectVisual: 'Project visual',
    searchResultMoreTechnologies: 'more',
    searchResultMatchDetailsAria: 'Match details',
    searchContextTitle: 'Why it matches',
    searchContextEvidenceTitle: 'Evidence used',
    searchContextExplanationPrefix: 'Matches your search',
    searchContextExplanationConnector: 'in',
    searchContextAnd: 'and',
    searchContextRelevantProjectFields: 'relevant project fields',
    searchEvidenceFieldTitle: 'Project title',
    searchEvidenceFieldSummary: 'Summary',
    searchEvidenceFieldDescription: 'Description',
    searchEvidenceFieldClient: 'Client',
    searchEvidenceFieldCategory: 'Category',
    searchEvidenceFieldTechnology: 'Technology',
    searchEvidenceFieldTechnologies: 'Technologies',
    searchEvidenceFieldSolution: 'Implemented solution',
    searchEvidenceFieldArchitecture: 'Architecture',
    searchEvidenceFieldBusinessGoal: 'Business goal',
    searchEvidenceFieldAIUsage: 'AI usage',
    searchEvidenceFieldTechnicalDecisions: 'Technical decisions',
    searchEvidenceFieldResults: 'Results',
    searchMatchTypeFTS: 'Matching text',
    searchMatchTypeFuzzy: 'Approximate match',
    searchMatchTypeSemantic: 'Semantic match',
    searchMatchTypeStructured: 'Structured match',
    searchResultOpenCaseStudy: 'Open case study',
    detailCaseStudyEyebrow: 'Case study',
    detailErrorEyebrow: 'Error',
    detailProjectHighlightsAria: 'Project highlights',
    detailAdminMarkdownSource: 'Admin markdown source',
    detailTechnologiesUsedAria: 'Technologies used',
    detailHeroGalleryAria: 'Main project gallery',
    detailGalleryFeatured: 'Featured',
    detailGalleryOpenImage: 'Open image',
    detailGalleryControlsAria: 'Gallery controls',
    detailGalleryPreviousImage: 'Previous image',
    detailGalleryNextImage: 'Next image',
    detailGalleryFallbackCaption: 'Project visual',
    detailGalleryViewFull: 'View full size',
    detailAssistantAccessRequirementsAria: 'Assistant access requirements',
    detailAssistantEyebrow: 'Project assistant',
    detailAssistantLoginPrompt: 'Log in to unlock project-specific chat.',
    detailAssistantLoginCta: 'Log in',
    detailAssistantVerifyPrompt: 'Verify your email to keep your local account eligible for the assistant.',
    detailAssistantVerifyCta: 'Verify email',
    detailAssistantCompleteProfilePrompt: 'Complete your profile with your full name and company to enable the assistant.',
    detailAssistantCompleteProfileCta: 'Complete profile',
    detailAssistantGoogleRestriction: 'Google sign-in requires a verified email before the assistant can be enabled.',
    detailLightboxAria: 'Image preview',
    detailLightboxClose: 'Close image preview',
    detailAssistantToggleOpen: 'Ask project assistant',
    detailAssistantToggleClose: 'Close assistant',
    detailAssistantPanelAria: 'Project assistant',
    detailAssistantConversationResume: 'Continue the conversation with context from this browser session only.',
    detailAssistantConversationIntro: 'Ask detailed questions grounded in the project documentation.',
    detailAssistantEmpty: 'Try asking about architecture, results, integrations, or tradeoffs.',
    detailAssistantRoleAssistant: 'Assistant',
    detailAssistantRoleYou: 'You',
    detailAssistantThinking: 'Thinking…',
    detailAssistantPlaceholder: 'Ask a detailed question about this project',
    detailAssistantClear: 'Clear chat',
    detailAssistantSend: 'Send',
    detailAssistantUnavailable: 'Assistant unavailable.',
  },
  de: {
    headerTitle: 'Projektportfolio',
    headerSummary: 'Strategie, Umsetzung und technisches Urteilsvermögen.',
    navHome: 'Start',
    navLogin: 'Login',
    navSearch: 'Suche',
    navAdmin: 'Admin',
    navLogout: 'Abmelden',
    headerCaption: 'Ing. Marlon Ly Bellido',
    searchPlaceholder: 'Projekte nach Technologie, Kunde oder Konzept suchen…',
    searchButton: 'Suchen',
    searchClear: 'Suche löschen',
    searchSuggestionsLabel: 'Suchvorschläge',
    landingSearchEyebrow: 'GEFÜHRTE SUCHE',
    landingSearchTitle: '',
    landingSearchLead: 'Suche Projekte, Fälle und reale Arbeit.',
    landingSearchPlaceholder: 'Suche ein Projekt, eine Technologie oder ein Thema...',
    landingSearchContextHint: '',
    landingQuickPrompts: [
      { label: 'Zeig mir die PLC-Migration von Printer 05', query: 'Printer 05' },
      { label: 'Ich suche Arbeiten mit Allen-Bradley und CompactLogix', query: 'CompactLogix' },
      { label: 'Finde Industrial-Automation-Projekte mit Ethernet/IP', query: 'Ethernet/IP' },
      { label: 'Suche Motion-Control mit SEW Eurodrive', query: 'SEW Eurodrive' },
    ],
    landingEyebrow: 'Ausgewählte digitale Arbeit',
    landingTitle: 'Eine editorischere Startseite für Produkt, Architektur und Umsetzung.',
    landingLead:
      'PortfolioForge bündelt öffentliche Projekte mit besserer Hierarchie, klareren Blöcken und einer besser lesbaren Darstellung des Werts hinter jeder Lieferung.',
    landingPrimaryCta: 'Projekte erkunden',
    landingSecondaryCta: 'Katalog ansehen',
    landingDesignIntent: 'Designabsicht',
    landingPrinciples: [
      'Case Studies mit klarer Narration und Fokus auf Entscheidungen.',
      'Modulare Komposition für Desktop, Tablet und Mobile.',
      'Dark-Palette mit editorischem Kontrast und visuellem Rhythmus.',
    ],
    landingHighlights: [
      { value: '01', label: 'Öffentliches Portfolio mit strukturierten Projektgeschichten' },
      { value: '02', label: 'Responsiver modularer Katalog für ausgewählte Arbeiten' },
      { value: '03', label: 'Suche und Detailseiten im selben visuellen System' },
    ],
    landingPortfolioSystem: 'Portfolio-System',
    landingShowcaseTitle: 'Eine modulare Komposition, damit sich das Portfolio wie eine Publikation liest, nicht wie eine Liste.',
    landingShowcaseCopy:
      'Die Landing trennt jetzt Botschaft, Exploration und Katalog in großzügigere Blöcke, damit der Inhalt auf Desktop besser atmet, ohne auf Mobile an Klarheit zu verlieren.',
    landingQuoteEyebrow: 'Editorialer Rhythmus',
    landingQuote: '„Solide Blöcke, bessere Nutzung der Breite und eine visuelle Struktur, die Kontext, Narration und Exploration priorisiert.“',
    catalogEyebrow: 'Projektindex',
    catalogTitle: 'Ausgewählte Fallstudien',
    catalogIntro: '',
    catalogSearchLabel: 'Projekte suchen',
    catalogSearchPlaceholder: 'Nach Projektnamen suchen…',
    catalogCategoryPlaceholder: 'Kategorie',
    catalogClearFilters: 'Filter löschen',
    catalogViewCaseStudy: 'Fallstudie ansehen',
    catalogOpenProject: 'Projekt öffnen',
    catalogNoImage: 'Kein Bild',
    catalogNoSummary: 'Keine Zusammenfassung verfügbar.',
    catalogLoading: 'Projekte werden geladen…',
    catalogError: 'Fehler',
    catalogNoMatches: 'Keine Projekte passen zu deiner Suche.',
    catalogNoFilterMatches: 'Keine Projekte für die aktuellen Filter.',
    detailBack: '← Zurück zu Projekten',
    detailLoading: 'Projekt wird geladen…',
    detailNotFound: 'Projekt nicht gefunden.',
    detailProjectOverview: 'Projektüberblick',
    detailCategory: 'Kategorie',
    detailClient: 'Kontext',
    detailUpdated: 'Aktualisiert',
    detailTechnologies: 'Technologien',
    detailIndependent: 'Unabhängig / Intern',
    detailRecentlyCurated: 'Kürzlich kuratiert',
    detailNotSpecified: 'Nicht angegeben',
    detailStrategyLayer: 'Strategie',
    detailExecutionLayer: 'Umsetzung',
    detailTechnicalLayer: 'Technik',
    detailBusinessGoal: 'Geschäftskontext / Ziel',
    detailProblem: 'Problem',
    detailSolution: 'Lösung',
    detailDeliveryScope: 'Lieferumfang',
    detailResponsibilityScope: 'Verantwortungsbereich',
    detailArchitecture: 'Architektur',
    detailAIUsage: 'KI-Einsatz',
    detailIntegrations: 'Integrationen',
    detailTechnicalDecisions: 'Technische Entscheidungen',
    detailChallenges: 'Herausforderungen',
    detailResults: 'Ergebnisse',
    detailMetrics: 'Metriken',
    detailTimeline: 'Timeline',
    detailClientContext: 'Kundenkontext',
    detailVisualUnavailable: 'Visual der Fallstudie nicht verfügbar',
    searchResultsEyebrow: 'Öffentliche Suche',
    searchResultsTitle: 'Suche Projekte nach Technologie, Kunde oder Konzept.',
    searchResultsIntro: 'Die öffentliche Suche nutzt jetzt dieselbe editorische Sprache wie Katalog und Detailansicht, damit sich die Exploration konsistent anfühlt.',
    searchResultsCountSingular: 'Projekt gefunden',
    searchResultsCountPlural: 'Projekte gefunden',
    searchResultsError: 'Fehler bei der Projektsuche.',
    searchResultsSearching: 'Suche läuft',
    searchResultsNoResults: 'Keine Projekte gefunden für',
    searchResultsViewCatalog: 'Gesamten Katalog ansehen',
    searchResultsLoadMore: 'Mehr laden',
    searchResultsLoadingMore: 'Wird geladen…',
    searchResultsMinCharacters: 'Gib mindestens 2 Zeichen für die Suche ein',
    searchFiltersTitle: 'Filter',
    searchFiltersCategory: 'Kategorie',
    searchFiltersClient: 'Kunde',
    searchFiltersTechnologies: 'Technologien',
    searchFiltersClear: 'Filter löschen',
    searchResultProjectVisual: 'Projektvisual',
    searchResultMoreTechnologies: 'mehr',
    searchResultMatchDetailsAria: 'Trefferdetails',
    searchContextTitle: 'Warum es passt',
    searchContextEvidenceTitle: 'Verwendete Evidenz',
    searchContextExplanationPrefix: 'Passt zu deiner Suche',
    searchContextExplanationConnector: 'in',
    searchContextAnd: 'und',
    searchContextRelevantProjectFields: 'relevanten Projektfeldern',
    searchEvidenceFieldTitle: 'Projekttitel',
    searchEvidenceFieldSummary: 'Zusammenfassung',
    searchEvidenceFieldDescription: 'Beschreibung',
    searchEvidenceFieldClient: 'Kunde',
    searchEvidenceFieldCategory: 'Kategorie',
    searchEvidenceFieldTechnology: 'Technologie',
    searchEvidenceFieldTechnologies: 'Technologien',
    searchEvidenceFieldSolution: 'Umgesetzte Lösung',
    searchEvidenceFieldArchitecture: 'Architektur',
    searchEvidenceFieldBusinessGoal: 'Geschäftsziel',
    searchEvidenceFieldAIUsage: 'KI-Einsatz',
    searchEvidenceFieldTechnicalDecisions: 'Technische Entscheidungen',
    searchEvidenceFieldResults: 'Ergebnisse',
    searchMatchTypeFTS: 'Passender Text',
    searchMatchTypeFuzzy: 'Ungefähre Übereinstimmung',
    searchMatchTypeSemantic: 'Semantische Übereinstimmung',
    searchMatchTypeStructured: 'Strukturierte Übereinstimmung',
    searchResultOpenCaseStudy: 'Fallstudie öffnen',
    detailCaseStudyEyebrow: 'Fallstudie',
    detailErrorEyebrow: 'Fehler',
    detailProjectHighlightsAria: 'Projekthighlights',
    detailAdminMarkdownSource: 'Admin-Markdown-Quelle',
    detailTechnologiesUsedAria: 'Verwendete Technologien',
    detailHeroGalleryAria: 'Hauptgalerie des Projekts',
    detailGalleryFeatured: 'Hervorgehoben',
    detailGalleryOpenImage: 'Bild öffnen',
    detailGalleryControlsAria: 'Galeriesteuerung',
    detailGalleryPreviousImage: 'Vorheriges Bild',
    detailGalleryNextImage: 'Nächstes Bild',
    detailGalleryFallbackCaption: 'Projektvisual',
    detailGalleryViewFull: 'Voll anzeigen',
    detailAssistantAccessRequirementsAria: 'Zugriffsanforderungen für den Assistenten',
    detailAssistantEyebrow: 'Projektassistent',
    detailAssistantLoginPrompt: 'Melde dich an, um den projektspezifischen Chat freizuschalten.',
    detailAssistantLoginCta: 'Anmelden',
    detailAssistantVerifyPrompt: 'Verifiziere deine E-Mail, damit dein lokales Konto für den Assistenten berechtigt bleibt.',
    detailAssistantVerifyCta: 'E-Mail verifizieren',
    detailAssistantCompleteProfilePrompt: 'Vervollständige dein Profil mit deinem vollständigen Namen und Unternehmen, um den Assistenten zu aktivieren.',
    detailAssistantCompleteProfileCta: 'Profil vervollständigen',
    detailAssistantGoogleRestriction: 'Google-Anmeldung erfordert eine verifizierte E-Mail, bevor der Assistent aktiviert werden kann.',
    detailLightboxAria: 'Bildvorschau',
    detailLightboxClose: 'Bildvorschau schließen',
    detailAssistantToggleOpen: 'Projektassistent fragen',
    detailAssistantToggleClose: 'Assistent schließen',
    detailAssistantPanelAria: 'Projektassistent',
    detailAssistantConversationResume: 'Setze die Unterhaltung nur mit dem Kontext dieser Browser-Sitzung fort.',
    detailAssistantConversationIntro: 'Stelle detaillierte Fragen auf Basis der Projektdokumentation.',
    detailAssistantEmpty: 'Frage zum Beispiel nach Architektur, Ergebnissen, Integrationen oder Tradeoffs.',
    detailAssistantRoleAssistant: 'Assistent',
    detailAssistantRoleYou: 'Du',
    detailAssistantThinking: 'Denkt nach…',
    detailAssistantPlaceholder: 'Stelle eine detaillierte Frage zu diesem Projekt',
    detailAssistantClear: 'Chat leeren',
    detailAssistantSend: 'Senden',
    detailAssistantUnavailable: 'Assistent nicht verfügbar.',
  },
};

export function getMessages(locale: PublicLocale): Messages {
  return {
    ...(baseMessages[locale] ?? baseMessages.es),
    ...(authMessages[locale] ?? authMessages.es),
  };
}
