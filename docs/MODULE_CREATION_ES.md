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

#### 3. Renderizado HTML y Registro de Módulo (RenderHTML)

Para que el sistema de navegación y SSR reconozca una estructura como un módulo, **es obligatorio** usar el `module.New()`. Cualquier handler que no use este constructor no aparecerá en el menú principal ni se inyectará automáticamente en la carga inicial.

**Uso de Componentes (TinyWasm Components):**

La arquitectura exige la reutilización estricta de componentes. **No se debe escribir HTML crudo** en los módulos para elementos estructurales o de UI comunes.

1.  **Catálogo de Componentes**: Consultar `tinywasm/components` para ver los componentes disponibles (Botones, Paneles, Listas, Headers, etc.).
2.  **Creación de Nuevos Componentes**: Si un componente necesario no existe (por ejemplo, un tipo específico de Grilla o Card), **debe crearse primero en `tinywasm/components`** y luego importarse en el módulo. Esto asegura que la identidad visual y lógica se mantenga centralizada.
3.  **Implementación**: El método `RenderHTML` del módulo debe orquestar estos componentes llamando a sus métodos `Render()` o `String()`.

Ejemplo Correcto:

```go
// RenderHTML usa el Builder Pattern para ensamblar el módulo
func (m *MiModulo) RenderHTML() string {
    return module.New(m).
        SetTitle("Mi Título").
        SetPublic(). 
        WithHeader(anyControlsComponent).
        WithSearchList("my_list_id", itemsHTML).
        WithForm(
            form.New("my-form").
                AddInput("sku", "SKU", form.Text, true).
                AddInput("name", "Name", form.Text, true),
        ).
        WithButtons(
            &button.Button{Name: "btn_save", Title: "Guardar", Icon: "icon-btn-save"},
            &button.Button{Name: "btn_cancel", Title: "Cancelar"},
        ).
        Render()
}
```

> [!TIP]
> **Preferencia de Interfaces**: Use siempre los métodos que aceptan componentes (`WithForm`, `WithButtons`, `WithHeader`) en lugar de los que terminan en `Html`. Esto permite que el sistema detecte automáticamente el CSS e Iconos de cada componente sin cargarlos todos por defecto.

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

### 4. Separación de Responsabilidades (SSR.go)

**CRÍTICO:** Todo código relacionado con CSS, SVG e inyección de JS inicial debe residir en un archivo separado `ssr.go` (o `back.go`) con la etiqueta de compilación `//go:build !wasm`.

Esto es fundamental para:
1.  **Reducir el tamaño del binario WASM**: El cliente no necesita cargar strings de CSS o SVG que ya fueron renderizados por el servidor.
2.  **Seguridad**: Evita exponer lógica de servidor en el cliente.

**Ejemplo `ssr.go`:**

```go
//go:build !wasm

package mimodulo

import "github.com/tinywasm/site"

func (m *MiModulo) RenderCSS() string {
    // Usar variables estándar (ver tinywasm/components/palette.go)
    return ".mi-estilo { color: var(" + components.VarColorPrimary + "); }"
}

// ...
```

### 5. Estilos y Variables (Palette)

**REGLA:** No hardcodear colores (ej: `#fff`, `red`). Usar siempre las variables CSS definidas por el sistema.
Las variables disponibles están declaradas en `tinywasm/components/palette.go`:

-   `var(--color-primary)`
-   `var(--color-secondary)`
-   `var(--color-selection)`
-   etc.
