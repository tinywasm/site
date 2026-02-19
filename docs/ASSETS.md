# Assets

`tinywasm/site` leverages `tinywasm/assetmin` to bundle all resources into a single set of files, minimizing HTTP requests and optimizing performance.

## Bundling Strategy

All registered Modules and their sub-components contribute to the global asset bundles during registration.

- **CSS**: Collected from `RenderCSS()`.
- **JS**: Collected from `RenderJS()`.
- **Icons**: Collected from `IconSvg()`.

### Global Bundles

- `style.css`: All component styles combined.
- `script.js`: Global initialization and component-specific logic.
- `sprite.svg`: A single SVG sprite containing all icons, accessible via `<use href="#icon-id">`.

## Extraction Logic

The `site` library automatically manages the lifecycle of `assetmin` during the registration phase in `register_ssr.go`.

1. **Type Assertion**: Registered handlers are checked against `CSSProvider`, `JSProvider`, and `IconSvgProvider` interfaces.
2. **Collection**: Content from `RenderCSS()`, `RenderJS()`, and `IconSvg()` is collected.
3. **Bundling**: The collected content is passed to `assetmin` for minification and bundling.
4. **Injection**: Correct `<link>`, `<script>`, and SVG sprite reference tags are automatically included in the SSR output.

Note: `IconSvg()` returns a `map[string]string` (not a slice), providing one icon per provider.
