# Guía de Creación de Módulos (TinyWasm)

Esta guía documenta el proceso estándar para crear nuevos módulos en la plataforma `tinywasm`.

## Estructura General

Un módulo típico reside en `modules/<nombre_modulo>/` y contiene al menos los siguientes archivos:

-   `<nombre_modulo>.go`: Define la estructura principal y la implementación de interfaces.
-   `front.go`: Lógica específica de WebAssembly (compilada con `GOOS=js GOARCH=wasm`).
-   `back.go`: Lógica de backend y definiciones de interfaces (compilada por defecto).
-   `README.md`: Documentación del propósito y diseño del módulo.

## Requisitos de Implementación

Para que un módulo sea cargado por `tinywasm/site`, debe implementar al menos la interfaz básica que espera el registro de manejadores.

### 1. Definición del Struct

```go
package mimodulo

type MiModulo struct {
    // dependencias o estado
}

// Constructor para inyección de dependencias
func Add() []any {
    m := &MiModulo{}
    return []any{m}
}
```

### 2. Identificación del Módulo (HandlerName)

Debe implementar `HandlerName() string` para devolver un identificador único que se usará en el DOM (id) y en el enrutamiento.

```go
func (m *MiModulo) HandlerName() string {
    return "mimodulo"
}
```

### 3. Renderizado HTML (RenderHTML)

Debe implementar `RenderHTML() string` que devuelve el HTML inicial.

**Uso de Componentes (TinyWasm Components):**

La arquitectura exige la reutilización estricta de componentes. **No se debe escribir HTML crudo** en los módulos para elementos estructurales o de UI comunes.

1.  **Catálogo de Componentes**: Consultar `tinywasm/components` para ver los componentes disponibles (Botones, Paneles, Listas, Headers, etc.).
2.  **Creación de Nuevos Componentes**: Si un componente necesario no existe (por ejemplo, un tipo específico de Grilla o Card), **debe crearse primero en `tinywasm/components`** y luego importarse en el módulo. Esto asegura que la identidad visual y lógica se mantenga centralizada.
3.  **Implementación**: El método `RenderHTML` del módulo debe orquestar estos componentes llamando a sus métodos `Render()` o `String()`.

Ejemplo Correcto:

```go
func (m *MiModulo) RenderHTML() string {
    panel := &panel.Panel{
        ID: m.HandlerName(),
        Title: "Mi Título",
        Content: myContentComponent,
    }
    return panel.RenderHTML()
}
```

## Integración con TinyWasm Components

Se recomienda usar los componentes base de `tinywasm/components` cuando sea posible para botones, tarjetas y elementos comunes.

## Integración con CRUDP (Backend/Frontend Data Sync)

Como `tinywasm/site` utiliza `crudp` para la orquestación y el enrutamiento de datos, el struct principal del módulo debe cumplir con la interfaz `Handler` de `crudp`. Esto permite que el módulo sincronice datos automáticamente entre el servidor y el cliente WASM.

### Requisitos del Struct Principal

Además de `HandlerName` y `RenderHTML`, el struct debe implementar los siguientes métodos para el control de acceso y validación:

#### 1. ValidateData (Validación de Datos)

Se ejecuta antes de cualquier operación CRUD. Implementar aunque sea retornando `nil` si no aplica.

```go
// ValidateData valida los datos entrantes.
func (m *MiModulo) ValidateData(action byte, data ...any) error {
    // action: 'c' (create), 'r' (read), 'u' (update), 'd' (delete)
    return nil
}
```

#### 2. AllowedRoles (Control de Acceso)

Define qué roles tienen permiso para ejecutar acciones CRUD.

```go
// AllowedRoles retorna los roles permitidos para una acción.
// Retorna '*' para público, 'a' para admin, etc.
func (m *MiModulo) AllowedRoles(action byte) []byte {
    if action == 'r' {
        return []byte{'*'} // Lectura pública
    }
    return []byte{'a'} // Escritura solo admin
}
```

### Operaciones CRUD (Backend)

En el archivo `back.go`, se pueden implementar los métodos CRUD (`Create`, `Read`, `Update`, `Delete`) para manejar la lógica de negocio.

```go
//go:build !wasm

func (m *MiModulo) Read(data ...any) any {
    // Lógica para recuperar datos (BD, API, etc.)
    return []MiData{{ID: "1", Val: "Ejemplo"}}
}
```

### Actualización DOM (Frontend)

En el archivo `front.go`, estos mismos métodos reciben los datos del servidor para actualizar la UI.

```go
//go:build wasm

func (m *MiModulo) Read(data ...any) any {
    // Lógica para actualizar el DOM con los datos recibidos
    return nil
}
```
