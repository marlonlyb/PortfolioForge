import { useEffect, useMemo, useState } from 'react';
import { Link, useParams } from 'react-router-dom';

import { fetchProductById } from './api';
import type { ProductDetail } from '../../shared/types/product';
import { AppError } from '../../shared/api/errors';

export function ProductDetailPage() {
  const { id } = useParams<{ id: string }>();

  const [project, setProject] = useState<ProductDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [failedImage, setFailedImage] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    let cancelled = false;
    setLoading(true);
    setError(null);
    setProject(null);
    setFailedImage(null);

    fetchProductById(id)
      .then((data) => {
        if (!cancelled) {
          setProject(data);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof AppError ? err.message : 'Could not load the project.');
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  const imageCandidates = useMemo(() => {
    if (!project) return [] as string[];

    const variantImages = (project.variants ?? [])
      .map((variant) => variant.image_url)
      .filter((image): image is string => Boolean(image));

    return [...variantImages, ...project.images].filter((image, index, all) => all.indexOf(image) === index);
  }, [project]);

  const mainImage = imageCandidates.find((image) => image !== failedImage) ?? null;

  if (loading) {
    return (
      <section className="detail">
        <p className="catalog__loading">Loading project…</p>
      </section>
    );
  }

  if (error || !project) {
    return (
      <section className="detail">
        <div className="card card--muted">
          <p className="eyebrow">Error</p>
          <p>{error ?? 'Project not found.'}</p>
          <Link to="/">Back to projects</Link>
        </div>
      </section>
    );
  }

  return (
    <section className="detail">
      <Link className="detail__back" to="/">
        &larr; Back to projects
      </Link>

      <div className="detail__layout">
        <div className="detail__image-col">
          {mainImage ? (
            <img
              className="detail__image"
              src={mainImage}
              alt={project.name}
              onError={() => setFailedImage(mainImage)}
            />
          ) : (
            <div className="detail__image detail__image--placeholder">Project image pending</div>
          )}
        </div>

        <div className="detail__info">
          {project.category ? <p className="eyebrow">{project.category}</p> : null}
          <h2>{project.name}</h2>
          {project.brand ? <p className="detail__brand">{project.brand}</p> : null}
          <p className="detail__description">{project.description}</p>

          <div className="card card--muted">
            <p className="eyebrow">Skeleton status</p>
            <p>
              This cloned base already supports public project browsing and admin CRUD. The
              next iteration should replace transitional product fields with dedicated portfolio
              fields such as architecture, AI usage, metrics, media, and contact leads.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
