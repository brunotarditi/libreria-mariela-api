# 📚 Librería Mariela API

![Go Version](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go&logoColor=white)
![Gin Framework](https://img.shields.io/badge/Gin-v1.10.0-008ECF?style=flat-square&logo=gin&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-v1.25.12-7952B3?style=flat-square)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker&logoColor=white)

**Librería Mariela API** es el backend RESTful diseñado para gestionar de manera integral las operaciones comerciales y de inventario de una librería y papelería comercial/escolar.

El sistema administra catálogo de productos, control de stock automatizado, registro transaccional de compras y ventas con cálculo de **costo promedio ponderado (PPP)**, listas de precios, presupuestos en PDF, importación/exportación masiva mediante hojas de cálculo Excel, registro de auditoría de peticiones y autenticación OAuth 2.0 delegada con Peak Auth.

---

## 📑 Tabla de Contenidos

- [Características Principales](#-características-principales)
- [Stack Tecnológico](#-stack-tecnológico)
- [Arquitectura del Proyecto](#-arquitectura-del-proyecto)
- [Requisitos Previos](#-requisitos-previos)
- [Variables de Entorno](#-variables-de-entorno)
- [Instalación y Puesta en Marcha](#-instalación-y-puesta-en-marcha)
  - [Opción 1: Desarrollo Local (Go + PostgreSQL)](#opción-1-desarrollo-local-go--postgresql)
  - [Opción 2: Despliegue con Docker Compose](#opción-2-despliegue-con-docker-compose)
- [Referencia de la API (Endpoints)](#-referencia-de-la-api-endpoints)
  - [Autenticación](#autenticación)
  - [Salud del Sistema](#salud-del-sistema)
  - [Dashboard](#dashboard)
  - [Productos](#productos)
  - [Categorías y Marcas](#categorías-y-marcas)
  - [Clientes y Proveedores](#clientes-y-proveedores)
  - [Compras](#compras)
  - [Ventas](#ventas)
  - [Stock y Movimientos](#stock-y-movimientos)
  - [Listas de Precios](#listas-de-precios)
  - [Presupuestos (Budgets)](#presupuestos-budgets)
- [Lógica de Negocio Destacada](#-lógica-de-negocio-destacada)
- [Licencia](#-licencia)

---

## ✨ Características Principales

- **Gestión de Catálogo**: Categorías, marcas y productos con código interno, SKU y márgenes de ganancia asignados.
- **Control Transaccional de Stock**:
  - Las compras incrementan el stock e insertan movimientos de entrada (`IN`).
  - Las ventas validan existencias, descuentan el stock e insertan movimientos de salida (`OUT`).
  - Posibilidad de anulación de compras y ventas con reversión automática y consistente de existencias.
- **Costeo Inteligente y Fijación de Precios**:
  - Cálculo de costo unitario promedio ponderado según el historial de compras del producto.
  - Determinación automática del precio de venta sugerido basado en el margen de ganancia configurado:
    $$\text{Precio} = \text{Costo Promedio} \times \left(1 + \frac{\text{Margen}}{100}\right)$$
- **Procesamiento de Archivos Excel (.xlsx)**:
  - Exportación dinámica del catálogo completo de productos con formato estilizado.
  - Importación y actualización masiva de productos desde archivos Excel vía `multipart/form-data`.
- **Generación de Presupuestos en PDF**:
  - Creación de presupuestos comerciales formateados con encabezado, logotipo institucional y tablas utilizando `maroto/v2`.
- **Auditoría Transaccional Automática**:
  - Middleware que registra asíncronamente (vía Goroutines) cada petición procesada: ruta, método HTTP, IP, marca temporal y usuario.
- **Autenticación OAuth 2.0**:
  - Endpoint de intercambio de código de autorización (`code exchange`) contra el servidor de identidad Peak Auth.

---

## 🛠️ Stack Tecnológico

| Componente | Tecnología |
| :--- | :--- |
| **Lenguaje** | [Go (Golang)](https://go.dev/) 1.27+ |
| **Framework Web** | [Gin Gonic](https://gin-gonic.com/) |
| **ORM** | [GORM](https://gorm.io/) |
| **Base de Datos** | [PostgreSQL](https://www.postgresql.org/) 16 |
| **Manipulación Excel** | [Excelize v2](https://github.com/qax-os/excelize) |
| **Generación de PDF** | [Maroto v2](https://github.com/johnfercher/maroto) |
| **Contenerización** | [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/) |

---

## 📂 Arquitectura del Proyecto

El proyecto sigue un patrón por capas desacoplado con soporte de genéricos en Go para operaciones CRUD comunes:

```text
libreria-mariela-api/
├── app/                  # Contenedor de la aplicación y dependencias globales
├── assets/               # Plantillas, recursos estáticos y logos (PDF/Excel)
│   ├── public/           # Archivos multimedia públicos (logos)
│   └── templates/        # Plantillas de Excel generadas o exportadas
├── common/               # Operaciones y controladores genéricos CRUD (Operations & Handlers)
├── constants/            # Constantes del dominio (tipos de movimientos de stock, etc.)
├── controllers/          # Manejadores HTTP que reciben las peticiones y orquestan servicios
├── db/                   # Inicialización, conexión de base de datos y AutoMigrate de GORM
├── middlewares/          # Middlewares de Gin (CORS, registro de auditoría en BD)
├── models/               # Entidades de la base de datos mapeadas con GORM
├── repositories/         # Capa de acceso a datos y consultas específicas
├── requests/             # DTOs de entrada con reglas de binding y validaciones
├── responses/            # DTOs de salida y estructuras serializables para el cliente
├── services/             # Lógica de negocio (transacciones, cálculo de stock/costos, PDF, Excel)
├── utils/                # Utilidades de cálculo matemático/financiero y renderizado
├── Dockerfile            # Construcción multi-etapa optimizada para producción
├── docker-compose.yml    # Orquestación de API, PostgreSQL 16 y pgAdmin
├── main.go               # Punto de entrada de la aplicación y configuración de middlewares
└── routes.go             # Declaración y mapeo de rutas públicas y privadas
```

---

## 📋 Requisitos Previos

- **Go**: 1.27 o superior instalado en el sistema.
- **PostgreSQL**: 13 o superior (o Docker para ejecutarlo en un contenedor).
- **Docker y Docker Compose** (opcional, para ejecución contenerizada).

---

## ⚙️ Variables de Entorno

Copia el archivo de ejemplo `.env.example` a `.env` en la raíz del proyecto:

```bash
cp .env.example .env
```

Configura las siguientes variables según tu entorno:

| Variable | Descripción | Ejemplo / Valor por Defecto |
| :--- | :--- | :--- |
| `ENV` | Entorno de ejecución (`development` / `production`) | `development` |
| `DB_HOST` | Host de la base de datos PostgreSQL | `localhost` (o `db` en Docker) |
| `DB_PORT` | Puerto de PostgreSQL | `5432` |
| `DB_USER` | Usuario de la base de datos | `postgres` |
| `DB_PASSWORD` | Contraseña del usuario de la base de datos | `secret` |
| `DB_NAME` | Nombre de la base de datos | `libreria_mariela` |
| `ALLOWED_ORIGINS` | Orígenes permitidos para CORS (separados por coma) | `http://localhost:4200,http://localhost:3000` |
| `PEAK_AUTH_URL` | URL base del servidor de autenticación Peak Auth | `http://localhost:9009` |
| `PEAK_AUTH_CLIENT_ID` | Identificador de cliente registrado en Peak Auth | `libreria-mariela` |
| `PEAK_AUTH_CLIENT_SECRET` | Secreto de cliente OAuth registrado en Peak Auth | `tu_client_secret` |
| `ROOT_EMAIL` | *(Opcional - Docker)* Email de administrador para pgAdmin | `admin@admin.com` |
| `ROOT_PASSWORD` | *(Opcional - Docker)* Contraseña para pgAdmin | `admin123` |

---

## 🚀 Instalación y Puesta en Marcha

### Opción 1: Desarrollo Local (Go + PostgreSQL)

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/brunotarditi/libreria-mariela-api.git
   cd libreria-mariela-api
   ```

2. **Instalar dependencias de Go:**
   ```bash
   go mod download
   ```

3. **Configurar el archivo `.env`:**
   Asegúrate de que PostgreSQL esté corriendo y la base de datos exista.

4. **Iniciar el servidor:**
   ```bash
   go run main.go routes.go
   ```
   *La API correrá en `http://localhost:8080` y GORM ejecutará automáticamente las migraciones (`AutoMigrate`) de las tablas necesarias.*

---

### Opción 2: Despliegue con Docker Compose

Docker Compose levantará automáticamente la API de Go, la base de datos PostgreSQL y la interfaz web de pgAdmin.

1. **Asegurar el archivo `.env` configurado.**

2. **Construir y levantar los contenedores:**
   ```bash
   docker compose up --build -d
   ```

3. **Servicios disponibles:**
   - **API REST**: `http://localhost:8080`
   - **Base de datos PostgreSQL**: `localhost:5432`
   - **pgAdmin 4**: `http://localhost:5050` (accede con `ROOT_EMAIL` y `ROOT_PASSWORD`)

4. **Detener los servicios:**
   ```bash
   docker compose down
   ```

---

## 📡 Referencia de la API (Endpoints)

El prefijo global de la API es `/api/v1`.

### Autenticación

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `POST` | `/api/v1/auth/exchange` | Canjea un código de autorización OAuth 2.0 por un token de acceso en Peak Auth. |

**Body de ejemplo:**
```json
{
  "code": "codigo_de_autorizacion_oauth"
}
```

---

### Salud del Sistema

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/healthy` | Comprueba la disponibilidad de la API y el estado del servicio. |

---

### Dashboard

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/dashboard` | Retorna contadores generales (total de productos, clientes, proveedores) y registros recientes de auditoría. |

---

### Productos

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/products` | Lista todos los productos incluyendo sus categorías y marcas asociadas. |
| `GET` | `/api/v1/products/:id` | Obtiene un producto por su ID. |
| `POST` | `/api/v1/products` | Crea un nuevo producto. |
| `PUT` | `/api/v1/products/:id` | Actualiza un producto existente. |
| `DELETE` | `/api/v1/products/:id` | Elimina lógicamente un producto (soft-delete). |
| `GET` | `/api/v1/products/export` | Descarga el catálogo de productos en formato Excel (`.xlsx`). |
| `POST` | `/api/v1/products/import` | Importa productos masivamente desde un archivo Excel (`multipart/form-data`, campo `file`). |

**Body de creación de producto:**
```json
{
  "code": "CUAD-A4-01",
  "sku": "7791234567890",
  "name": "Cuaderno Universitario Rayado 100 Hojas",
  "profit_margin": 35.0,
  "description": "Tapa dura encuadernada",
  "category_id": 1,
  "brand_id": 2
}
```

---

### Categorías y Marcas

#### Categorías (`/api/v1/categories`)
- `GET /` - Listar todas las categorías.
- `GET /:id` - Obtener categoría por ID.
- `POST /` - Crear una categoría (`{"name": "Escritura"}`).
- `POST /list` - Crear múltiples categorías en lote (`[{"name": "Escritura"}, {"name": "Papelería"}]`).
- `PUT /:id` - Actualizar categoría.
- `DELETE /:id` - Eliminar categoría.

#### Marcas (`/api/v1/brands`)
- `GET /` - Listar todas las marcas.
- `GET /:id` - Obtener marca por ID.
- `POST /` - Crear una marca (`{"name": "Rivadavia"}`).
- `POST /list` - Crear múltiples marcas en lote (`[{"name": "Bic"}, {"name": "Faber-Castell"}]`).
- `PUT /:id` - Actualizar marca.
- `DELETE /:id` - Eliminar marca.

---

### Clientes y Proveedores

#### Clientes (`/api/v1/customers`)
- `GET /` - Listar clientes.
- `GET /:id` - Obtener cliente por ID.
- `POST /` - Registrar cliente:
  ```json
  {
    "name": "Colegio San Martín",
    "contact_info": "compras@colegiosanmartin.edu.ar - Tel: 011-4567-8900"
  }
  ```
- `PUT /:id` - Actualizar datos del cliente.
- `DELETE /:id` - Eliminar cliente.

#### Proveedores (`/api/v1/suppliers`)
- `GET /` - Listar proveedores.
- `GET /:id` - Obtener proveedor por ID.
- `POST /` - Registrar proveedor:
  ```json
  {
    "name": "Distribuidora Papelera Central",
    "contact_info": "ventas@papeleracentral.com - Tel: 011-4433-2211"
  }
  ```
- `PUT /:id` - Actualizar datos del proveedor.
- `DELETE /:id` - Eliminar proveedor.

---

### Compras

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/purchases` | Lista el historial completo de compras realizadas. |
| `GET` | `/api/v1/purchases/:id` | Obtiene el detalle de una compra específica. |
| `POST` | `/api/v1/purchases` | Registra una compra. Aumenta automáticamente el stock del producto e ingresa un movimiento de tipo `IN`. |
| `DELETE` | `/api/v1/purchases/:id` | Anula una compra. Revierte el stock (movimiento `OUT` de devolución de compra). |

**Body de registro de compra:**
```json
{
  "product_id": 1,
  "supplier_id": 2,
  "cost": 1500.50,
  "quantity": 50
}
```

---

### Ventas

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/sells` | Lista el historial de ventas realizadas. |
| `GET` | `/api/v1/sells/:id` | Obtiene el detalle de una venta por ID. |
| `POST` | `/api/v1/sells` | Registra una venta, calcula el costo promedio, descuenta stock y registra un movimiento `OUT`. |
| `DELETE` | `/api/v1/sells/:id` | Anula una venta. Revierte el stock (movimiento `IN` de devolución de venta). |

**Body de registro de venta:**
```json
{
  "product_id": 1,
  "customer_id": 1,
  "quantity": 5
}
```
*Nota: El precio unitario de la venta (`price`) y el costo promedio (`average_cost`) son determinados automáticamente por la API durante la transacción.*

---

### Stock y Movimientos

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/stocks` | Lista todos los movimientos de stock registrados (entradas, salidas y motivos). |
| `GET` | `/api/v1/stocks/:id` | Obtiene el detalle de un movimiento de stock por ID. |

---

### Listas de Precios

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/prices` | Lista los registros de precios. |
| `GET` | `/api/v1/prices/:id` | Obtiene un precio por ID. |
| `POST` | `/api/v1/prices` | Registra un nuevo precio para un producto con fecha de vigencia (`effective_at`). |

**Body:**
```json
{
  "product_id": 1,
  "price": 2800.00
}
```

---

### Presupuestos (Budgets)

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `GET` | `/api/v1/budget` | Genera y descarga dinámicamente un documento PDF (`presupuesto.pdf`) de presupuesto comercial. |
| `POST` | `/api/v1/budges` | Guarda un presupuesto en la base de datos (con cliente y lista de ítems en JSONB). |
| `DELETE` | `/api/v1/budges/:id` | Elimina un presupuesto registrado. |

**Body de guardado de presupuesto:**
```json
{
  "client_name": "Instituto Manuel Belgrano",
  "description": "Presupuesto para útiles escolares de inicio de ciclo",
  "items": ["Cuaderno A4 x 50", "Lapicera Azul x 100", "Regla 30cm x 50"]
}
```

---

## 💡 Lógica de Negocio Destacada

1. **Cálculo de Costo Promedio Ponderado**:
   Al registrar una venta, la API analiza las compras históricas cronológicas del producto y el stock disponible para determinar el valor unitario de reposición/costo actual.
2. **Margen de Utilidad Dinámico**:
   El precio final de la venta se calcula con el porcentaje `profit_margin` asignado en la entidad `Product`, asegurando la rentabilidad deseada en cada transacción.
3. **Consistencia Transaccional (ACID)**:
   Todas las operaciones de compra y venta utilizan transacciones de base de datos (`tx.Begin()`, `tx.Commit()`, `tx.Rollback()`). Si la actualización de stock o el movimiento fallan, toda la operación se descarta de inmediato para prevenir inconsistencias en el inventario.
4. **Auditoría No Bloqueante**:
   Gracias al modelo de concurrencia de Go, el middleware de auditoría delega la inserción de registros en la tabla `audit_logs` a Goroutines (`go db.Create(&audit)`), permitiendo responder a las solicitudes HTTP con mínima latencia.

---

## 📄 Licencia

Este proyecto es de uso privado para el negocio familiar **Librería Mariela**. Todos los derechos reservados.
