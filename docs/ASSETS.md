# Assets

`tinywasm/site` leverages `tinywasm/assetmin` to bundle all resources into a single set of files, minimizing HTTP requests and optimizing performance.

## Bundling Strategy

All registered Modules and their sub-components contribute to the global asset bundles during registration.

- **CSS**: Collected from `RenderCSS()`.
- **JS**: Collected from `RenderJS()` (if available).
- **Icons**: Collected from `IconSvg()`.

### Global Bundles

- `style.css`: All component styles combined.
- `script.js`: Global initialization and component-specific logic.
- `sprite.svg`: A single SVG sprite containing all icons, accessible via `<use href="#icon-id">`.

## Integration

The `site` library automatically manages the lifecycle of `assetmin`:

1. Detects `RenderCSS`, `RenderJS`, and `IconSvg` via reflection.
2. Passes the content to `assetmin` for minification and bundling.
3. Automatically includes the correct `<link>` and `<script>` tags in the SSR output.

---
**Status**: No Implemented
