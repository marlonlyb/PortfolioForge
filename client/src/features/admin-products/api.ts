export {
  createAdminProject as createProduct,
  fetchAdminProjectById,
  fetchAdminProjects as fetchAdminProducts,
  fetchProjectLocalizations,
  fetchProjectReadiness,
  reembedProject,
  reembedStale,
  saveProjectLocalizations,
  updateAdminProject as updateProduct,
  updateAdminProjectStatus as updateProductStatus,
  updateProjectEnrichment,
} from '../admin-projects/api';

export type {
  AdminProjectLocalizationLocale,
  AdminProjectLocalizationsResponse,
  CreateAdminProjectPayload as CreateProductPayload,
  CreateAdminProjectVariantPayload as CreateVariantPayload,
  LocalizedAdminField,
  ProjectReadiness,
  ReembedResponse,
  SaveProjectLocalizationsPayload,
  UpdateAdminProjectPayload as UpdateProductPayload,
  UpdateAdminProjectStatusPayload as UpdateProductStatusPayload,
  UpdateAdminProjectVariantPayload as UpdateVariantPayload,
  UpdateProjectEnrichmentPayload,
  UpdateProjectEnrichmentProfilePayload,
} from '../admin-projects/api';

export type { AdminProjectListResponse as AdminProductListResponse } from '../../shared/types/admin-project';
