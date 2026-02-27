# ✏️ Librería Mariela API - Sistema para negocio familiar

![Go Version](https://img.shields.io/badge/go-1.23.0-blue.svg)
![Gin Framework](https://img.shields.io/badge/gin-v1.10.0-blue.svg)
![PostgreSQL](https://img.shields.io/badge/postgresql-v1.6.0-blue.svg)

**Librería Mariela API** es la api para el negocio Librería Mariela, este negocio de útiles escolares actualmente cuenta con módulos de inventario, importar productos desde una hoja de excel y llevar el control del stock de los productos.

---

## 🛠️ Stack Tecnológico

- **Lenguaje**: Go (1.23.0)
- **Web Framework**: Gin
- **ORM**: GORM
- **Base de Datos**: PostgreSQL

## 🚀 Instalación y Desarrollo Local

### Requisitos Previos

- Go 1.23.0
- PostgreSQL 13+

### Pasos de instalación

1. **Clonar el repositorio:**

   ```bash
   git clone https://github.com/brunotarditi/libreria-mariela-api.git
   cd peak-auth
   ```

2. **Instalar dependencias:**

   ```bash
   go mod download
   ```

3. **Configurar variables de entorno:**
   Copia el archivo de ejemplo y configura tus datos (base de datos, puerto, etc.).

   ```bash
   cp .env.example .env
   ```

4. **Ejecutar el servidor:**
   ```bash
   go run main.go
   ```

