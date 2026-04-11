# Guía de creación de proyectos y relación con el buscador IA

## Objetivo

Este documento explica:

1. cómo crear o editar un proyecto desde el panel admin,
2. qué significa cada campo,
3. cómo impacta cada campo en el buscador híbrido con IA.

---

## Flujo recomendado

### 1. Crear primero las tecnologías

Antes de crear un proyecto, ve a:

- `/admin/technologies`

Ahí defines el catálogo de tecnologías que luego podrás asociar a cada proyecto.

Cada tecnología ayuda al buscador a entender mejor el stack real del proyecto y mejora:

- filtros estructurados,
- coincidencias fuzzy,
- coincidencias semánticas.

---

## 2. Crear un nuevo proyecto

Ve a:

- `/admin/projects/new`

La pantalla está dividida en dos bloques:

1. **Project profile**
2. **Rich Profile (Search Enrichment)**

La idea es esta:

- **Project profile** = datos visibles base del proyecto
- **Rich Profile** = datos que enriquecen la búsqueda y la explicación IA

---

## 3. Significado de cada campo

## A. Project profile

### Title

**Qué es**
- El nombre principal del proyecto.

**Cómo lo usa el sistema**
- Es uno de los campos con más peso en la búsqueda.
- Ayuda mucho en coincidencias directas por nombre.

**Recomendación**
- Usa un título claro y profesional.
- Evita nombres genéricos como `Dashboard 1` o `Proyecto final`.

**Ejemplo**
- `SCADA Commissioning Dashboard`
- `Industrial Monitoring Portal for Siemens`

---

### Summary / Description

**Qué es**
- La descripción principal del proyecto.

**Cómo lo usa el sistema**
- Participa en la búsqueda FTS.
- También entra al texto compuesto usado para embeddings.
- Es requisito para que el proyecto tenga readiness útil.

**Recomendación**
- Escribe al menos una descripción clara de 2 a 4 líneas.
- Incluye dominio, problema y solución general.

**Ejemplo**
- `Aplicación web para monitoreo industrial y soporte de commissioning de tableros Siemens, con visualización de eventos, estado de señales y seguimiento operativo.`

---

### Category

**Qué es**
- Clasificación general del proyecto.

**Cómo lo usa el sistema**
- Se usa como filtro estructurado.
- Ayuda al usuario a agrupar proyectos por tipo.

**Recomendación**
- Usa categorías consistentes.

**Ejemplos útiles**
- `industrial`
- `automation`
- `saas`
- `data`
- `ai`

---

### Client / Context

**Qué es**
- Empresa, cliente, marca o contexto principal del proyecto.

**Cómo lo usa el sistema**
- Tiene mucho peso en la búsqueda.
- Si alguien escribe `SIEMENS`, este campo es clave.

**Recomendación**
- Si hay cliente real, usa su nombre.
- Si no hay cliente, usa un contexto claro.

**Ejemplos**
- `SIEMENS`
- `Industrial automation lab`
- `Internal R&D`

---

### Main images

**Qué es**
- URLs separadas por coma para las imágenes principales.

**Cómo lo usa el sistema**
- Impacta la presentación visual del catálogo y detalle.
- No mejora directamente el ranking del buscador.

**Recomendación**
- Usa una imagen hero clara y profesional.
- Evita dejarlo vacío si el proyecto ya está listo para mostrar.

---

### Published

**Qué es**
- Define si el proyecto aparece públicamente.

**Cómo lo usa el sistema**
- Solo los proyectos publicados deben aparecer en el sitio público y en la búsqueda pública.

**Recomendación**
- Déjalo activado solo cuando el contenido mínimo esté listo.

---

## B. Rich Profile (Search Enrichment)

Estos campos son los más importantes para el buscador inteligente.

### Technologies

**Qué es**
- Lista de tecnologías asociadas al proyecto.

**Cómo lo usa el sistema**
- Mejora filtros estructurados.
- Mejora búsqueda fuzzy.
- Entra al texto compuesto usado para embeddings.
- Mejora mucho la búsqueda por stack.

**Ejemplos**
- `React`
- `Go`
- `PostgreSQL`
- `Siemens`
- `SCADA`

**Importante**
- Si no asignas tecnologías, el proyecto pierde fuerza en búsquedas técnicas.

---

### Business Goal

**Qué es**
- El objetivo de negocio o necesidad operativa del proyecto.

**Cómo lo usa el sistema**
- Entra al documento de búsqueda.
- Ayuda al embedding a entender el contexto.

**Ejemplo**
- `Reducir el tiempo de puesta en marcha y mejorar la visibilidad operativa del commissioning.`

---

### Problem Statement

**Qué es**
- El problema concreto que resolvió el proyecto.

**Cómo lo usa el sistema**
- Entra al documento de búsqueda y al texto de embedding.
- Es especialmente útil cuando alguien busca por problema y no por tecnología.

**Ejemplos de búsqueda que ayuda a resolver**
- `commissioning`
- `troubleshooting`
- `industrial monitoring`

**Ejemplo**
- `El proceso de commissioning dependía de revisión manual dispersa y no existía una vista central del estado de señales, alarmas y validaciones.`

---

### Solution Summary

**Qué es**
- Resumen técnico y funcional de la solución.

**Cómo lo usa el sistema**
- Es uno de los campos de mayor peso para la búsqueda.
- También es clave para la explicación generada por IA.
- Es parte del readiness principal del proyecto.

**Recomendación**
- Este es uno de los mejores campos para escribir bien.
- Resume tecnología, solución y uso.

**Ejemplo**
- `Se desarrolló una interfaz React para soporte de commissioning Siemens, con visualización de señales, estados operativos, trazabilidad y asistencia para validación en campo.`

---

### Architecture

**Qué es**
- Descripción de la arquitectura o integración técnica.

**Cómo lo usa el sistema**
- Participa en FTS y embeddings.
- Ayuda mucho cuando el usuario busca por arquitectura, integración o diseño técnico.

**Ejemplo**
- `Frontend React, backend Go, PostgreSQL para persistencia y servicios de integración para eventos de monitoreo industrial.`

---

### AI Usage

**Qué es**
- Explica si el proyecto usó IA y cómo.

**Cómo lo usa el sistema**
- Entra al documento de búsqueda.
- Ayuda cuando alguien busca términos como `AI`, `LLM`, `automation assistance`, etc.

**Ejemplo**
- `Se utilizó IA para generación de resúmenes operativos y apoyo contextual en búsqueda por evidencia.`

---

## 4. Cómo interactúa todo esto con el buscador

El buscador no depende de un solo campo.

Combina varias capas:

### 1. Filtros estructurados
- categoría
- cliente/contexto
- tecnologías

### 2. Búsqueda léxica
Busca coincidencias en texto usando PostgreSQL FTS.

### 3. Búsqueda fuzzy
Tolera errores o variantes de escritura.

Ejemplo:
- `commisioning` puede aproximarse a `commissioning`

### 4. Búsqueda semántica
Convierte el contenido del proyecto en embeddings para encontrar similitud conceptual.

Esto permite encontrar proyectos incluso si la búsqueda del usuario no coincide palabra por palabra.

---

## 5. Qué campos entran realmente al documento de búsqueda

Actualmente el texto compuesto para búsqueda y embeddings se arma con:

- `name`
- `brand` / client-context
- `description`
- `solution_summary`
- `architecture`
- `business_goal`
- `problem_statement`
- `ai_usage`
- nombres de `technologies`

Además, el índice FTS da pesos distintos:

### Peso alto
- `solution_summary`
- `name`
- `brand`

### Peso medio
- `architecture`
- `description`
- tecnologías

### Peso menor
- `business_goal`
- `problem_statement`
- `ai_usage`

Eso significa que, si quieres mejorar realmente el buscador, los campos más importantes para redactar bien son:

1. `Title`
2. `Client / Context`
3. `Summary / Description`
4. `Solution Summary`
5. `Technologies`

---

## 6. Cómo funciona la explicación IA

Cuando el sistema encuentra un proyecto relevante:

1. recupera evidencia del proyecto,
2. toma la query del usuario,
3. genera una sola oración usando OpenAI,
4. pero solo con evidencia entregada al modelo.

La IA **no debería inventar** ni improvisar contexto fuera del proyecto.

Ejemplo esperado:

> Para la búsqueda «siemens commissioning», este proyecto coincide porque describe trabajo con Siemens en una interfaz para soporte de commissioning y monitoreo industrial.

---

## 7. Qué mejora más la calidad del buscador

### Muy importante
- poner un cliente/contexto claro
- asociar tecnologías reales
- escribir bien `solution_summary`
- escribir bien `problem_statement`
- describir la arquitectura con términos útiles

### Menos importante
- imágenes
- detalles puramente visuales

---

## 8. Readiness de búsqueda

El sistema clasifica el proyecto en tres niveles:

### Incomplete
Falta lo mínimo para buscar bien:
- nombre
- descripción útil

### Basic
Tiene lo mínimo, pero todavía le faltan señales fuertes.

### Complete
Tiene:
- nombre
- descripción suficiente
- categoría
- tecnologías
- `solution_summary`

**Objetivo recomendado:** llevar todos los proyectos públicos al menos a `Complete`.

---

## 9. Recomendación práctica para crear un buen proyecto

Usa este orden:

1. **Title**
2. **Summary / Description**
3. **Category**
4. **Client / Context**
5. **Technologies**
6. **Problem Statement**
7. **Solution Summary**
8. **Architecture**
9. **Business Goal**
10. **AI Usage**
11. Guardar
12. Ejecutar **Actualizar búsqueda** si estás editando uno existente

---

## 10. Ejemplo mínimo de proyecto bien cargado

### Base
- **Title**: `Siemens Commissioning Dashboard`
- **Summary / Description**: `Plataforma para soporte de commissioning y monitoreo industrial con foco en visibilidad operativa y validación en campo.`
- **Category**: `industrial`
- **Client / Context**: `SIEMENS`

### Enrichment
- **Technologies**: `React`, `Go`, `PostgreSQL`, `SCADA`, `Siemens`
- **Business Goal**: `Reducir tiempos de puesta en marcha y centralizar la validación operativa.`
- **Problem Statement**: `La validación de señales y estados se hacía manualmente y sin trazabilidad central.`
- **Solution Summary**: `Se construyó una interfaz para visualización y validación de señales, estados y eventos durante commissioning.`
- **Architecture**: `Frontend React, backend Go, PostgreSQL y servicios de integración industrial.`
- **AI Usage**: `Uso de IA para búsqueda contextual y resúmenes de resultados.`

Con eso ya deberías poder recuperar el proyecto con búsquedas como:

- `siemens`
- `commissioning`
- `react scada`
- `monitoring industrial`

---

## 11. Resumen final

Si quieres que el buscador funcione bien, no basta con “crear el proyecto”.

Debes pensar cada proyecto como una combinación de:

- nombre claro,
- contexto de cliente,
- tecnologías reales,
- problema,
- solución,
- arquitectura,
- y lenguaje útil para búsqueda.

Mientras mejor redactado esté el enrichment, mejor responderá el buscador IA.
