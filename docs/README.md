# Startlight Docs

## Project Structure

Inside of your Astro + Starlight project, you'll see the following folders and files:

```text
.
├── public/
├── src/
│   ├── assets/
│   ├── content/
│   │   └── docs/
│   └── content.config.ts
├── astro.config.mjs
├── package.json
└── tsconfig.json
```

Starlight looks for `.md` or `.mdx` files in the `src/content/docs/` directory. Each file is exposed as a route based on its file name.

Images can be added to `src/assets/` and embedded in Markdown with a relative link.

Static assets, like favicons, can be placed in the `public/` directory.

## Embedding Videos

The documentation site supports video embeds using the [astro-embed](https://astro-embed.netlify.app/) package. To embed videos in MDX files:

### YouTube Videos

```mdx
---
title: Your Page
---

import { YouTube } from 'astro-embed';

<YouTube id="dQw4w9WgXcQ" />
```

### Vimeo Videos

```mdx
import { Vimeo } from 'astro-embed';

<Vimeo id="1084537" />
```

### Other Supported Embeds

The astro-embed package also supports:
- Twitter/X posts (`Tweet` component)
- Mastodon posts (`MastodonPost` component)
- GitHub Gists (`Gist` component)
- Bluesky posts (`BlueskyPost` component)
- Generic link previews (`LinkPreview` component)

For complete documentation and component options, see the [astro-embed documentation](https://astro-embed.netlify.app/).

## 🧞 Commands

All commands are run from the root of the project, from a terminal:

| Command                   | Action                                           |
| :------------------------ | :----------------------------------------------- |
| `npm install`             | Installs dependencies                            |
| `npm run dev`             | Starts local dev server at `localhost:4321`      |
| `npm run build`           | Build your production site to `./dist/`          |
| `npm run preview`         | Preview your build locally, before deploying     |
| `npm run astro ...`       | Run CLI commands like `astro add`, `astro check` |
| `npm run astro -- --help` | Get help using the Astro CLI                     |

## Want to learn more?

Check out [Starlight’s docs](https://starlight.astro.build/), read [the Astro documentation](https://docs.astro.build), or jump into the [Astro Discord server](https://astro.build/chat).
