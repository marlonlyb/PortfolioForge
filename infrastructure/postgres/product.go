package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/marlonlyb/portfolioforge/model"
)

const pTable = "products"

var pFields = []string{
	"id",
	"product_name",
	"price",
	"images",
	"description",
	"features",
	"name",
	"slug",
	"category",
	"brand",
	"active",
	"source_markdown_url",
	"created_at",
	"updated_at",
}

var (
	pPsqlInsert = BuildSQLInsert(pTable, pFields)
	pPsqlUpdate = BuildSQLUpdatedByID(pTable, pFields)
	pPsqlDelete = BuildSQLDelete(pTable)
)

// ProjectCatalogRepository is the canonical write/admin repository for projects.
// It still persists against the legacy `products` table and related product_* tables.
type ProjectCatalogRepository struct {
	db *pgxpool.Pool
}

func NewProjectCatalogRepository(db *pgxpool.Pool) ProjectCatalogRepository {
	return ProjectCatalogRepository{db: db}
}

func (p ProjectCatalogRepository) Create(m *model.AdminProjectWrite) error {
	legacy := m.ToLegacyProduct(true)
	_, err := p.db.Exec(
		context.Background(),
		pPsqlInsert,
		legacy.ID,
		legacy.ProductName,
		legacy.Price,
		legacy.Images,
		legacy.Description,
		legacy.Features,
		NullIfEmpty(legacy.Name),
		NullIfEmpty(legacy.Slug),
		NullIfEmpty(legacy.Category),
		NullIfEmpty(legacy.Brand),
		legacy.Active,
		NullIfEmpty(legacy.SourceMarkdownURL),
		legacy.CreatedAt,
		Int64ToNull(legacy.UpdatedAt),
	)

	if err != nil {
		return err
	}
	return nil
}

func (p ProjectCatalogRepository) Update(m *model.AdminProjectWrite) error {
	legacy := m.ToLegacyProduct(true)
	_, err := p.db.Exec(
		context.Background(),
		pPsqlUpdate,
		legacy.ProductName,
		legacy.Price,
		legacy.Images,
		legacy.Description,
		legacy.Features,
		NullIfEmpty(legacy.Name),
		NullIfEmpty(legacy.Slug),
		NullIfEmpty(legacy.Category),
		NullIfEmpty(legacy.Brand),
		legacy.Active,
		NullIfEmpty(legacy.SourceMarkdownURL),
		legacy.UpdatedAt,
		legacy.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) Delete(ID uuid.UUID) error {
	_, err := p.db.Exec(
		context.Background(),
		pPsqlDelete,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) UpdateActive(ID uuid.UUID, active bool) error {
	_, err := p.db.Exec(
		context.Background(),
		"UPDATE products SET active = $1, updated_at = EXTRACT(EPOCH FROM NOW())::int WHERE id = $2",
		active,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) CreateVariants(projectID uuid.UUID, variants []model.AdminProjectVariantInput) error {
	if len(variants) == 0 {
		return nil
	}

	query := `
		INSERT INTO product_variants (id, product_id, sku, color, size, price, stock, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())`

	for _, variant := range variants {
		v := variant.ToLegacy(projectID)
		variantID := v.ID
		if variantID == uuid.Nil {
			var err error
			variantID, err = uuid.NewUUID()
			if err != nil {
				return err
			}
		}

		imageURL := sql.NullString{}
		if v.ImageURL != "" {
			imageURL.String = v.ImageURL
			imageURL.Valid = true
		}

		_, err := p.db.Exec(
			context.Background(),
			query,
			variantID,
			projectID,
			v.SKU,
			v.Color,
			v.Size,
			v.Price,
			v.Stock,
			imageURL,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p ProjectCatalogRepository) DeleteVariantsByProductID(productID uuid.UUID) error {
	_, err := p.db.Exec(
		context.Background(),
		"DELETE FROM product_variants WHERE product_id = $1",
		productID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) ReplaceMedia(productID uuid.UUID, media []model.ProjectMedia) error {
	return replaceProjectMedia(context.Background(), p.db, productID, media)
}

func (p ProjectCatalogRepository) DeleteVariantByID(ID uuid.UUID) error {
	_, err := p.db.Exec(
		context.Background(),
		"DELETE FROM product_variants WHERE id = $1",
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) UpdateVariant(variant model.AdminProjectVariantInput, projectID uuid.UUID) error {
	v := variant.ToLegacy(projectID)
	imageURL := sql.NullString{}
	if v.ImageURL != "" {
		imageURL.String = v.ImageURL
		imageURL.Valid = true
	}

	_, err := p.db.Exec(
		context.Background(),
		`UPDATE product_variants
		 SET sku = $1,
		     color = $2,
		     size = $3,
		     price = $4,
		     stock = $5,
		     image_url = $6,
		     updated_at = NOW()
		 WHERE id = $7 AND product_id = $8`,
		v.SKU,
		v.Color,
		v.Size,
		v.Price,
		v.Stock,
		imageURL,
		v.ID,
		v.ProductID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p ProjectCatalogRepository) GetStoreByID(ID uuid.UUID) (model.StoreProduct, error) {
	products, err := p.getStoreProducts("WHERE p.id = $1", ID)
	if err != nil {
		return model.StoreProduct{}, err
	}

	if len(products) == 0 {
		return model.StoreProduct{}, pgx.ErrNoRows
	}

	return products[0], nil
}

func (p ProjectCatalogRepository) GetAdminByID(ID uuid.UUID) (model.AdminProject, error) {
	products, err := p.getStoreProducts("WHERE p.id = $1", ID)
	if err != nil {
		return model.AdminProject{}, err
	}

	if len(products) == 0 {
		return model.AdminProject{}, pgx.ErrNoRows
	}

	adminProject := model.AdminProjectFromStoreProduct(products[0])
	projectDetails, err := NewProjectRepository(p.db).GetByID(context.Background(), ID)
	if err == nil {
		adminProject.ApplyEnrichment(projectDetails)
	}

	return adminProject, nil
}

func (p ProjectCatalogRepository) GetStoreAll() ([]model.StoreProduct, error) {
	return p.getStoreProducts("WHERE p.active = TRUE", nil)
}

func (p ProjectCatalogRepository) GetAdminAll() ([]model.AdminProject, error) {
	products, err := p.getStoreProducts("", nil)
	if err != nil {
		return nil, err
	}

	return model.AdminProjectsFromStoreProducts(products), nil
}

func (p ProjectCatalogRepository) getStoreProducts(whereClause string, arg interface{}) ([]model.StoreProduct, error) {
	query := `
		SELECT
			p.id,
			COALESCE(NULLIF(p.name, ''), p.product_name) AS name,
			COALESCE(NULLIF(p.slug, ''), regexp_replace(lower(COALESCE(NULLIF(p.name, ''), p.product_name)), '[^a-z0-9]+', '-', 'g')) AS slug,
			p.description,
			COALESCE(NULLIF(p.category, ''), 'general') AS category,
			COALESCE(p.brand, '') AS brand,
			COALESCE(p.source_markdown_url, '') AS source_markdown_url,
			p.images,
			p.active,
			v.id,
			v.product_id,
			v.sku,
			v.color,
			v.size,
			v.price,
			v.stock,
			COALESCE(v.image_url, '')
		FROM products p
		LEFT JOIN product_variants v ON v.product_id = p.id
	` + whereClause + `
		ORDER BY p.created_at DESC, v.color ASC, v.size ASC`

	var rows pgx.Rows
	var err error
	if arg == nil {
		rows, err = p.db.Query(context.Background(), query)
	} else {
		rows, err = p.db.Query(context.Background(), query, arg)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productsMap := map[uuid.UUID]*model.StoreProduct{}
	orderedIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var (
			productID         uuid.UUID
			name              string
			slug              string
			description       string
			category          string
			brand             string
			sourceMarkdownURL string
			imagesRaw         []byte
			active            bool
			variantID         uuid.NullUUID
			variantProd       uuid.NullUUID
			sku               sql.NullString
			color             sql.NullString
			size              sql.NullString
			price             sql.NullFloat64
			stock             sql.NullInt64
			imageURL          sql.NullString
		)

		if err = rows.Scan(
			&productID,
			&name,
			&slug,
			&description,
			&category,
			&brand,
			&sourceMarkdownURL,
			&imagesRaw,
			&active,
			&variantID,
			&variantProd,
			&sku,
			&color,
			&size,
			&price,
			&stock,
			&imageURL,
		); err != nil {
			return nil, err
		}

		productData, exists := productsMap[productID]
		if !exists {
			images := []string{}
			if len(imagesRaw) > 0 {
				_ = json.Unmarshal(imagesRaw, &images)
			}

			productData = &model.StoreProduct{
				ID:                productID,
				Name:              name,
				Slug:              slugify(slug),
				Description:       description,
				Category:          category,
				Brand:             strings.TrimSpace(brand),
				SourceMarkdownURL: strings.TrimSpace(sourceMarkdownURL),
				Images:            images,
				Active:            active,
				Variants:          []model.StoreProductVariant{},
			}
			productsMap[productID] = productData
			orderedIDs = append(orderedIDs, productID)
		}

		if variantID.Valid && variantProd.Valid && sku.Valid && color.Valid && size.Valid && price.Valid && stock.Valid {
			variant := model.StoreProductVariant{
				ID:        variantID.UUID,
				ProductID: variantProd.UUID,
				SKU:       sku.String,
				Color:     color.String,
				Size:      size.String,
				Price:     price.Float64,
				Stock:     int(stock.Int64),
			}
			if imageURL.Valid {
				variant.ImageURL = imageURL.String
			}
			productData.Variants = append(productData.Variants, variant)
		}
	}

	products := make([]model.StoreProduct, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		productData := productsMap[id]
		media, mediaErr := fetchProjectMedia(context.Background(), p.db, id)
		if mediaErr == nil {
			productData.Media = media
			productData.Images = model.BuildProjectImageList(media, productData.Images)
		}
		decorateStoreProduct(productData)
		products = append(products, *productData)
	}

	return products, nil
}

func decorateStoreProduct(productData *model.StoreProduct) {
	colors := map[string]struct{}{}
	sizes := map[string]struct{}{}
	priceFrom := 0.0
	for i, variant := range productData.Variants {
		if i == 0 || variant.Price < priceFrom {
			priceFrom = variant.Price
		}
		colors[variant.Color] = struct{}{}
		sizes[variant.Size] = struct{}{}
	}

	if len(productData.Variants) > 0 {
		productData.PriceFrom = priceFrom
		productData.AvailableColors = mapKeys(colors)
		productData.AvailableSizes = mapKeys(sizes)
	}
}

func mapKeys(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func slugify(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	return strings.Trim(value, "-")
}
